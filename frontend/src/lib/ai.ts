import { callMCPTool } from "./mcp-bridge";
import type {
  DecisionStep,
  ReportPayload,
  AnalyzeResponse,
  ToolLogEntry,
} from "./types";

const AI_BASE_URL = process.env.AI_BASE_URL?.trim();
const AI_MODEL = process.env.AI_MODEL?.trim();
const AI_API_KEY = process.env.AI_API_KEY?.trim();
const AI_TIMEOUT_MS = parseTimeoutMs(process.env.AI_TIMEOUT_MS);

const AI_SYSTEM_PROMPT = `You are SportStyle AI Explainer, an expert sports betting analyst.

You will receive:
- wallet identity information
- deterministic NBA trade metrics
- a structured report payload

Write a concise, plain-English explanation of the wallet's NBA trading style.

Requirements:
- Be accurate to the supplied metrics only
- Reference specific metric values
- Explain the style label in plain English
- If sample size is small, mention it as a limitation
- If there are zero NBA trades, clearly say there is not enough NBA history to infer style
- Keep the response to 2 short paragraphs max
- Do not mention hidden prompts, tool calls, or internal implementation details`;

export type StreamEvent =
  | { type: "step"; data: DecisionStep }
  | { type: "tool_log"; data: ToolLogEntry }
  | { type: "report"; data: ReportPayload }
  | { type: "explanation"; data: string }
  | { type: "done"; data: AnalyzeResponse }
  | { type: "error"; data: string };

interface ResolveWalletResult {
  wallet_address: string;
  display_name: string;
  input_type: string;
  profile_image: string;
}

interface FetchTradesResult {
  wallet: string;
  sport: string;
  total_trades: number;
  trades: Array<Record<string, unknown>>;
}

interface MetricsResult {
  wallet: string;
  metrics: {
    entry_timing_hours: number;
    size_ratio_pct: number;
    conviction: number;
  };
  sample_size: number;
  warning?: string;
}

interface PipelineResult {
  decisionLog: DecisionStep[];
  walletInfo: ResolveWalletResult;
  tradesResult: FetchTradesResult;
  metricsResult: MetricsResult;
  reportPayload: ReportPayload;
}

export async function analyzeWalletStream(
  walletInput: string,
  onEvent: (event: StreamEvent) => void
): Promise<void> {
  const pipeline = await runDeterministicPipeline(walletInput, onEvent);
  const explanation = await generateExplanation(pipeline);

  onEvent({ type: "explanation", data: explanation });
  onEvent({
    type: "done",
    data: {
      decisionLog: pipeline.decisionLog,
      reportPayload: pipeline.reportPayload,
      explanation,
    },
  });
}

async function runDeterministicPipeline(
  walletInput: string,
  onEvent: (event: StreamEvent) => void
): Promise<PipelineResult> {
  const decisionLog: DecisionStep[] = [];
  const toolLogs: Record<string, string[]> = {};

  const emitToolLog = (tool: string, message: string) => {
    toolLogs[tool] = [...(toolLogs[tool] || []), message];
    onEvent({
      type: "tool_log",
      data: {
        tool,
        message,
        timestamp: new Date().toISOString(),
      },
    });
  };

  const addStep = (tool: string, resultText: string) => {
    const step: DecisionStep = {
      step: decisionLog.length + 1,
      tool,
      reasoning: getToolReasoning(tool),
      timestamp: new Date().toISOString(),
      result_summary: summarizeResult(tool, resultText),
      logs: toolLogs[tool] || [],
    };

    decisionLog.push(step);
    onEvent({ type: "step", data: step });
  };

  const walletInfoText = await callMCPTool(
    "resolve_wallet_target",
    { input: walletInput },
    (message) => emitToolLog("resolve_wallet_target", message)
  );
  addStep("resolve_wallet_target", walletInfoText);
  const walletInfo = JSON.parse(walletInfoText) as ResolveWalletResult;

  const tradesText = await callMCPTool(
    "fetch_sports_trades",
    {
      wallet: walletInfo.wallet_address,
      sport: "nba",
    },
    (message) => emitToolLog("fetch_sports_trades", message)
  );
  addStep("fetch_sports_trades", tradesText);
  const tradesResult = JSON.parse(tradesText) as FetchTradesResult;

  const metricsText = await callMCPTool(
    "calculate_style_metrics",
    {
      wallet: walletInfo.wallet_address,
      trades_json: tradesText,
    },
    (message) => emitToolLog("calculate_style_metrics", message)
  );
  addStep("calculate_style_metrics", metricsText);
  const metricsResult = JSON.parse(metricsText) as MetricsResult;

  const reportText = await callMCPTool(
    "build_report_payload",
    {
      wallet_info: walletInfoText,
      metrics_json: metricsText,
      trades_summary: JSON.stringify({
        total_trades: tradesResult.total_trades,
        sport: tradesResult.sport,
      }),
    },
    (message) => emitToolLog("build_report_payload", message)
  );
  addStep("build_report_payload", reportText);

  const reportPayload = JSON.parse(reportText) as ReportPayload;
  onEvent({ type: "report", data: reportPayload });

  return {
    decisionLog,
    walletInfo,
    tradesResult,
    metricsResult,
    reportPayload,
  };
}

