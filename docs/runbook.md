# Easy Arbitra Runbook

## 1. Local Start
1. Copy env: `cp .env.example .env`
2. Start stack: `docker compose up -d`
3. Check services: `docker compose ps`
4. Backend health: `curl -sf http://localhost:8080/healthz`
5. Backend ready: `curl -sf http://localhost:8080/readyz`

## 2. Seed Data
1. Apply migrations: `cd backend && make migrate`
2. Run seed: `cd backend && make seed`

## 3. One-Command Bootstrap
- Run `./scripts/bootstrap_dev.sh` to start containers, migrate DB, seed data, and execute smoke tests.

## 4. Common Operations
- Restart backend only: `docker compose restart backend`
- View logs: `docker compose logs -f backend`
- Rebuild frontend: `docker compose build frontend && docker compose up -d frontend`

## 5. Backup & Restore
- Backup: `./ops/backup/backup.sh`
- Restore: `./ops/backup/restore.sh <backup.sql.gz>`

## 6. Incident Checklist
1. Confirm DB connectivity and migration status.
2. Confirm API health endpoint.
3. Inspect recent deploy/CI artifacts.
4. Check worker logs for ingestion failures.
5. Disable workers (`WORKER_ENABLED=false`) if external APIs are unstable.
