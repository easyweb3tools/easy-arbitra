# Backend

The backend is a Go service that wraps Polymarket data access and deterministic trading metrics behind MCP tools and a simple REST bridge. It is designed to keep the market data pipeline explicit, testable, and independent from the LLM.

## Responsibilities

- Resolve wallet addresses and Polymarket profile URLs
- Fetch public profile metadata from Polymarket
- Discover NBA markets and filter a wallet's trade history to sports-specific trades
- Compute deterministic style metrics from enriched trade data
- Build a structured payload for the frontend dashboard
- Expose the tool chain through MCP and HTTP
- Sync NBA leaderboard wallets into Postgres on a schedule
- Generate style tags for homepage grouping

## Service Endpoints

- `:8081` SSE MCP server
- `:8082` REST bridge
- `GET /api/health` health check
- `POST /api/tools/call` tool invocation endpoint
- `GET /api/style-wallets` homepage style-group feed
- `POST /api/style-wallets/sync` manual sync trigger

## Tool Pipeline

The frontend calls Bedrock, and Bedrock is instructed to invoke these backend tools in order:

1. `resolve_wallet_target`
2. `fetch_sports_trades`
3. `calculate_style_metrics`
4. `build_report_payload`

## Metrics Produced

- `entry_timing_hours`: average time between market start and trade execution
- `size_ratio_pct`: average trade size relative to market volume
- `conviction`: average buy price on a 0 to 1 scale

These metrics are then normalized into the frontend radar chart and combined into a style label such as `Early Whale` or `Contrarian Hunter`.

## Package Layout

```text
backend/
├── main.go           Service bootstrap and HTTP wiring
├── polymarket/       Polymarket API client and data models
├── metrics/          Deterministic metric calculation
└── tools/            MCP tool handlers and report builder
```

## Local Development

Key environment variables:

- `DATABASE_URL` enables Postgres persistence and scheduled syncing
- `AI_BASE_URL`
- `AI_MODEL`
- `AI_API_KEY`
- `AI_TIMEOUT_MS` optional, used for batch style tagging
- `LEADERBOARD_SYNC_INTERVAL` optional, defaults to `4h`
- `LEADERBOARD_TOP_LIMIT` optional, defaults to `100`
- `WALLET_ANALYSIS_LIMIT` optional, defaults to `3000`

```bash
go run .
```

Build:

```bash
go build ./...
```

## Discover Demo Wallets

Use the wallet discovery script to find or rank candidate wallets for demos:

```bash
go run ./cmd/discover-wallets -recent-limit 400 -recent-pages 4 -output 10
```

If recent global trades do not include enough NBA activity, score your own curated wallet list instead:

```bash
go run ./cmd/discover-wallets -wallets-file ./wallets.txt -output 10
```

The script ranks wallets by presentation quality rather than profit. It favors larger NBA samples, broader market coverage, and more legible style metrics for demos.

## Container Build

The backend image is built from `backend/Dockerfile`.

Local build:

```bash
docker build -f Dockerfile .
```

GitHub Actions publishes the image to GitHub Container Registry through `.github/workflows/backend-ghcr.yml`.