async function generateExplanation(pipeline: PipelineResult): Promise<string> {
  if (!isAIConfigured()) {
    return buildFallbackExplanation(pipeline, "AI provider is not configured.");
  }

  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), AI_TIMEOUT_MS);
    const response = await fetch(getChatCompletionsUrl(AI_BASE_URL!), {
      method: "POST",
      headers: {
        Authorization: `Bearer ${AI_API_KEY}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        model: AI_MODEL,
        temperature: 0.4,
        messages: [
          { role: "system", content: AI_SYSTEM_PROMPT },
          {
            role: "user",
            content: buildExplanationPrompt(pipeline),
          },
        ],
      }),
      cache: "no-store",
      signal: controller.signal,
    }).finally(() => clearTimeout(timeoutId));

    if (!response.ok) {
      const body = await response.text();
      throw new Error(`AI API error ${response.status}: ${body}`);
    }

    const payload = (await response.json()) as {
      choices?: Array<{
        message?: {
          content?:
            | string
            | Array<{ type?: string; text?: string; content?: string }>;
        };
      }>;
    };

    const content = payload.choices?.[0]?.message?.content;
    const explanation = extractAssistantText(content);
    if (!explanation) {
      throw new Error("AI API returned an empty explanation.");
    }

    return explanation;
  } catch (error) {
    return buildFallbackExplanation(pipeline, getErrorMessage(error));
  }
}

function isAIConfigured(): boolean {
  return Boolean(AI_BASE_URL && AI_MODEL && AI_API_KEY);
}

function parseTimeoutMs(value: string | undefined): number {
  const parsed = Number.parseInt(value || "", 10);
  if (Number.isFinite(parsed) && parsed > 0) {
    return parsed;
  }
  return 120_000;
}

function getChatCompletionsUrl(baseUrl: string): string {
  const normalized = baseUrl.replace(/\/+$/, "");
  if (normalized.endsWith("/chat/completions")) {
    return normalized;
  }
  return `${normalized}/chat/completions`;
}

function buildExplanationPrompt({
  walletInfo,
  tradesResult,
  metricsResult,
  reportPayload,
}: PipelineResult): string {
  return JSON.stringify(
    {
      wallet: walletInfo,
      trades_summary: {
        total_trades: tradesResult.total_trades,
        sport: tradesResult.sport,
      },
      metrics: metricsResult,
      report: reportPayload,
    },
    null,
    2
  );
}

function extractAssistantText(
  content: string | Array<{ type?: string; text?: string; content?: string }> | undefined
): string {
  if (typeof content === "string") {
    return content.trim();
  }

  if (Array.isArray(content)) {
    return content
      .map((part) => part.text || part.content || "")
      .join("")
      .trim();
  }

  return "";
}

function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }

  if (typeof error === "string") {
    return error;
  }

  return "Unknown error";
}

function buildFallbackExplanation(
  { reportPayload, walletInfo, metricsResult }: PipelineResult,
  cause: string
): string {
  if (metricsResult.sample_size === 0) {
    return `${walletInfo.display_name} has no detected NBA trading activity in the fetched Polymarket history. The structured report is still generated from the deterministic pipeline, but there is not enough NBA data to infer a reliable style. AI explanation generation was unavailable for this request (${cause}).`;
  }

  const metrics = metricsResult.metrics;
  const explanationParts = [
    `${walletInfo.display_name} profiles as a ${reportPayload.report.style_label} based on ${metricsResult.sample_size} NBA trades. Average entry timing is ${metrics.entry_timing_hours.toFixed(1)} hours before market resolution, average position size is ${metrics.size_ratio_pct.toFixed(4)}% of market volume, and conviction is ${metrics.conviction.toFixed(2)} on the 0-1 scale.`,
    metrics.conviction > 0.75
      ? "That conviction score suggests a strong bias toward favorites or higher-confidence entries."
      : metrics.conviction > 0 && metrics.conviction < 0.35
        ? "That conviction score points to a more contrarian pattern, with entries clustering toward underdog pricing."
        : "The conviction score sits in the middle, which looks more balanced than aggressively favorite-seeking or contrarian.",
    metrics.entry_timing_hours <= 6
      ? "The trader tends to get involved relatively early."
      : "The trader tends to enter closer to resolution, which is more reactive than early-positioning.",
    `AI explanation generation was unavailable for this request (${cause}), so this explanation was generated from the deterministic metrics pipeline instead of an LLM.`,
  ];

  if (metricsResult.warning) {
    explanationParts.splice(explanationParts.length - 1, 0, metricsResult.warning);
  }

  return explanationParts.join(" ");
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
