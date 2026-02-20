#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: restore.sh <backup.sql.gz>"
  exit 1
fi

FILE=$1
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-easy_arbitra}
DB_USER=${DB_USER:-postgres}

zcat "$FILE" | PGPASSWORD=${DB_PASSWORD:-postgres} psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME"
echo "restore done"
