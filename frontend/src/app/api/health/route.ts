import { NextResponse } from "next/server";

const MCP_BRIDGE_URL = process.env.MCP_BRIDGE_URL || "http://localhost:8082";

function getBridgeHealthEndpoints(baseUrl: string): string[] {
  const normalized = baseUrl.replace(/\/+$/, "");

  if (normalized.endsWith("/api/health")) {
    return [normalized];
  }

  if (normalized.endsWith("/api")) {
    return [`${normalized}/health`, `${normalized}/api/health`];
  }

  return [`${normalized}/api/health`, `${normalized}/health`];
}

export async function GET() {
  const endpoints = getBridgeHealthEndpoints(MCP_BRIDGE_URL);
  const checks: Array<{
    url: string;
    ok: boolean;
    status: number | null;
    body: string;
  }> = [];

  for (const url of endpoints) {
    try {
      const response = await fetch(url, {
        method: "GET",
        headers: { Accept: "application/json" },
        cache: "no-store",
      });

      const body = await response.text();
      checks.push({
        url,
        ok: response.ok,
        status: response.status,
        body,
      });

      if (response.ok) {
        return NextResponse.json({
          ok: true,
          service: "frontend",
          mcpBridgeUrl: MCP_BRIDGE_URL,
          checks,
        });
      }
    } catch (error) {
      checks.push({
        url,
        ok: false,
        status: null,
        body: error instanceof Error ? error.message : "Unknown error",
      });
    }
  }

  return NextResponse.json(
    {
      ok: false,
      service: "frontend",
      mcpBridgeUrl: MCP_BRIDGE_URL,
      checks,
    },
    { status: 502 }
  );
}
