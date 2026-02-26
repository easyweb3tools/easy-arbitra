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
  pool_tier: "star" | "strategy" | "observation";
  has_ai_report: boolean;
  nl_summary: string;
  summary: string;
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
  pool_tier?: "star" | "strategy" | "observation";
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
  risk_warnings?: string[];
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
    pool_tier: "star" | "strategy" | "observation";
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
  recent_events?: WalletEvent[];
};

export type WalletShareCard = {
  wallet: Wallet;
  total_trades: number;
  realized_pnl: number;
  smart_score: number;
  info_edge_level: string;
  strategy_type: string;
  pool_tier: "star" | "strategy" | "observation";
  has_ai_report: boolean;
  nl_summary: string;
  follower_count: number;
  new_followers_7d: number;
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
  recommend_reason: string;
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
  recommend_reason: string;
  last_analyzed_at?: string;
};

export type OpsHighlights = {
  as_of: string;
  new_potential_wallets_24h: number;
  top_realized_pnl_24h: OpsTopRealizedWallet[];
  top_ai_confidence: OpsTopAIConfidenceWallet[];
};

export type WatchlistItem = {
  watchlist_id: number;
  watchlisted_at: string;
  wallet: Wallet;
  total_trades: number;
  trading_pnl: number;
  maker_rebates: number;
  realized_pnl: number;
  smart_score: number;
  info_edge_level: string;
  strategy_type: string;
  pool_tier: "star" | "strategy" | "observation";
  has_ai_report: boolean;
  nl_summary: string;
  last_analyzed_at?: string;
};

export type WatchlistFeedItem = {
  event_id: number;
  wallet: Wallet;
  event_type: string;
  event_payload: Record<string, unknown>;
  action_required: boolean;
  suggestion: string;
  suggestion_zh: string;
  event_time: string;
};

export type WalletDecisionCard = {
  wallet_id: number;
  pool_tier: "star" | "strategy" | "observation";
  suitable_for: "conservative" | "aggressive" | "event_driven" | string;
  risk_level: "low" | "medium" | "high" | string;
  suggested_position: string;
  momentum: "heating" | "stable" | "cooling" | string;
  status_7d: string;
  recommendation: string;
  recommendation_zh: string;
  disclaimer: string;
  disclaimer_zh: string;
  last_updated: string;
};

export type WalletEvent = {
  event_id: number;
  event_type: string;
  event_payload: Record<string, unknown>;
  action_required: boolean;
  suggestion: string;
  suggestion_zh: string;
  event_time: string;
};

export type WalletShareLanding = {
  wallet: Wallet;
  pool_tier: "star" | "strategy" | "observation";
  strategy_type: string;
  smart_score: number;
  pnl_7d: number;
  pnl_30d: number;
  max_drawdown_7d: number;
  stability_score: number;
  nl_summary: string;
  follower_count: number;
  new_followers_7d: number;
  decision_card: WalletDecisionCard;
  updated_at: string;
};

export type PortfolioItem = {
  id: number;
  name: string;
  name_zh: string;
  description: string;
  risk_level: "low" | "medium" | "high" | string;
  expected_return: string;
  max_drawdown: string;
  wallet_ids: number[];
  wallets: Wallet[];
};

export type WatchlistSummary = {
  followed_wallets: number;
  style_distribution: Record<string, number>;
  action_required: number;
  healthy_wallets: number;
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

export type PnLHistoryPoint = {
  date: string;
  pnl_7d: number;
  pnl_30d: number;
  pnl_90d: number;
  trade_count_30d: number;
  active_days_30d: number;
  avg_edge: number;
};

export type TradeHistoryItem = {
  id: number;
  block_time: string;
  market_title: string;
  market_slug: string;
  outcome: "Yes" | "No";
  action: "Buy" | "Sell";
  price: number;
  size: number;
  fee_paid: number;
  is_maker: boolean;
};

export type WalletPosition = {
  market_id: number;
  market_title: string;
  market_slug: string;
  category: string;
  net_size: number;
  avg_price: number;
  total_volume: number;
  trade_count: number;
  last_trade_at: string;
};

// ── Copy Trading Types ──

export type CopyTradingConfig = {
  id: number;
  wallet_id: number;
  wallet_address: string;
  wallet_pseudonym?: string;
  enabled: boolean;
  max_position_usdc: number;
  risk_preference: "conservative" | "moderate" | "aggressive";
  total_pnl: number;
  total_copies: number;
  open_positions: number;
  created_at: string;
};

export type CopyTradeDecision = {
  id: number;
  decision: "copy" | "skip";
  confidence: number;
  market_title: string;
  outcome: string;
  action: string;
  price: number;
  size_usdc: number;
  stop_loss_price?: number;
  reasoning: string;
  reasoning_en: string;
  risk_notes: string[];
  status: "pending" | "executed" | "stopped" | "expired";
  realized_pnl?: number;
  created_at: string;
};

export type CopyTradeDashboard = {
  total_pnl: number;
  win_rate: number;
  total_copies: number;
  total_skipped: number;
  open_positions: number;
  active_configs: number;
  configs: CopyTradingConfig[];
  recent_decisions: CopyTradeDecision[];
};

export type CopyTradePerformance = {
  total_pnl: number;
  win_rate: number;
  total_copies: number;
  daily_points: { date: string; pnl: number; copies: number }[];
};

// ── Copy Trade Monitoring Types ──

export type SyncerRun = {
  id: number;
  job_name: string;
  started_at: string;
  ended_at?: string;
  status: string;
  stats: Record<string, number>;
  error_text?: string;
};

export type HourlyStat = {
  hour: string;
  runs: number;
  wallets_checked: number;
  new_trades: number;
  decisions_copy: number;
  decisions_skip: number;
  errors: number;
};

export type CopyableWallet = {
  id: number;
  pseudonym?: string;
  address: string;
  smart_score: number;
  pool_tier: string;
  strategy_type: string;
  risk_level: string;
  momentum: string;
  pnl_30d: number;
  trade_count_30d: number;
};

export type CopyTradeMonitor = {
  enabled_configs: number;
  recent_runs: SyncerRun[];
  hourly_stats: HourlyStat[];
  copyable_wallets: CopyableWallet[];
};
