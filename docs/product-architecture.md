# Easy Arbitra — Product Architecture

Last updated: 2026-03-02  
Status: Simplified architecture — focused on leaderboard + daily pick

## 1. Product Definition

Easy Arbitra is a Polymarket wallet intelligence product with one clear value proposition:  
**Discover profitable traders, recommend the best one daily, and track follow PnL.**

Core capabilities:
- **Leaderboard** — ranked list of profitable wallets by SmartScore
- **Daily Pick** — AI-recommended best trader each day (powered by Amazon Nova)
- **Follow PnL Tracking** — next-day performance of recommended trader's trades
- **Wallet Profiles** — detailed stats, positions, and trade history per wallet

## 2. Frontend Routes

Next.js 14 App Router (6 routes):

| Route | Purpose |
|-------|---------|
| `/` | Home — daily pick banner + leaderboard preview |
| `/daily-picks` | Today's recommended trader + history with follow PnL |
| `/leaderboard` | Full leaderboard with sorting |
| `/wallets` | Wallet explorer with strategy/tier filters |
| `/wallets/[id]` | Wallet profile: stats, positions, trade history |
| `/markets` | Market browser |

## 3. System Architecture

### 3.1 Stack

- **Frontend**: Next.js 14 + TypeScript + Tailwind CSS
- **API**: Gin (Go)
- **Database**: PostgreSQL + GORM
- **Workers**: scheduled background jobs
- **AI**: Amazon Nova (via Bedrock or Dev API) for daily pick reasoning
- **Reverse proxy**: Nginx
- **Runtime**: Docker Compose

### 3.2 Request Flow

```
Browser → Nginx (:3000) → Next.js (SSR) or Gin API (/api/v1)
                                              ↓
                                     Service → Repository → PostgreSQL
```

### 3.3 Backend Layout

```
backend/
├── cmd/server/          # Entrypoint: API + workers
├── config/              # Configuration (YAML + env)
├── internal/
│   ├── ai/              # Amazon Nova integration (Bedrock/DevAPI)
│   ├── api/handler/     # HTTP handlers
│   ├── api/middleware/   # Request ID, CORS, rate limit
│   ├── client/          # External API clients (Polymarket)
│   ├── model/           # GORM models
│   ├── repository/      # SQL queries
│   ├── service/         # Business logic
│   └── worker/          # Background jobs
└── pkg/                 # Shared utilities (logger, polyaddr, response)
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
| `daily_pick` | Daily recommended trader + follow PnL |
| `ingest_cursor` | Resumable sync state |
| `ingest_run` | Worker run audit log |

### Key Fields

- `wallet_score.smart_score` — core ranking metric
- `wallet_score.strategy_type` — classification label
- `wallet_score.pool_tier` — `star` / `strategy` / `observation`
- `daily_pick.reason_json` — Nova-generated analysis
- `daily_pick.follow_pnl` — next-day follow result

## 6. Workers

| Worker | Interval | Purpose |
|--------|----------|---------|
| `MarketSyncer` | 10m | Sync markets from Polymarket |
| `TradeSyncer` | 5m | Sync recent trades |
| `TradeBackfillSyncer` | 5m | Backfill historical trades |
| `FeatureBuilder` | 30m | Compute wallet features |
| `ScoreCalculator` | 1h | Classify & score wallets |
| `DailyRecommender` | 6h | Pick best trader + Nova reasoning + backfill PnL |

### DailyRecommender Flow

1. **Backfill** — calculate yesterday's pick follow PnL
2. **Select** — top 10 from leaderboard, exclude last 7 days' picks
3. **Analyze** — call Amazon Nova for recommendation reasoning
4. **Write** — save `daily_pick` record

## 7. AI Integration

Amazon Nova is used for generating daily pick recommendations:
- Provider: Bedrock SDK or Dev API (configurable)
- Model: `nova-pro-v1`
- Input: wallet stats (PnL, trades, strategy, score, volume)
- Output: structured analysis JSON + natural language summary
- Fallback: template-based summary if Nova is unavailable

## 8. Runtime

### Docker Compose

```yaml
services: [postgres, backend, frontend, nginx]
```

Bootstrap:
1. `docker compose up -d`
2. Auto-migrate on startup
3. Verify: `GET /healthz`, `GET /readyz`

### Configuration

Key env vars / config:
- `WORKER_ENABLED` — enable background workers
- `WORKER_DAILY_RECOMMENDER_INTERVAL` — pick frequency (default 6h)
- `NOVA_ENABLED` — enable AI analysis
- `NOVA_PROVIDER` — `bedrock` or `devapi`
- `NOVA_API_KEY` — API key for Dev API mode

## 9. Removed Features

The following were removed during the Mar 2026 simplification:
- ~~Authentication (Google OAuth + JWT)~~
- ~~Copy-trading (AI agent auto-execution)~~
- ~~Watchlist + event feed~~
- ~~Anomaly detection~~
- ~~Portfolio starter packs~~
- ~~Share landing pages~~
- ~~Wallet decision cards~~
- ~~AI batch analysis worker~~
