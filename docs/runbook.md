# Easy Arbitra Runbook

## 1. Local Start
1. Copy env: `cp .env.example .env`
2. Start stack: `docker compose up -d`
3. Check services: `docker compose ps`
4. Open UI (via nginx): `http://localhost:3000`
5. Backend health: `curl -sf http://localhost:8080/healthz`
6. Backend ready: `curl -sf http://localhost:8080/readyz`

## 2. Database Migration
1. Apply migrations: `cd backend && make migrate`

## 3. One-Command Bootstrap
- Run `./scripts/bootstrap_dev.sh` to start containers, migrate DB, and execute smoke tests.

## 4. Common Operations
- Restart backend only: `docker compose restart backend`
- Recreate frontend + nginx after env/proxy changes: `docker compose up -d --force-recreate frontend nginx`
- View logs: `docker compose logs -f backend`
- Rebuild frontend: `docker compose build frontend && docker compose up -d frontend`

## 5. Operations API Smoke Tests
- Highlights: `curl -sf 'http://localhost:3000/api/v1/ops/highlights?limit=5' | jq '.data.as_of,.data.new_potential_wallets_24h'`
- Potential wallets: `curl -sf 'http://localhost:3000/api/v1/wallets/potential?page=1&page_size=5&min_trades=100&min_realized_pnl=0' | jq '.data.pagination,.data.items[0].wallet.id'`
- Share card: `curl -sf 'http://localhost:3000/api/v1/wallets/<wallet_id>/share-card' | jq '.data.wallet.id,.data.smart_score,.data.has_ai_report'`

### Watchlist APIs
1. Add:
   - `curl -sf -X POST 'http://localhost:3000/api/v1/watchlist' -H 'Content-Type: application/json' -H 'X-User-Fingerprint: demo-fp' -d '{"wallet_id":11443}'`
2. List:
   - `curl -sf 'http://localhost:3000/api/v1/watchlist?page=1&page_size=20' -H 'X-User-Fingerprint: demo-fp' | jq '.data.pagination,.data.items[0]'`
3. Feed:
   - `curl -sf 'http://localhost:3000/api/v1/watchlist/feed?page=1&page_size=20' -H 'X-User-Fingerprint: demo-fp' | jq '.data.pagination,.data.items[0]'`
4. Remove:
   - `curl -sf -X DELETE 'http://localhost:3000/api/v1/watchlist/11443' -H 'X-User-Fingerprint: demo-fp'`

## 6. Enable Amazon Nova
1. Set env in `.env`:
   - `NOVA_ENABLED=true`
   - `NOVA_PROVIDER=devapi`
   - `NOVA_API_BASE_URL=https://api.nova.amazon.com/v1`
   - `NOVA_API_KEY=<your_secret_key>`
   - `NOVA_ANALYSIS_MODEL=nova-pro-v1`
2. Optional (Bedrock mode):
   - `NOVA_PROVIDER=bedrock`
   - `NOVA_REGION=us-east-1`
   - `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN` (if using STS)
3. Restart backend: `docker compose restart backend`
4. Trigger analysis: `curl -X POST http://localhost:8080/api/v1/ai/analyze/<wallet_id>`
5. Read report: `curl -sf http://localhost:8080/api/v1/ai/report/<wallet_id>`

### Batch AI Analyzer (recommended)
- The backend now runs `ai_batch_analyzer` automatically and only analyzes wallets that satisfy:
  - total trades `>= 100`
  - realized PnL `> 0`
- Key envs:
  - `WORKER_AI_BATCH_ENABLED=true`
  - `WORKER_AI_BATCH_ANALYZER_INTERVAL=10m`
  - `WORKER_AI_BATCH_SIZE=3`
  - `WORKER_AI_BATCH_COOLDOWN=24h`
  - `WORKER_AI_BATCH_REQUEST_SPACING=25s`
  - `WORKER_AI_BATCH_MIN_TRADES=100`
  - `WORKER_AI_BATCH_MIN_REALIZED_PNL=0`
- Inspect progress: `docker compose logs -f backend | rg ai_batch_analyzer`

## 7. Backup & Restore
- Backup: `./ops/backup/backup.sh`
- Restore: `./ops/backup/restore.sh <backup.sql.gz>`

## 8. Incident Checklist
1. Confirm DB connectivity and migration status.
2. Confirm API health endpoint.
3. Inspect recent deploy/CI artifacts.
4. Check worker logs for ingestion failures.
5. Disable workers (`WORKER_ENABLED=false`) if external APIs are unstable.

## 9. Troubleshooting
1. `watchlist/feed` empty:
   - Confirm user fingerprint is consistent in every request (`X-User-Fingerprint`).
   - Confirm watchlist exists: `SELECT * FROM watchlist WHERE user_fingerprint='demo-fp';`
2. AI analyze returns `wallet has fewer than 100 trades`:
   - Use `/wallets/potential` and pick wallets with `total_trades >= 100` and positive PnL.
3. Backend log has `relation "wallet_update_event" does not exist`:
   - Apply migration `004_watchlist.sql` to the active DB:
     - `docker compose exec -T postgres psql -U postgres -d easy_arbitra -f /dev/stdin < backend/migrations/004_watchlist.sql`
4. Frontend sees stale data:
   - Force refresh and ensure API via nginx endpoint (`http://localhost:3000/api/v1/...`), not direct `backend:8080` in browser.
