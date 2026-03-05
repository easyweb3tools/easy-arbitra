# Easy Arbitra — Product Architecture

Last updated: 2026-03-05  
Status: Nova-as-Brain architecture — Nova drives analysis and daily pick decisions

## 1. Product Definition

Easy Arbitra is a Polymarket wallet intelligence product with one clear value proposition:  
**Amazon Nova is the analytical brain that continuously evaluates profitable traders, decides the best one to recommend daily, and validates its own picks.**

Core capabilities:
- **Leaderboard** — ranked list of profitable wallets by SmartScore
- **Nova-Driven Daily Pick** — Nova analyzes candidates hourly, builds memory across rounds, and self-determines when to make its final recommendation
- **Follow PnL Tracking** — next-day validation of Nova's pick, fed back as learning signal
- **Wallet Profiles** — detailed stats, positions, and trade history per wallet
- **Nova Thinking Timeline** — real-time visibility into Nova's analysis process

## 2. Frontend Routes

Next.js 14 App Router (6 routes):

| Route | Purpose |
|-------|---------|
| `/` | Home — daily pick banner + leaderboard preview |
| `/daily-picks` | Nova thinking timeline + today's pick + history with follow PnL |
| `/leaderboard` | Full leaderboard with sorting |
| `/wallets` | Wallet explorer with strategy/tier filters |
| `/wallets/[id]` | Wallet profile: stats, positions, trade history |
| `/markets` | Market browser |

## 3. System Architecture

### 3.1 Stack

- **Frontend**: Next.js 14 + TypeScript + Tailwind CSS
- **API**: Gin (Go)
- **Database**: PostgreSQL + GORM
- **AI Brain**: Amazon Nova (via Bedrock or Dev API)
- **Workers**: scheduled background jobs
- **Reverse proxy**: Nginx
- **Runtime**: Docker Compose

### 3.2 Nova-as-Brain Architecture

```
┌──────────────────────────────────────────────────────┐
│                   ⏰ Hourly Timer                     │
│               (UTC 08:00 – 22:00)                    │
└──────────────┬───────────────────────────────────────┘
               ▼
┌──────────────────────────────────────────────────────┐
│           NovaOrchestrator Worker                     │
│                                                      │
│  1. Backfill yesterday's follow PnL                  │
│  2. Load Nova's memory (prior rounds today)          │
│  3. Collect Top 20 candidates from leaderboard       │
│  4. Load yesterday's validation result               │
│  5. ─── Call Nova Orchestrate ──────────────────┐    │
│                                                 │    │
│     ┌───────────────────────────────────────┐   │    │
│     │  🧠 Amazon Nova                       │   │    │
│     │                                       │   │    │
│     │  Receives: candidates + memory +      │   │    │
│     │            yesterday's feedback       │   │    │
│     │                                       │   │    │
│     │  Decides:  "analyzing" → save memo    │   │    │
│     │            "final" → pick winner      │   │    │
│     └───────────────────────────────────────┘   │    │
│                                                 │    │
│  6. Save nova_session (memory)          ◄───────┘    │
│  7. If final → write daily_pick                      │
└──────────────────────────────────────────────────────┘
```

**Key design**: Nova is called hourly and **it decides** whether to continue analyzing or make its final pick. It can finalize early if confident, or use all rounds up to the window end.

### 3.3 Backend Layout

```
backend/
├── cmd/server/          # Entrypoint: API + workers
├── config/              # Configuration (YAML + env)
├── internal/
│   ├── ai/              # Amazon Nova integration
│   │   ├── bedrock_client.go   # Analyzer: AnalyzeWallet + Orchestrate
│   │   ├── orchestrate.go      # Types, prompts, response parser
│   │   └── orchestrate_impl.go # Orchestrate for Bedrock/DevAPI/Mock
│   ├── api/handler/     # HTTP handlers
│   ├── api/middleware/   # Request ID, CORS, rate limit
│   ├── client/          # External API clients (Polymarket)
│   ├── model/           # GORM models
│   ├── repository/      # SQL queries
│   ├── service/         # Business logic
│   └── worker/          # Background jobs (NovaOrchestrator)
└── pkg/                 # Shared utilities
```

## 4. API Endpoints

