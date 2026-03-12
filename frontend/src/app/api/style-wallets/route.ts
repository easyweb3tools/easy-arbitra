import { NextResponse } from "next/server";

const MCP_BRIDGE_URL = process.env.MCP_BRIDGE_URL || "http://localhost:8082";

export async function GET(request: Request) {
  try {
    const url = new URL(request.url);
    const limitPerGroup = url.searchParams.get("limit_per_group") || "6";
    const response = await fetch(
      `${MCP_BRIDGE_URL}/api/style-wallets?limit_per_group=${encodeURIComponent(limitPerGroup)}`,
      {
        method: "GET",
        cache: "no-store",
      }
    );

    const text = await response.text();
    if (!response.ok) {
      return NextResponse.json(
        { error: text || "Failed to load style wallets" },
        { status: response.status }
      );
    }

    return new Response(text, {
      headers: { "Content-Type": "application/json" },
    });
  } catch (error) {
    return NextResponse.json(
      {
        error:
          error instanceof Error ? error.message : "Internal server error",
      },
      { status: 500 }
    );
  }
}
