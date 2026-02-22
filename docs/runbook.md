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

## 5. Enable Amazon Nova
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

## 6. Backup & Restore
- Backup: `./ops/backup/backup.sh`
- Restore: `./ops/backup/restore.sh <backup.sql.gz>`

## 7. Incident Checklist
1. Confirm DB connectivity and migration status.
2. Confirm API health endpoint.
3. Inspect recent deploy/CI artifacts.
4. Check worker logs for ingestion failures.
5. Disable workers (`WORKER_ENABLED=false`) if external APIs are unstable.
