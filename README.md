# Easy Arbitra

An AI-powered Polymarket wallet intelligence platform focused on operational conversion:
**discover high-quality wallets -> understand why -> follow quickly -> get actionable updates**.

## Product Architecture (Current)

### Frontend (Next.js 14)
Main routes:
- `/`: discovery home (featured wallets, starter portfolios, real-time signals)
- `/wallets`: wallet exploration with strategy/pool tier/AI filters + sorting
- `/wallets/[id]`: decision card, AI report, evidence snapshot, share card
- `/watchlist`: portfolio summary + action-required feed + normal updates
- `/markets`: market list
- `/anomalies` and `/anomalies/[id]`: anomaly feed and detail
- `/s/[id]`: share landing page optimized for external conversion
- `/leaderboard`: redirected to `/wallets?sort_by=smart_score&order=desc`

### Backend (Go + Gin)
Base path: `/api/v1`

Core APIs:
- Wallet exploration and profile:
  - `GET /wallets`
  - `GET /wallets/potential`
  - `GET /wallets/:id/profile`
  - `GET /wallets/:id/decision-card`
  - `GET /wallets/:id/share-card`
  - `GET /wallets/:id/share-landing`
  - `GET /wallets/:id/explanations`
  - `GET /wallets/:id/info-edge`
- Watchlist conversion and action feed:
  - `GET /watchlist`
  - `POST /watchlist`
  - `POST /watchlist/batch`
  - `DELETE /watchlist/:wallet_id`
  - `GET /watchlist/feed`
  - `GET /watchlist/summary`
- Discovery:
  - `GET /ops/highlights`
  - `GET /portfolios`
  - `GET /stats/overview`
  - `GET /leaderboard`
- Market and anomaly:
  - `GET /markets`
  - `GET /markets/:id`
  - `GET /anomalies`
  - `GET /anomalies/:id`
  - `PATCH /anomalies/:id/acknowledge`
- AI:
  - `POST /ai/analyze/:wallet_id`
  - `GET /ai/report/:wallet_id`
  - `GET /ai/report/:wallet_id/history`

### Data and workers
Key tables:
- `wallet`, `market`, `token`, `trade_fill`, `offchain_event`
- `wallet_features_daily`, `wallet_score`, `ai_analysis_report`
- `anomaly_alert`, `watchlist`, `wallet_update_event`, `portfolio`
- `ingest_cursor`, `ingest_run`

Important product fields:
- `wallet_score.pool_tier`, `suitable_for`, `risk_level`, `suggested_position`, `momentum`
- `wallet_update_event.action_required`, `suggestion`, `suggestion_zh`
- `portfolio.wallet_ids`

Workers continuously ingest market/trade/offchain data and update features, scores, anomalies, and AI-ready signals.

## Tech Stack
- Frontend: Next.js, TypeScript, Tailwind CSS
- Backend: Go, Gin, GORM
- Database: PostgreSQL
- Runtime: Docker Compose + Nginx

## Repository Structure
- `backend/`: API, services, repositories, workers, migrations
- `frontend/`: App Router UI, components, typed API client
- `scripts/`: bootstrap and smoke scripts
- `ops/`: nginx config and operational helpers
- `docs/`: product and architecture docs

## Local Development

### Prerequisites
- Docker + Docker Compose
- (Optional) Go and Node.js for non-container local runs

### One-command bootstrap (recommended)
```bash
./scripts/bootstrap_dev.sh
```
This will:
1. Start `postgres`, `backend`, `frontend` containers.
2. Run DB migrations.
3. Verify backend health/readiness.
4. Run API/UI smoke scripts.

### Manual commands
Backend:
```bash
cd backend
make run
make test
make migrate
```

Frontend:
```bash
cd frontend
npm run dev
npm run build
```

## Docker Compose Notes
- Public app endpoint: `http://localhost:3000`
- Backend health endpoints:
  - `http://localhost:8080/healthz`
  - `http://localhost:8080/readyz`
- If frontend runtime reports missing module errors after dependency updates, recreate frontend service with volume cleanup (`node_modules` volume can be stale).

## Verification Checklist
Minimum local gate before PR:
1. `cd backend && go test ./...`
2. `cd frontend && npm run build`
3. `./scripts/e2e_api_smoke.sh`
4. `./scripts/e2e_ui_smoke.sh`

## Migrations and Config
- Keep `DATABASE_AUTO_MIGRATE=false` in Docker.
- Use `cmd/migrate` / migration files for deterministic schema changes.
- Never commit secrets; use `.env` based on `.env.example`.
