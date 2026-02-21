# Repository Guidelines

## Project Structure & Module Organization
- `backend/`: Go API, workers, data ingestion, and AI analysis.
  - `cmd/server`: HTTP server entrypoint.
  - `cmd/migrate`: SQL migration runner.
  - `internal/`: domain modules (`api`, `service`, `repository`, `worker`, `client`, `ai`).
  - `migrations/`: ordered SQL migrations (`001_*.sql`, `002_*.sql`, ...).
- `frontend/`: Next.js App Router UI.
  - `src/app/`: route pages (`/wallets`, `/markets`, `/anomalies`, etc.).
  - `src/components/`: reusable UI components.
  - `src/lib/`: API client and shared types.
- `scripts/`: local smoke/bootstrap scripts.
- `ops/`: backup and monitoring helpers.
- `.github/workflows/`: CI, release, and smoke pipelines.

## Build, Test, and Development Commands
- Backend:
  - `cd backend && make run`: run API locally.
  - `cd backend && make test`: run Go tests.
  - `cd backend && make migrate`: apply SQL migrations.
- Frontend:
  - `cd frontend && npm run dev`: run UI dev server.
  - `cd frontend && npm run build`: production build check.
  - `cd frontend && npm run test:e2e`: Playwright smoke tests.
- Full stack (recommended):
  - `./scripts/bootstrap_dev.sh`: start Docker stack, migrate, run API/UI smoke.

## Coding Style & Naming Conventions
- Go: run `gofmt` before commit; keep packages focused by domain.
- TypeScript/React: strict types, small components, server-first pages unless interactivity is needed.
- Naming:
  - Files: snake_case in Go (`info_edge_service.go`), route folders in Next.js (`wallets/[id]`).
  - Migrations: numeric prefix + purpose (`003_anomaly_and_feature_ext.sql`).

## Testing Guidelines
- Go tests use `testing` package; files end with `_test.go`.
- Keep test names explicit: `Test<Module><Behavior>`.
- Add tests when changing parsing, scoring, routing, or error handling.
- Minimum local gate before PR: backend tests + frontend build + smoke scripts.

## Commit & Pull Request Guidelines
- Commit format: short imperative summary (e.g., `Add info-edge API and anomaly detail route`).
- Prefer small, focused commits by concern (API, UI, migration, CI).
- PR must include:
  - What changed and why.
  - Migration/config impact.
  - Verification steps and command output summary.
  - Screenshots/GIFs for UI changes (`/wallets`, `/anomalies`, etc.).

## Security & Configuration Tips
- Never commit secrets. Use env vars (`DATABASE_*`, `AWS_*`) and `.env` from `.env.example`.
- Default to `DATABASE_AUTO_MIGRATE=false` in Docker; use `cmd/migrate` for deterministic schema changes.
