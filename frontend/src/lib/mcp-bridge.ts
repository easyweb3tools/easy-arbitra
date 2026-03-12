import type { ToolCallRequest, MCPToolResult } from "./types";

const MCP_BRIDGE_URL = process.env.MCP_BRIDGE_URL || "http://localhost:8082";

export async function callMCPTool(
  toolName: string,
  args: Record<string, unknown>,
  onLog?: (message: string) => void
): Promise<string> {
  const request: ToolCallRequest = { tool: toolName, args };
  const endpoints = onLog
    ? getMcpBridgeStreamEndpoints(MCP_BRIDGE_URL)
    : getMcpBridgeEndpoints(MCP_BRIDGE_URL);

  let lastError = "";
  for (const endpoint of endpoints) {
    const response = await fetch(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(request),
    });

    if (onLog && response.ok) {
      return readStreamedToolResult(response, onLog);
    }

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

async function readStreamedToolResult(
  response: Response,
  onLog: (message: string) => void
): Promise<string> {
  const reader = response.body?.getReader();
  if (!reader) {
    throw new Error("MCP bridge stream missing response body");
  }

  const decoder = new TextDecoder();
  let buffer = "";
  let resultText = "";

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    buffer += decoder.decode(value, { stream: true });
    const lines = buffer.split("\n\n");
    buffer = lines.pop() || "";

    for (const line of lines) {
      const dataLine = line.replace(/^data: /, "").trim();
      if (!dataLine) continue;

      const event = JSON.parse(dataLine) as {
        type: string;
        message?: string;
        result?: MCPToolResult;
      };

      if (event.type === "log" && event.message) {
        onLog(event.message);
        continue;
      }

      if (event.type === "error") {
        throw new Error(event.message || "Unknown stream error");
      }

      if (event.type === "result" && event.result) {
        if (event.result.isError) {
          const errorText =
            event.result.content?.[0]?.text || "Unknown tool error";
          throw new Error(errorText);
        }

        resultText = event.result.content?.[0]?.text || "";
      }
    }
  }

  if (!resultText) {
    throw new Error("MCP bridge stream ended without a result");
  }

  return resultText;
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

function getMcpBridgeStreamEndpoints(baseUrl: string): string[] {
  const normalized = baseUrl.replace(/\/+$/, "");

  if (normalized.endsWith("/api/tools/call")) {
    return [`${normalized}-stream`];
  }

  if (normalized.endsWith("/api")) {
    return [`${normalized}/tools/call-stream`, `${normalized}/api/tools/call-stream`];
  }

  return [`${normalized}/api/tools/call-stream`, `${normalized}/tools/call-stream`];
}
