import {
  BedrockRuntimeClient,
  ConverseCommand,
  type Message,
  type ContentBlock,
  type ToolResultContentBlock,
} from "@aws-sdk/client-bedrock-runtime";
import { FetchHttpHandler } from "@smithy/fetch-http-handler";
import { callMCPTool } from "./mcp-bridge";
import { toolConfig } from "./tool-definitions";
import type { DecisionStep, ReportPayload, AnalyzeResponse } from "./types";

const client = new BedrockRuntimeClient({
  region: process.env.AWS_REGION || "us-east-1",
  requestHandler: new FetchHttpHandler(),
});

const MODEL_ID = process.env.BEDROCK_MODEL_ID || "us.amazon.nova-lite-v1:0";

const SYSTEM_PROMPT = `You are SportStyle AI Explainer, an expert sports betting analyst. Your task is to analyze a Polymarket wallet's NBA trading style.

You MUST follow this exact 4-step tool calling sequence:
1. Call resolve_wallet_target with the user's input to standardize the wallet address
2. Call fetch_sports_trades with the resolved wallet address to get NBA trades
3. Call calculate_style_metrics with the wallet and trades JSON to compute style metrics
4. Call build_report_payload with wallet info and metrics to generate the final report

After completing all 4 tool calls, provide a natural language explanation of the wallet's trading style. Your explanation should:
- Reference specific metric values (entry timing hours, position size %, conviction score)
- Describe the trading pattern in plain English
- Conviction measures average buy price: high (>0.7) = bets on favorites, low (<0.4) = contrarian/underdog hunter
- Assign a style characterization (e.g., "Early Whale", "Quick Scout", "Favorite Backer", "Contrarian Hunter")
- Be 2-3 paragraphs long

If the wallet has zero NBA trades, explain that no NBA trading activity was found and suggest the user try a different wallet.

IMPORTANT: Always call tools in the exact order above. Pass data between steps:
- Step 2 uses wallet_address from Step 1
- Step 3 uses wallet from Step 1 and trades array from Step 2
- Step 4 uses wallet_info (full JSON from Step 1), metrics_json (full JSON from Step 3)`;

// SSE event types
export type StreamEvent =
  | { type: "step"; data: DecisionStep }
  | { type: "report"; data: ReportPayload }
  | { type: "explanation"; data: string }
  | { type: "done"; data: AnalyzeResponse }
  | { type: "error"; data: string };

export async function analyzeWalletStream(
  walletInput: string,
  onEvent: (event: StreamEvent) => void
): Promise<void> {
  const decisionLog: DecisionStep[] = [];
  let reportPayload: ReportPayload | null = null;
  let explanation = "";

  const messages: Message[] = [
    {
      role: "user",
      content: [
        {
          text: `Analyze the NBA trading style for this Polymarket wallet: ${walletInput}`,
        },
      ],
    },
  ];

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const bedrockTools = toolConfig.tools.map((t) => ({
    toolSpec: {
      name: t.toolSpec.name,
      description: t.toolSpec.description,
      inputSchema: { json: t.toolSpec.inputSchema.json },
    },
  })) as any[];

  let stepCount = 0;
  const MAX_ITERATIONS = 10;

  while (stepCount < MAX_ITERATIONS) {
    stepCount++;

    const command = new ConverseCommand({
      modelId: MODEL_ID,
      system: [{ text: SYSTEM_PROMPT }],
      messages,
      toolConfig: { tools: bedrockTools },
    });

    const response = await client.send(command);
    const output = response.output;

    if (!output?.message?.content) {
      break;
    }

    messages.push(output.message);

    const toolUseBlocks: ContentBlock[] = [];
    let hasText = false;

    for (const block of output.message.content) {
      if (block.text) {
        hasText = true;
        explanation += block.text;
      }
      if (block.toolUse) {
        toolUseBlocks.push(block);
      }
    }

    if (toolUseBlocks.length === 0) {
      break;
    }

    const toolResults: ContentBlock[] = [];

    for (const block of toolUseBlocks) {
      if (!block.toolUse) continue;

      const { toolUseId, name, input } = block.toolUse;
      const stepNum = decisionLog.length + 1;
      const timestamp = new Date().toISOString();

      let resultText: string;
      let resultSummary: string;

      try {
        resultText = await callMCPTool(
          name!,
          (input as Record<string, unknown>) || {}
        );
        resultSummary = summarizeResult(name!, resultText);
      } catch (error) {
        resultText = JSON.stringify({
          error: error instanceof Error ? error.message : "Unknown error",
        });
        resultSummary = `Error: ${error instanceof Error ? error.message : "Unknown error"}`;
      }

      const step: DecisionStep = {
        step: stepNum,
        tool: name!,
        reasoning: getToolReasoning(name!),
        timestamp,
        result_summary: resultSummary,
      };
      decisionLog.push(step);

      // Stream the step immediately
      onEvent({ type: "step", data: step });

      // Extract report payload
      if (name === "build_report_payload") {
        try {
          reportPayload = JSON.parse(resultText) as ReportPayload;
          onEvent({ type: "report", data: reportPayload });
        } catch {
          // Not valid JSON, skip
        }
      }

      const toolResultContent: ToolResultContentBlock[] = [{ text: resultText }];
      toolResults.push({
        toolResult: {
          toolUseId: toolUseId!,
          content: toolResultContent,
        },
      });
    }

    messages.push({
      role: "user",
      content: toolResults,
    });

    if (response.stopReason === "end_turn" && hasText) {
      break;
    }
  }

  // Stream the explanation
  if (explanation) {
    onEvent({ type: "explanation", data: explanation });
  }

  // Fallback report payload
  if (!reportPayload) {
    reportPayload = {
      wallet_card: {
        address: walletInput,
        display_name: walletInput.substring(0, 10) + "...",
        profile_image: "",
        sport: "NBA",
        total_trades: 0,
      },
      radar_chart: { entry_timing: 0, size_ratio: 0, conviction: 0 },
      report: {
        style_label: "Unknown",
        summary_context: "No data available",
      },
    };
    onEvent({ type: "report", data: reportPayload });
  }

  // Final done event with complete data
  onEvent({
    type: "done",
    data: {
      decisionLog,
      reportPayload,
      explanation: explanation || "Analysis complete. See the dashboard for results.",
    },
  });
}

function getToolReasoning(toolName: string): string {
  switch (toolName) {
    case "resolve_wallet_target":
      return "Standardizing the wallet input to a verified address with profile info";
    case "fetch_sports_trades":
      return "Fetching NBA-specific trade history from Polymarket";
    case "calculate_style_metrics":
      return "Computing entry timing, position sizing, and conviction metrics";
    case "build_report_payload":
      return "Building structured report with radar chart data and style classification";
    default:
      return "Processing...";
  }
}

function summarizeResult(toolName: string, resultText: string): string {
  try {
    const data = JSON.parse(resultText);
    switch (toolName) {
      case "resolve_wallet_target":
        return `Resolved: ${data.display_name} (${data.input_type})`;
      case "fetch_sports_trades":
        return `Found ${data.total_trades} NBA trades`;
      case "calculate_style_metrics":
        return `Metrics: timing=${data.metrics?.entry_timing_hours}h, size=${data.metrics?.size_ratio_pct}%, conviction=${data.metrics?.conviction}`;
      case "build_report_payload":
        return `Style: ${data.report?.style_label}`;
      default:
        return "Done";
    }
  } catch {
    return "Result received";
  }
}
