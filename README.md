# Easy Arbitra

Easy Arbitra is an AI-assisted Polymarket wallet analyzer focused on NBA trading behavior. A user submits a wallet address or Polymarket profile URL, the system resolves the identity, fetches historical NBA trades, computes deterministic style metrics, and asks an OpenAI-compatible model to explain the trader's behavior in plain English.

## Architecture

The project is split into two services:

- `frontend/`: a Next.js 16 application deployed to Cloudflare Workers
- `backend/`: a Go service that exposes MCP tools and a REST bridge for Polymarket data access

High-level request flow:

1. The user submits a wallet or profile URL in the frontend.
2. The frontend calls its own `/api/analyze` route.
3. The route runs a fixed deterministic 4-step MCP pipeline.
4. The backend resolves the wallet, fetches NBA trades, computes metrics, and returns a structured report payload.
5. The frontend sends the structured result to an OpenAI-compatible chat completions API for the final explanation.
6. The frontend renders the decision log, wallet card, radar chart, and natural-language explanation.

## Feature List

- Wallet input supports both raw Ethereum addresses and Polymarket profile URLs
- Polymarket public profile lookup for display name and avatar enrichment
- NBA-only trade filtering using Polymarket sports tags, events, and market metadata
- Deterministic trading metrics:
  - entry timing
  - size ratio versus market volume
  - conviction based on average buy price
- LLM-generated style explanation powered by an OpenAI-compatible API
- Structured UI output including a decision log, summary card, and radar chart
- Cloudflare Worker deployment for the frontend
- GitHub Container Registry build pipeline for the backend image

## Repository Layout

```text
.
├── backend/              Go MCP server and REST bridge
├── frontend/             Next.js app and Cloudflare Worker config
└── .github/workflows/    CI/CD for frontend deploy and backend image builds
```

## Core Backend Tool Sequence

The analysis route enforces a four-tool pipeline:

1. `resolve_wallet_target`
2. `fetch_sports_trades`
3. `calculate_style_metrics`
4. `build_report_payload`

This keeps the market data and metric calculations deterministic while using the LLM only for the final explanation layer.

## Local Development

Backend:

```bash
cd backend
go run .
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

The frontend expects these runtime values:

- `AI_BASE_URL`
- `AI_MODEL`
- `AI_API_KEY`
- `AI_TIMEOUT_MS` optional, defaults to `120000`
- `MCP_BRIDGE_URL`
- `DATABASE_URL` on the backend for Postgres-backed wallet catalog sync

The frontend calls an OpenAI-compatible `chat/completions` endpoint for the final explanation. If the AI provider is unavailable, the app still completes the deterministic 4-step tool pipeline and falls back to a locally generated narrative.

For local development, `MCP_BRIDGE_URL` should usually point to `http://localhost:8082`.

## Deployment

- Frontend: GitHub Actions builds and deploys the Next.js app to Cloudflare Workers
- Backend: GitHub Actions builds and publishes a Docker image to `ghcr.io`
- Runtime topology: Cloudflare Worker calls the backend running on your EC2 instance
- If `DATABASE_URL` is set on the backend, it will sync the top 100 NBA wallets from Polymarket Analytics every 4 hours and store AI style tags for homepage grouping

## Tech Stack

- Next.js 16
- React 19
- Tailwind CSS 4
- Go 1.23
- OpenAI-compatible chat completions API
- Cloudflare Workers
- GitHub Actions
- GitHub Container Registry
