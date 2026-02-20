#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${BASE_URL:-http://localhost:3000}

curl -sf "$BASE_URL/" | grep -qi "Easy Arbitra"
curl -sf "$BASE_URL/wallets" | grep -qi "Wallets"
curl -sf "$BASE_URL/markets" | grep -qi "Markets"
curl -sf "$BASE_URL/leaderboard" | grep -qi "Leaderboard"
curl -sf "$BASE_URL/anomalies" | grep -qi "Anomaly"

# optional page, should return 200
curl -sf "$BASE_URL/methodology" >/dev/null

echo "UI smoke passed"
