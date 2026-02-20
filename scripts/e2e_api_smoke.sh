#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${BASE_URL:-http://localhost:8080}

curl -sf "$BASE_URL/healthz" >/dev/null
curl -sf "$BASE_URL/readyz" >/dev/null
curl -sf "$BASE_URL/api/v1/wallets?page=1&page_size=1" >/dev/null
curl -sf "$BASE_URL/api/v1/markets?page=1&page_size=1" >/dev/null
curl -sf "$BASE_URL/api/v1/stats/overview" >/dev/null

# Optional endpoints (may return empty but should not 5xx)
curl -sf "$BASE_URL/api/v1/leaderboard?page=1&page_size=1" >/dev/null || true
curl -sf "$BASE_URL/api/v1/anomalies?page=1&page_size=1" >/dev/null || true

echo "API smoke passed"
