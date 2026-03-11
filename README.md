# Easy Arbitra

Easy Arbitra is an AI-assisted Polymarket wallet analyzer focused on NBA trading behavior. A user submits a wallet address or Polymarket profile URL, the system resolves the identity, fetches historical NBA trades, computes deterministic style metrics, and asks Amazon Bedrock to explain the trader's behavior in plain English.

## Architecture

The project is split into two services:

- `frontend/`: a Next.js 16 application deployed to Cloudflare Workers
- `backend/`: a Go service that exposes MCP tools and a REST bridge for Polymarket data access

High-level request flow:

1. The user submits a wallet or profile URL in the frontend.
2. The frontend calls its own `/api/analyze` route.
3. The route uses Amazon Bedrock Converse with tool calling.
4. Bedrock invokes the backend MCP bridge in a fixed 4-step sequence.
5. The backend resolves the wallet, fetches NBA trades, computes metrics, and returns a structured report payload.
6. The frontend renders the decision log, wallet card, radar chart, and natural-language explanation.

## Feature List

- Wallet input supports both raw Ethereum addresses and Polymarket profile URLs
- Polymarket public profile lookup for display name and avatar enrichment
- NBA-only trade filtering using Polymarket sports tags, events, and market metadata
- Deterministic trading metrics:
  - entry timing
  - size ratio versus market volume
  - conviction based on average buy price
- LLM-generated style explanation powered by Amazon Bedrock
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

The Bedrock prompt enforces a four-tool pipeline:

1. `resolve_wallet_target`
2. `fetch_sports_trades`
3. `calculate_style_metrics`
4. `build_report_payload`

This keeps the LLM orchestration constrained while leaving the market data and metric calculations deterministic.

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

- `AWS_REGION`
- `BEDROCK_MODEL_ID`
- `AWS_BEARER_TOKEN_BEDROCK`
- `MCP_BRIDGE_URL`

The frontend now prefers Amazon Bedrock API key auth via `AWS_BEARER_TOKEN_BEDROCK`. If that token is present, the Bedrock client uses bearer auth instead of SigV4 credentials.

If the AWS account used by the frontend is not allowlisted for Amazon Bedrock yet, or the bearer token cannot invoke Bedrock, the app falls back to the deterministic 4-step tool pipeline and still returns a report. In that mode, only the final narrative is generated locally instead of by Bedrock.

For local development, `MCP_BRIDGE_URL` should usually point to `http://localhost:8082`.

## Deployment

- Frontend: GitHub Actions builds and deploys the Next.js app to Cloudflare Workers
- Backend: GitHub Actions builds and publishes a Docker image to `ghcr.io`
- Runtime topology: Cloudflare Worker calls the backend running on your EC2 instance

Bedrock note:

- Some AWS accounts must be manually allowlisted before Bedrock API access is granted.
- If production shows an AWS message about `bedrock-allowlisting`, submit the AWS support request for that account.
- Until AWS approves access, the frontend will use deterministic fallback analysis instead of failing the request.

## Tech Stack

- Next.js 16
- React 19
- Tailwind CSS 4
- Go 1.23
- Amazon Bedrock
- Cloudflare Workers
- GitHub Actions
- GitHub Container Registry
