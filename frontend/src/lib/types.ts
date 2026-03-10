// TypeScript types for SportStyle AI Explainer

export interface WalletCard {
  address: string;
  display_name: string;
  profile_image: string;
  sport: string;
  total_trades: number;
}

export interface RadarChartData {
  entry_timing: number;
  size_ratio: number;
  roi: number;
}

export interface ReportData {
  style_label: string;
  summary_context: string;
}

export interface ReportPayload {
  wallet_card: WalletCard;
  radar_chart: RadarChartData;
  report: ReportData;
}

export interface DecisionStep {
  step: number;
  tool: string;
  reasoning: string;
  timestamp: string;
  result_summary?: string;
}

export interface AnalyzeResponse {
  decisionLog: DecisionStep[];
  reportPayload: ReportPayload;
  explanation: string;
}

export interface ToolCallRequest {
  tool: string;
  args: Record<string, unknown>;
}

export interface MCPToolResult {
  content: Array<{ type: string; text: string }>;
  isError?: boolean;
}

// Bedrock Converse types
export interface BedrockMessage {
  role: "user" | "assistant";
  content: BedrockContent[];
}

export interface BedrockContent {
  text?: string;
  toolUse?: {
    toolUseId: string;
    name: string;
    input: Record<string, unknown>;
  };
  toolResult?: {
    toolUseId: string;
    content: Array<{ text: string }>;
    status?: "success" | "error";
  };
}
