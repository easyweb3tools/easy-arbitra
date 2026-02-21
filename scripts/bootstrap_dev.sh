#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT_DIR"

cp -f .env.example .env

docker compose up -d postgres backend frontend

for i in {1..60}; do
  if curl -sf http://localhost:8080/healthz >/dev/null; then
    break
  fi
  sleep 2
  if [[ $i -eq 60 ]]; then
    echo "backend healthz timeout"
    exit 1
  fi
done

docker compose exec -T backend sh -lc 'cd /app && DATABASE_HOST=postgres DATABASE_PORT=5432 DATABASE_USER=postgres DATABASE_PASSWORD=postgres DATABASE_DBNAME=easy_arbitra DATABASE_SSLMODE=disable /usr/local/go/bin/go run ./cmd/migrate'

curl -sf http://localhost:8080/readyz >/dev/null
./scripts/e2e_api_smoke.sh
./scripts/e2e_ui_smoke.sh

echo "Dev bootstrap complete"
