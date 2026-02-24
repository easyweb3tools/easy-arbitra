# Easy Arbitra — Product Architecture

Last updated: 2026-02-24  
Status: Implemented architecture (backend + frontend + Docker Compose)

## 1. Product Definition

Easy Arbitra is a Polymarket wallet intelligence product focused on one operational goal:  
help users discover high-quality wallets and convert from browsing to follow actions.

The product combines:
- Wallet performance facts (PnL, activity, volume)
- Strategy labeling (e.g., market maker, event trader)
- Timing/edge analysis and anomaly detection
- AI-generated explanations and summaries
- Follow workflows (single wallet + batch portfolio follow)

## 2. Current Product Surfaces

Frontend routes (Next.js App Router):
- `/`: discovery home (featured wallets, starter portfolios, signals)
- `/wallets`: exploration and filtering of potential wallets
- `/wallets/[id]`: wallet profile with decision card, AI analysis, evidence
- `/watchlist`: tracked wallets + action-oriented event feed
- `/markets`: market browser
- `/anomalies` and `/anomalies/[id]`: anomaly feed and detail
- `/s/[id]`: share landing page optimized for conversion
- `/leaderboard`: redirected to `/wallets?sort_by=smart_score&order=desc`

## 3. System Architecture

## 3.1 High-level

- **Frontend**: Next.js 14 + TypeScript
- **API layer**: Gin (Go)
- **Data layer**: PostgreSQL + GORM + SQL-heavy repositories
- **Workers**: scheduled ingestion and analytics jobs
- **Reverse proxy**: Nginx (routes `/api/v1` to backend, UI to frontend)
- **Runtime**: Docker Compose for local full-stack orchestration

Request path:
1. Browser -> Nginx (`:3000`)
2. Nginx -> Next.js frontend or Gin API
3. API -> service layer -> repositories -> PostgreSQL
4. Background workers continuously update raw and derived data

## 3.2 Backend module layout

`backend/internal` domains:
- `api/handler`: HTTP handlers and request parsing
- `service`: business logic and view composition
- `repository`: SQL queries and persistence
- `worker`: scheduled syncers and detectors
- `client`: external data source adapters
- `ai`: AI analysis integration layer

Entrypoints:
- `cmd/server`: API + workers
- `cmd/migrate`: SQL migration runner

## 4. API Architecture (Current)

Base path: `/api/v1`

Core wallet APIs:
- `GET /wallets`
- `GET /wallets/potential`
- `GET /wallets/:id`
- `GET /wallets/:id/profile`
- `GET /wallets/:id/share-card`
- `GET /wallets/:id/share-landing`
- `GET /wallets/:id/decision-card`
- `GET /wallets/:id/explanations`
- `GET /wallets/:id/info-edge`

Watchlist and conversion APIs:
- `GET /watchlist`
- `POST /watchlist`
- `POST /watchlist/batch`
- `DELETE /watchlist/:wallet_id`
- `GET /watchlist/feed`
- `GET /watchlist/summary`

Portfolio and discovery APIs:
- `GET /portfolios`
- `GET /ops/highlights`
- `GET /stats/overview`
- `GET /leaderboard`

Market and anomaly APIs:
- `GET /markets`
- `GET /markets/:id`
- `GET /anomalies`
- `GET /anomalies/:id`
- `PATCH /anomalies/:id/acknowledge`

AI APIs:
- `POST /ai/analyze/:wallet_id`
- `GET /ai/report/:wallet_id`
- `GET /ai/report/:wallet_id/history`

## 5. Data Model (Operationally Relevant)

Main tables:
- `wallet`, `market`, `token`, `trade_fill`, `offchain_event`
- `wallet_features_daily`
- `wallet_score`
- `ai_analysis_report`
- `anomaly_alert`
- `watchlist`
- `wallet_update_event`
- `portfolio`
- `ingest_cursor`, `ingest_run`

Important product fields currently used in APIs:
- `wallet_score.pool_tier` (`observation`, `strategy`, `star`)
- `wallet_score.suitable_for`, `risk_level`, `suggested_position`, `momentum`
- `wallet_update_event.action_required`, `suggestion`, `suggestion_zh`
- `portfolio.wallet_ids` (JSON list for starter packs)

## 6. Data and Intelligence Pipelines

## 6.1 Ingestion pipeline

Workers ingest markets, trades, and offchain events via dedicated clients.  
Cursor/state tables are used for resumable incremental sync.

## 6.2 Feature + scoring pipeline

Daily/periodic jobs generate wallet features and classification scores.  
Classification currently writes both strategy and product-oriented metadata (pool tier, risk, suggested position, momentum).

## 6.3 AI analysis pipeline

AI analysis is wallet-centric:
- Triggered on demand (`/ai/analyze/:wallet_id`) and optionally by batch worker
- Stores model output in `ai_analysis_report`
- Emits watchlist update events (`ai_report`) for downstream feed UX

## 6.4 Anomaly pipeline

Anomaly detector scans wallets and produces:
- `anomaly_alert` records
- Watchlist feed entries mapped to action-oriented categories

## 7. Frontend Architecture

Code organization:
- `src/app`: route-level server components
- `src/components`: UI modules by domain (`wallet`, `watchlist`, `portfolio`, `share`, `ui`)
- `src/lib/api.ts`: typed API client
- `src/lib/types.ts`: API contract types

Current frontend behavior aligned with operations:
- Wallet exploration supports strategy/pool/AI filters and sorting
- Wallet detail starts with a decision card before deep analytics
- Watchlist separates action-required events from normal updates
- Share flow uses dedicated landing page (`/s/[id]`)
- Batch follow is exposed through starter portfolio cards

## 8. Runtime and Deployment

## 8.1 Local stack

`docker-compose.yml` services:
- `postgres`
- `backend`
- `frontend`
- `nginx`

Standard local bootstrap:
1. `docker compose up -d postgres backend frontend`
2. Run migrations (`cmd/migrate`)
3. Verify health:
   - `GET /healthz`
   - `GET /readyz`
4. Run API/UI smoke scripts

## 8.2 Known compose nuance

Frontend mounts `/app/node_modules` as a volume.  
After dependency changes, stale anonymous volumes can cause runtime module errors.  
If that happens, recreate frontend with volume cleanup.

## 9. Quality and Guardrails

Current guardrails in architecture:
- API middleware: request ID, CORS, rate limit, centralized error handling
- Strong typed API contracts on frontend
- Migration-first schema evolution (ordered SQL files)
- Health/readiness endpoints for operational checks
- Runtime smoke scripts for API and UI sanity

## 10. Current Limitations

- SQL-heavy aggregations can be expensive under larger datasets (several multi-CTE queries in critical paths)
- Leaderboard still exists as an API route while UI route is now redirect-based
- Some recommendation logic is rule-based and should evolve toward calibrated models
- Event suggestion text is template-driven; personalization depth is limited

## 11. Near-term Architecture Priorities

1. Add caching/materialized views for heavy wallet aggregates.
2. Introduce explicit query SLOs and slow-query budget per endpoint.
3. Evolve decision card generation into a dedicated precompute pipeline.
4. Improve observability (structured metrics + dashboards per service and worker).
5. Harden Docker smoke scripts to avoid false negatives from shell pipeline behavior.
