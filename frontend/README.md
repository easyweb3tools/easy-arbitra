# Frontend

This frontend is a Next.js 16 application that runs on Cloudflare Workers. It provides the user-facing wallet analysis flow, runs the deterministic analysis pipeline from the `/api/analyze` route, and renders the final trader profile dashboard.

## Responsibilities

- Collect wallet input from the user
- Validate Ethereum addresses and Polymarket profile URLs
- Call an OpenAI-compatible chat completions API from the server route
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
│   ├── api/analyze/route.ts   Server-side analysis orchestration
│   ├── dashboard/page.tsx     Analysis result screen
│   └── page.tsx               Landing and input flow
├── components/                UI building blocks and report views
└── lib/                       AI client, MCP bridge, types, tool definitions
```

## Runtime Configuration

The frontend expects these values at runtime:

- `AI_BASE_URL`
- `AI_MODEL`
- `AI_API_KEY`
- `MCP_BRIDGE_URL`

The app calls an OpenAI-compatible `chat/completions` endpoint for the final explanation. If the AI provider is unavailable, the `/api/analyze` route still completes the deterministic tool pipeline and falls back to a local explanation.

`MCP_BRIDGE_URL` should be the base URL of the Go backend, for example:

```text
http://localhost:8082
https://api.example.com
```

The frontend tolerates `MCP_BRIDGE_URL` values with or without a trailing `/api`, but the recommended format is the service root URL.

## Local Development

```bash
npm install
npm run dev
```

Open `http://localhost:3000`.

For deployment debugging, the frontend also exposes `GET /api/health`. It returns the Worker-visible `MCP_BRIDGE_URL` and the result of probing the backend health endpoint.

## Cloudflare Deployment

This repo is configured to deploy the frontend as a Cloudflare Worker using OpenNext.

Relevant files:

- `open-next.config.ts`
- `wrangler.jsonc`
- `.github/workflows/frontend-cloudflare.yml`

GitHub Actions is the source of truth for Worker configuration:

- GitHub Variables become Cloudflare Worker `vars`
- GitHub Secrets are pushed as Cloudflare Worker secrets during deployment
- Set `AI_API_KEY` as a GitHub Actions secret for production deploys

## Feature Summary

- Wallet input with client-side validation
- OpenAI-compatible explanation flow
- Tool-driven data pipeline instead of direct LLM hallucination
- Streaming-friendly architecture for progressive analysis updates
- Responsive dashboard for trader style breakdown
