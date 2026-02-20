#!/usr/bin/env bash
set -euo pipefail

OUT_DIR=${1:-./backups}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-easy_arbitra}
DB_USER=${DB_USER:-postgres}

mkdir -p "$OUT_DIR"
TS=$(date -u +"%Y%m%dT%H%M%SZ")
FILE="$OUT_DIR/${DB_NAME}_${TS}.sql.gz"

PGPASSWORD=${DB_PASSWORD:-postgres} pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" | gzip > "$FILE"
echo "backup written: $FILE"
