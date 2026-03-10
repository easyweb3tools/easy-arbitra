import {
  BedrockRuntimeClient,
  ConverseCommand,
  type Message,
  type ContentBlock,
  type ToolResultContentBlock,
} from "@aws-sdk/client-bedrock-runtime";
import { callMCPTool } from "./mcp-bridge";
import { toolConfig } from "./tool-definitions";
import type { DecisionStep, ReportPayload, AnalyzeResponse } from "./types";

const client = new BedrockRuntimeClient({
  region: process.env.AWS_REGION || "us-east-1",
});

const MODEL_ID = process.env.BEDROCK_MODEL_ID || "us.amazon.nova-lite-v1:0";

const SYSTEM_PROMPT = `You are SportStyle AI Explainer, an expert sports betting analyst. Your task is to analyze a Polymarket wallet's NBA trading style.

You MUST follow this exact 4-step tool calling sequence:
1. Call resolve_wallet_target with the user's input to standardize the wallet address
2. Call fetch_sports_trades with the resolved wallet address to get NBA trades
3. Call calculate_style_metrics with the wallet and trades JSON to compute style metrics
4. Call build_report_payload with wallet info and metrics to generate the final report

After completing all 4 tool calls, provide a natural language explanation of the wallet's trading style. Your explanation should:
- Reference specific metric values (entry timing hours, position size %, ROI %)
- Describe the trading pattern in plain English
- Assign a style characterization (e.g., "Early Whale", "Quick Scout", "Sharp Shooter")
- Be 2-3 paragraphs long

If the wallet has zero NBA trades, explain that no NBA trading activity was found and suggest the user try a different wallet.

IMPORTANT: Always call tools in the exact order above. Pass data between steps:
- Step 2 uses wallet_address from Step 1
- Step 3 uses wallet from Step 1 and trades array from Step 2
- Step 4 uses wallet_info (full JSON from Step 1), metrics_json (full JSON from Step 3)`;

export async function analyzeWallet(
  walletInput: string
): Promise<AnalyzeResponse> {
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

    // Add assistant message to conversation
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

    // If no tool calls, we're done
    if (toolUseBlocks.length === 0) {
      break;
    }

    // Process each tool call
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

      decisionLog.push({
        step: stepNum,
        tool: name!,
        reasoning: getToolReasoning(name!),
        timestamp,
        result_summary: resultSummary,
      });

      // Try to extract report payload from build_report_payload result
      if (name === "build_report_payload") {
        try {
          reportPayload = JSON.parse(resultText) as ReportPayload;
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

    // Add tool results as user message
    messages.push({
      role: "user",
      content: toolResults,
    });

    // Check stop reason
    if (response.stopReason === "end_turn" && hasText) {
      break;
    }
  }

  // Fallback report payload if none was extracted
  if (!reportPayload) {
    reportPayload = {
      wallet_card: {
        address: walletInput,
        display_name: walletInput.substring(0, 10) + "...",
        profile_image: "",
        sport: "NBA",
        total_trades: 0,
      },
      radar_chart: { entry_timing: 0, size_ratio: 0, roi: 0 },
      report: {
        style_label: "Unknown",
        summary_context: "No data available",
      },
    };
  }

  return {
    decisionLog,
    reportPayload,
    explanation: explanation || "Analysis complete. See the dashboard for results.",
  };
}

function getToolReasoning(toolName: string): string {
  switch (toolName) {
    case "resolve_wallet_target":
      return "Standardizing the wallet input to a verified address with profile info";
    case "fetch_sports_trades":
      return "Fetching NBA-specific trade history from Polymarket";
    case "calculate_style_metrics":
      return "Computing entry timing, position sizing, and ROI metrics";
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
        return `Metrics: timing=${data.metrics?.entry_timing_hours}h, size=${data.metrics?.size_ratio_pct}%, ROI=${data.metrics?.roi}%`;
      case "build_report_payload":
        return `Style: ${data.report?.style_label}`;
      default:
        return "Done";
    }
  } catch {
    return "Result received";
  }
}
