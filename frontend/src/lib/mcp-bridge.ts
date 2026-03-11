import type { ToolCallRequest, MCPToolResult } from "./types";

const MCP_BRIDGE_URL = process.env.MCP_BRIDGE_URL || "http://localhost:8082";

export async function callMCPTool(
  toolName: string,
  args: Record<string, unknown>
): Promise<string> {
  const request: ToolCallRequest = { tool: toolName, args };
  const endpoints = getMcpBridgeEndpoints(MCP_BRIDGE_URL);

  let lastError = "";
  for (const endpoint of endpoints) {
    const response = await fetch(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      lastError = `MCP bridge error: ${response.status} ${await response.text()}`;
      if (response.status === 404) {
        continue;
      }
      throw new Error(lastError);
    }

    const result: MCPToolResult = await response.json();

    if (result.isError) {
      const errorText = result.content?.[0]?.text || "Unknown tool error";
      throw new Error(errorText);
    }

    return result.content?.[0]?.text || "";
  }

  throw new Error(
    `${lastError} (checked endpoints: ${endpoints.join(", ")})`
  );
}

function getMcpBridgeEndpoints(baseUrl: string): string[] {
  const normalized = baseUrl.replace(/\/+$/, "");

  if (normalized.endsWith("/api/tools/call")) {
    return [normalized];
  }

  if (normalized.endsWith("/api")) {
    return [`${normalized}/tools/call`, `${normalized}/api/tools/call`];
  }

  return [`${normalized}/api/tools/call`, `${normalized}/tools/call`];
}