Base path: `/api/v1`

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/wallets` | List tracked wallets |
| GET | `/wallets/potential` | List profitable wallets (with filters) |
| GET | `/wallets/:id/profile` | Wallet profile + strategy + facts |
| GET | `/wallets/:id/trades` | Trade history (paginated) |
| GET | `/wallets/:id/positions` | Current positions |
| GET | `/leaderboard` | SmartScore leaderboard |
| GET | `/markets` | Market browser |
| GET | `/stats/overview` | Platform stats |
| GET | `/daily-pick` | Today's recommended trader |
| GET | `/daily-pick/history` | Daily pick history with follow PnL |
| GET | `/nova/sessions` | Nova's analysis rounds for a given date |
| GET | `/healthz` | Health check |
| GET | `/readyz` | Readiness check |

## 5. Data Model

### Active Tables

| Table | Purpose |
|-------|---------|
| `wallet` | Tracked Polymarket wallets |
| `market` | Indexed markets |
| `token` | Market outcome tokens |
| `trade_fill` | Individual trades |
| `wallet_features_daily` | Computed daily features |
| `wallet_score` | SmartScore + strategy classification |
| `nova_session` | **Nova's working memory** — one row per hourly analysis round |
| `daily_pick` | Daily recommended trader + follow PnL |
| `ingest_cursor` | Resumable sync state |
| `ingest_run` | Worker run audit log |

### nova_session (Nova's Brain Memory)

| Field | Type | Purpose |
|-------|------|---------|
| `session_date` | date | Analysis day |
| `round` | int | Round number (1 = first hour) |
| `phase` | text | `analyzing` / `final` / `verified` |
| `candidates_json` | jsonb | Nova's ranking of candidates |
| `observations_json` | jsonb | Nova's observation notes (recalled next round) |
| `decision_json` | jsonb | Full Nova response |
| `picked_wallet_id` | int | Set when phase = `final` |

## 6. Workers

| Worker | Interval | Purpose |
|--------|----------|---------|
| `MarketSyncer` | 10m | Sync markets from Polymarket |
| `TradeSyncer` | 5m | Sync recent trades |
| `TradeBackfillSyncer` | 5m | Backfill historical trades |
| `FeatureBuilder` | 30m | Compute wallet features |
| `ScoreCalculator` | 1h | Classify & score wallets |
| **`NovaOrchestrator`** | **1h** | **Wake Nova → analyze → memory → pick** |

### NovaOrchestrator Flow (per hour)

1. **Backfill** — calculate yesterday's pick follow PnL
2. **Window check** — skip if outside UTC 08:00–22:00
3. **Final check** — skip if today already has a `phase=final` session
4. **Collect** — top 20 candidates from leaderboard with PnL data
5. **Memory** — load all prior `nova_session` rounds for today
6. **Feedback** — load yesterday's validation result
7. **Orchestrate** — call Nova with candidates + memory + feedback
8. **Store** — save `nova_session` record
9. **If final** — also write `daily_pick` with Nova's chosen wallet

## 7. AI Integration (Amazon Nova)

Nova is the **central decision maker**, not just an explainer.

### Analyzer Interface

```go
type Analyzer interface {
    AnalyzeWallet(ctx, in WalletAnalysisInput) (*WalletAnalysisOutput, error)
    Orchestrate(ctx, in OrchestrateInput) (*OrchestrateOutput, error)  // NEW
}
```

### Nova's Orchestrate Prompt Design

- **System**: "You are the analytical brain of a Polymarket trading intelligence system"
- **Inputs**: candidates (top 20), memory (prior rounds), yesterday's result, round info
- **Self-determination**: Nova decides `analyzing` vs `final` (forced final on last round)
- **Memory**: observations from each round stored in `nova_session`, recalled in next round
- **Feedback loop**: yesterday's follow PnL injected as learning signal

### Providers

| Provider | Config |
|----------|--------|
| Bedrock | AWS SDK Converse API |
| Dev API | OpenAI-compatible chat completion |
| Mock | Deterministic fallback for testing |

## 8. Runtime

### Configuration

Key env vars:
- `WORKER_NOVA_ORCHESTRATOR_INTERVAL` — analysis frequency (default `1h`)
- `WORKER_NOVA_ORCHESTRATOR_START_HOUR` — UTC start hour (default `8`)
- `WORKER_NOVA_ORCHESTRATOR_END_HOUR` — UTC end hour (default `22`)
- `NOVA_ENABLED` — enable AI analysis
- `NOVA_PROVIDER` — `bedrock` or `devapi`
- `NOVA_API_KEY` — API key for Dev API mode

### Docker Compose

```yaml
services: [postgres, backend, frontend, nginx]
```

