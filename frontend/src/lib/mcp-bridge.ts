import type { ToolCallRequest, MCPToolResult } from "./types";

const MCP_BRIDGE_URL = process.env.MCP_BRIDGE_URL || "http://localhost:8082";

export async function callMCPTool(
  toolName: string,
  args: Record<string, unknown>
): Promise<string> {
  const request: ToolCallRequest = { tool: toolName, args };

  const response = await fetch(`${MCP_BRIDGE_URL}/api/tools/call`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error(`MCP bridge error: ${response.status} ${await response.text()}`);
  }

  const result: MCPToolResult = await response.json();

  if (result.isError) {
    const errorText = result.content?.[0]?.text || "Unknown tool error";
    throw new Error(errorText);
  }

  return result.content?.[0]?.text || "";
}
