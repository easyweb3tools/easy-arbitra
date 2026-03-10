# Frontend

This frontend is a Next.js 16 application that runs on Cloudflare Workers. It provides the user-facing wallet analysis flow, orchestrates Bedrock from the `/api/analyze` route, and renders the final trader profile dashboard.

## Responsibilities

- Collect wallet input from the user
- Validate Ethereum addresses and Polymarket profile URLs
- Call Amazon Bedrock Converse from the server route
- Forward tool calls to the backend MCP bridge
- Render the analysis result as a dashboard with:
  - decision log
  - wallet summary card
  - trading style radar chart
  - final natural-language explanation

## App Structure

```text
src/
├── app/
│   ├── api/analyze/route.ts   Server-side Bedrock orchestration
│   ├── dashboard/page.tsx     Analysis result screen
│   └── page.tsx               Landing and input flow
├── components/                UI building blocks and report views
└── lib/                       Bedrock client, MCP bridge, types, tool definitions
```

## Runtime Configuration

The frontend expects these values at runtime:

- `AWS_REGION`
- `BEDROCK_MODEL_ID`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `MCP_BRIDGE_URL`

`MCP_BRIDGE_URL` should be the base URL of the Go backend, for example:

```text
http://localhost:8082
https://api.example.com
```

## Local Development

```bash
npm install
npm run dev
```

Open `http://localhost:3000`.

## Cloudflare Deployment

This repo is configured to deploy the frontend as a Cloudflare Worker using OpenNext.

Relevant files:

- `open-next.config.ts`
- `wrangler.jsonc`
- `.github/workflows/frontend-cloudflare.yml`

GitHub Actions is the source of truth for Worker configuration:

- GitHub Variables become Cloudflare Worker `vars`
- GitHub Secrets are pushed as Cloudflare Worker secrets during deployment

## Feature Summary

- Wallet input with client-side validation
- Bedrock-powered explanation flow
- Tool-driven data pipeline instead of direct LLM hallucination
- Streaming-friendly architecture for progressive analysis updates
- Responsive dashboard for trader style breakdown
