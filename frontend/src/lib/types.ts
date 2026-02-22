export type Wallet = {
  id: number;
  address: string;
  pseudonym?: string;
  tracked: boolean;
};

export type PotentialWallet = {
  wallet: Wallet;
  total_trades: number;
  trading_pnl: number;
  maker_rebates: number;
  realized_pnl: number;
  smart_score: number;
  info_edge_level: string;
  strategy_type: string;
  has_ai_report: boolean;
  last_analyzed_at?: string;
};

export type Market = {
  id: number;
  condition_id: string;
  slug: string;
  title: string;
  category: string;
  status: number;
  has_fee: boolean;
  volume: number;
  liquidity: number;
};

export type LeaderboardItem = {
  wallet_id: number;
  address: string;
  pseudonym?: string;
  strategy_type: string;
  smart_score: number;
  info_edge_level: string;
  scored_at: string;
};

export type AnomalyAlert = {
  id: number;
  wallet_id: number;
  market_id?: number;
  alert_type: string;
  severity: number;
  evidence: Record<string, unknown>;
  description: string;
  acknowledged: boolean;
  created_at: string;
};

export type AIReport = {
  id: number;
  wallet_id: number;
  model_id: string;
  report: unknown;
  nl_summary: string;
  input_tokens: number;
  output_tokens: number;
  latency_ms: number;
  created_at: string;
};

export type WalletProfile = {
  wallet: Wallet;
  layer1_facts: {
    realized_pnl: number;
    unrealized_pnl: number;
    trading_pnl: number;
    maker_rebates: number;
    fees_paid: number;
    total_trades: number;
    volume_30d: number;
  };
  strategy?: {
    strategy_type: string;
    smart_score: number;
    info_edge_level: string;
    confidence: number;
    scored_at: string;
  };
  layer3_info_edge: {
    mean_delta_minutes: number;
    stddev_minutes: number;
    samples: number;
    p_value: number;
    label: string;
  };
  meta: {
    disclosures: string[];
  };
};

export type WalletShareCard = {
  wallet: Wallet;
  total_trades: number;
  realized_pnl: number;
  smart_score: number;
  info_edge_level: string;
  strategy_type: string;
  has_ai_report: boolean;
  nl_summary: string;
  updated_at: string;
};

export type WalletExplanation = {
  wallet_id: number;
  address: string;
  layer1: Record<string, unknown>;
  layer2: Record<string, unknown>;
  layer3: Record<string, unknown>;
  disclosures: string[];
  generated_at: string;
};

export type OverviewStats = {
  tracked_wallets: number;
  indexed_markets: number;
};

export type OpsTopRealizedWallet = {
  wallet: Wallet;
  trade_count: number;
  realized_pnl: number;
  realized_pnl_24h: number;
  has_ai_report: boolean;
  nl_summary: string;
  model_id: string;
  last_analyzed_at?: string;
};

export type OpsTopAIConfidenceWallet = {
  wallet: Wallet;
  trade_count: number;
  realized_pnl: number;
  smart_score: number;
  info_edge_level: string;
  strategy_type: string;
  nl_summary: string;
  last_analyzed_at?: string;
};

export type OpsHighlights = {
  as_of: string;
  new_potential_wallets_24h: number;
  top_realized_pnl_24h: OpsTopRealizedWallet[];
  top_ai_confidence: OpsTopAIConfidenceWallet[];
};

export type Paged<T> = {
  items: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
  };
};

export type ApiEnvelope<T> = {
  data: T;
  error?: string;
};
