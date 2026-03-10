# Backend

The backend is a Go service that wraps Polymarket data access and deterministic trading metrics behind MCP tools and a simple REST bridge. It is designed to keep the market data pipeline explicit, testable, and independent from the LLM.

## Responsibilities

- Resolve wallet addresses and Polymarket profile URLs
- Fetch public profile metadata from Polymarket
- Discover NBA markets and filter a wallet's trade history to sports-specific trades
- Compute deterministic style metrics from enriched trade data
- Build a structured payload for the frontend dashboard
- Expose the tool chain through MCP and HTTP

## Service Endpoints

- `:8081` SSE MCP server
- `:8082` REST bridge
- `GET /api/health` health check
- `POST /api/tools/call` tool invocation endpoint

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

```bash
go run .
```

Build:

```bash
go build ./...
```

## Container Build

The backend image is built from `backend/Dockerfile`.

Local build:

```bash
docker build -f Dockerfile .
```

GitHub Actions publishes the image to GitHub Container Registry through `.github/workflows/backend-ghcr.yml`.
