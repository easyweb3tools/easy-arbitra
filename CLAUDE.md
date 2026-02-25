# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

Easy Arbitra — an AI-powered Polymarket wallet intelligence platform. Product flow: discover high-quality wallets → understand why → follow → get actionable updates.

## Tech Stack

- **Frontend:** Next.js 14 (App Router), TypeScript 5.5, Tailwind CSS 3.4, lucide-react icons
- **Backend:** Go 1.23, Gin 1.10, GORM 1.30, PostgreSQL 16
- **AI:** AWS Bedrock (Nova models) for wallet analysis
- **Infrastructure:** Docker Compose, Nginx reverse proxy
- **Testing:** Go `testing` package, Playwright (E2E)
- **i18n:** English (`en`) and Simplified Chinese (`zh`) via cookie-based locale

## Development Commands

### Full Stack (recommended)
```bash
./scripts/bootstrap_dev.sh    # Start Docker stack, migrate, smoke test
```
App at `http://localhost:3000`, backend at `http://localhost:8080`.

### Backend (Go)
```bash
cd backend && make run        # Run API with air hot-reload
cd backend && make test       # Run all Go tests
cd backend && make migrate    # Apply SQL migrations
cd backend && make build      # Build binary to bin/
cd backend && make tidy       # go mod tidy
```

### Frontend (Next.js)
```bash
cd frontend && npm run dev        # Dev server on port 3000
cd frontend && npm run build      # Production build
cd frontend && npm run lint       # ESLint
cd frontend && npm run test:e2e   # Playwright tests
```

### Verification Before PR
```bash
cd backend && go test ./...
cd frontend && npm run build
./scripts/e2e_api_smoke.sh
./scripts/e2e_ui_smoke.sh
```

## Architecture

### Backend Layers
Handler (`internal/api/handler/`) → Service (`internal/service/`) → Repository (`internal/repository/`) → PostgreSQL

- **Handlers:** parse HTTP requests, validate params, call services
- **Services:** business logic (wallet, market, AI, anomaly, watchlist, portfolio, stats, explanation, info-edge, product)
- **Repositories:** GORM queries
- **Workers** (`internal/worker/`): background jobs for market sync, trade ingestion, feature building, score calculation, anomaly detection, AI batch analysis
- **Clients** (`internal/client/`): external API wrappers (Gamma, DataAPI, Offchain for Polymarket)
- **Config:** Viper-based (`config/config.go`), env vars override `config/config.yaml`

### Frontend Structure
- `src/app/` — App Router pages (server components by default)
- `src/components/` — domain-organized: `wallet/`, `watchlist/`, `portfolio/`, `share/`, `ai/`, `anomaly/`, `ui/`, `layout/`
- `src/lib/api.ts` — typed API client (server-side uses `API_SERVER_BASE_URL`, client uses `NEXT_PUBLIC_API_BASE_URL`)
- `src/lib/types.ts` — shared API contract types
- `src/lib/i18n.ts` — translation dictionaries (~250 keys, en/zh)
- `src/lib/fingerprint.ts` — user fingerprinting for watchlist (sent as `X-User-Fingerprint` header)

### Key Routes
| Frontend | Backend API |
|----------|-------------|
| `/` | `/ops/highlights`, `/portfolios`, `/stats/overview` |
| `/wallets` | `/wallets`, `/wallets/potential` |
| `/wallets/[id]` | `/wallets/:id/profile`, `/wallets/:id/decision-card` |
| `/watchlist` | `/watchlist`, `/watchlist/feed`, `/watchlist/summary` |
| `/s/[id]` | `/wallets/:id/share-landing`, `/wallets/:id/share-card` |

### Request Flow (Docker Compose)
Nginx (`:3000`) → routes `/api/` to backend (`:8080`), everything else to frontend (`:3000` internal)

### Database
- Migrations in `backend/migrations/` (numbered SQL: `001_*.sql` through `006_*.sql`)
- Run via `cmd/migrate`; `DATABASE_AUTO_MIGRATE=false` by default
- Key tables: `wallet`, `market`, `token`, `trade_fill`, `wallet_features_daily`, `wallet_score`, `ai_analysis_report`, `anomaly_alert`, `watchlist`, `wallet_update_event`, `portfolio`
- Key product fields: `wallet_score.pool_tier` (observation/strategy/star), `suitable_for`, `risk_level`, `momentum`

### No Auth
Uses `X-User-Fingerprint` header for watchlist identity — no traditional authentication system. Rate limiting: 30 req/60s per IP.

## Conventions

- **Go files:** snake_case (`info_edge_service.go`), run `gofmt` before commit
- **Go tests:** `Test<Module><Behavior>` naming, files end `_test.go`
- **TypeScript/React:** strict types, server components by default, client interactivity only when needed
- **Migrations:** numeric prefix + purpose (`006_product_rebuild.sql`)
- **Commits:** short imperative summary (e.g., `Add decision card and portfolio APIs`)
- **Path alias:** `@/*` maps to `src/*` in frontend TypeScript
- **Docker caveat:** frontend mounts `/app/node_modules` as anonymous volume; if deps change and frontend fails, recreate the service with volume cleanup
