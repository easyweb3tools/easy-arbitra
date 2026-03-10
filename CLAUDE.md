# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

**Backend (Go MCP Server):**
```bash
cd backend
go build ./...          # compile check
go run main.go          # starts SSE on :8081 + REST bridge on :8082
```

**Frontend (Next.js):**
```bash
cd frontend
npm install             # first time only
npm run dev             # dev server on :3000
npm run build           # production build (also runs TypeScript checking)
```

Both servers must run simultaneously. The frontend calls the Go backend REST bridge at `localhost:8082`.

**Deployment:**
```bash
# Frontend → Cloudflare Workers
cd frontend
npm run cf:build        # OpenNext Cloudflare build
npm run cf:preview      # local preview via wrangler
npm run cf:deploy       # deploy to Cloudflare Workers

# Backend → Docker / GHCR
docker build -f backend/Dockerfile .
```

CI/CD runs on `release/v*` branch pushes via GitHub Actions (`.github/workflows/`).

## Environment

Frontend requires `frontend/.env.local` (see `.env.local.example`):
- `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` — Bedrock credentials
- `BEDROCK_MODEL_ID` — defaults to `us.amazon.nova-lite-v1:0`
- `MCP_BRIDGE_URL` — defaults to `http://localhost:8082`

Backend has no env config — Polymarket API endpoints are hardcoded (`gamma-api.polymarket.com`, `data-api.polymarket.com`).

## Architecture

**SportStyle AI Explainer** — user inputs a Polymarket wallet address, Amazon Nova orchestrates 4 MCP tools via Bedrock Converse API, outputs NBA trading style analysis (radar chart + natural language).

```
Browser → Next.js /api/analyze (route.ts, SSE stream)
  → Bedrock Converse loop (bedrock.ts) with toolConfig, max 10 iterations
  → Nova returns tool_use → MCP Bridge (mcp-bridge.ts) → Go REST :8082
  → Tool result returns to Nova → loops 4 times → each step streamed to browser
  → Nova returns end_turn + explanation
  → Frontend renders Dashboard
```

**Dual-port Go server**: `:8081` is standard MCP SSE (for direct MCP clients), `:8082` is a REST bridge (for Next.js, avoids SSE session complexity).

**Page flow**: Home (`/`) collects wallet input → stores in `sessionStorage` → immediately navigates to `/dashboard` → dashboard starts streaming analysis → shows progress steps → renders results on completion. Results persist in `sessionStorage` for refresh resilience.

### Tool Calling Sequence (strict order, enforced by system prompt)

1. **resolve_wallet_target** — parse 0x address or Polymarket URL, fetch profile
2. **fetch_sports_trades** — accepts `wallet` (required), `sport` (default "nba"), `limit` (default 500) → get sport tag → paginate events → conditionIDs → filter user trades → fallback: check unmatched trades against market questions for sport keywords → enrich with market metadata
3. **calculate_style_metrics** — compute entry_timing_hours, size_ratio_pct, conviction (avg buy price 0-1) from enriched trades
4. **build_report_payload** — normalize to 0-1 radar values, assign style label, assemble wallet card + report

Each step passes its output as input to the next via the Bedrock Converse tool loop.

### Key Patterns

- **mcp-go v0.32.0**: Use `request.GetArguments()` to access tool args (not `request.Params.Arguments` which is typed `any`)
- **AWS SDK Tool types**: The Bedrock `Tool` type is a tagged union — cast with `as any[]` when building toolConfig
- **Streaming**: `/api/analyze` uses SSE. Events: `step`, `report`, `explanation`, `done`, `error`. Dashboard consumes the stream to show real-time tool progress.
- **Data between pages**: Wallet input and analysis results pass via `sessionStorage`; dashboard redirects to `/` if neither exists
- **Polymarket API rate limiting**: 200ms sleep between market metadata batch requests, batches of 20 conditionIDs

### Backend Structure

- `polymarket/` — API client layer: `gamma.go` (Gamma API: profiles, sports, events, markets), `data.go` (Data API: trades), `types.go` (shared structs)
- `tools/` — MCP tool handlers, each returns `*mcp.CallToolResult`
- `metrics/calculator.go` — pure deterministic functions (no API calls)
- `main.go` — server setup, tool registration with JSON schemas, REST bridge dispatcher

### Frontend Structure

- `src/lib/bedrock.ts` — core orchestration: Bedrock Converse loop, streams events via callback
- `src/lib/mcp-bridge.ts` — HTTP POST to Go REST bridge at `/api/tools/call`
- `src/lib/tool-definitions.ts` — Bedrock toolConfig mirroring MCP tool schemas
- `src/app/page.tsx` — home page: wallet input, navigates to dashboard on submit
- `src/app/dashboard/page.tsx` — runs streaming analysis on load, then renders results
- `src/app/api/analyze/route.ts` — Next.js API route, SSE endpoint
- `src/components/` — `wallet-input`, `decision-log`, `wallet-card`, `radar-chart`, `report-summary`

## Style Label Classification (build_report.go)

The third metric axis is **conviction** (avg buy price), not ROI. ROI cannot be computed honestly without settlement data.

| Condition | Label |
|-----------|-------|
| Early entry + Large position | Early Whale |
| Early entry + Small position | Quick Scout |
| Late entry + Large position | Late Whale |
| High conviction (>0.75) | Favorite Backer |
| Low conviction (<0.35) | Contrarian Hunter |
| Large position (>70%) | Heavy Hitter |
| Early entry (>50%) | Early Bird |
| Default | Steady Player |
