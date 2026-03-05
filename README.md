# Easy Arbitra

A Polymarket wallet intelligence platform powered by **Amazon Nova as the analytical brain**.  
Nova continuously evaluates profitable traders, decides the best one to recommend daily, and validates its own picks.

## How It Works

```
⏰ Every hour (UTC 08:00–22:00)
│
├─ 1. Wake Nova
├─ 2. Collect Top 20 candidates from leaderboard
├─ 3. Load Nova's memory (prior rounds today)
├─ 4. Load yesterday's validation result
├─ 5. Nova evaluates + ranks candidates
├─ 6. Nova decides: keep analyzing or make final pick
├─ 7. Save analysis to nova_session (memory)
│
└─ When Nova says "final":
   ├─ Write daily_pick
   └─ Next day: backfill follow PnL → feed back to Nova
```

## Key Features

| Feature | Description |
|---------|-------------|
| **Nova Thinking Timeline** | Watch Nova's hourly analysis rounds in real-time |
| **Daily Pick** | Nova-selected best trader with rationale (EN + ZH) |
| **Follow PnL** | Next-day validation of Nova's recommendation |
| **Leaderboard** | SmartScore-ranked profitable wallets |
| **Wallet Profiles** | Stats, positions, trade history per wallet |

## Tech Stack

| Layer | Technology |
|-------|-----------|
| AI Brain | Amazon Nova (Bedrock / Dev API) |
| Frontend | Next.js 14, TypeScript, Tailwind CSS |
| Backend | Go, Gin, GORM |
| Database | PostgreSQL |
| Runtime | Docker Compose + Nginx |

## Frontend Routes

| Route | Purpose |
|-------|---------|
| `/` | Home — daily pick banner + leaderboard |
| `/daily-picks` | Nova thinking timeline + pick + history |
| `/leaderboard` | Full leaderboard |
| `/wallets` | Wallet explorer |
| `/wallets/[id]` | Wallet detail |
| `/markets` | Market browser |

## API Endpoints

Base: `/api/v1`

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/wallets/potential` | Profitable wallets |
| GET | `/wallets/:id/profile` | Wallet profile |
| GET | `/wallets/:id/trades` | Trade history |
| GET | `/wallets/:id/positions` | Current positions |
| GET | `/leaderboard` | SmartScore leaderboard |
| GET | `/markets` | Market browser |
| GET | `/daily-pick` | Today's Nova pick |
| GET | `/daily-pick/history` | Pick history + follow PnL |
| GET | `/nova/sessions` | Nova's analysis rounds |
| GET | `/stats/overview` | Platform stats |

## Repository Structure

```
├── backend/
│   ├── cmd/server/           # Entrypoint
│   ├── config/               # YAML + env config
│   └── internal/
│       ├── ai/               # Nova integration (Orchestrate + AnalyzeWallet)
│       ├── api/              # Gin handlers + router
│       ├── client/           # Polymarket API clients
│       ├── model/            # GORM models (Wallet, NovaSession, DailyPick...)
│       ├── repository/       # SQL queries
│       ├── service/          # Business logic
│       └── worker/           # NovaOrchestrator + data syncers
├── frontend/
│   ├── src/app/              # Next.js App Router pages
│   ├── src/components/       # UI components
│   └── src/lib/              # API client, types, i18n
├── scripts/                  # Bootstrap + smoke tests
├── ops/                      # Nginx config
└── docs/                     # Architecture docs
```

## Local Development

### Prerequisites
- Docker + Docker Compose
- (Optional) Go 1.22+ and Node.js 20+ for non-container runs

### Quick Start
```bash
./scripts/bootstrap_dev.sh
```

### Manual
```bash
# Backend
cd backend && make run    # or: go run ./cmd/server
cd backend && make test

# Frontend
cd frontend && npm run dev
cd frontend && npm run build
```

## Configuration

Key environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `NOVA_ENABLED` | `false` | Enable Amazon Nova AI |
| `NOVA_PROVIDER` | `devapi` | `bedrock` or `devapi` |
| `NOVA_API_KEY` | — | API key (Dev API mode) |
| `WORKER_NOVA_ORCHESTRATOR_INTERVAL` | `1h` | Analysis frequency |
| `WORKER_NOVA_ORCHESTRATOR_START_HOUR` | `8` | UTC start hour |
| `WORKER_NOVA_ORCHESTRATOR_END_HOUR` | `22` | UTC end hour |
| `DATABASE_AUTO_MIGRATE` | `false` | Auto-create tables on startup |

## Verification

```bash
cd backend  && go build ./... && go vet ./... && go test ./...
cd frontend && npm run build
```

## Docker Compose

- App: `http://localhost:3000`
- Health: `http://localhost:8080/healthz`
- Ready: `http://localhost:8080/readyz`

> If frontend reports missing module errors after dependency changes, recreate the service with volume cleanup.
