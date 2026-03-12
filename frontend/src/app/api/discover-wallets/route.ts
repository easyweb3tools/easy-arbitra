import { NextResponse } from "next/server";

const MCP_BRIDGE_URL = process.env.MCP_BRIDGE_URL || "http://localhost:8082";

export async function POST(request: Request) {
  try {
    const body = await request.json();
    const response = await fetch(`${MCP_BRIDGE_URL}/api/discover-wallets`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    const text = await response.text();
    if (!response.ok) {
      return NextResponse.json(
        { error: text || "Failed to discover wallets" },
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
