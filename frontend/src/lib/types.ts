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
};

export type OverviewStats = {
  tracked_wallets: number;
  indexed_markets: number;
};

export type Paged<T> = {
  items: T[];
  pagination: {
    page: number;
    page_size: number;
    total: number;
  };
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

// ── Daily Pick ──

export type DailyPick = {
  id: number;
  pick_date: string;
  wallet_id: number;
  smart_score: number;
  realized_pnl: number;
  total_trades: number;
  win_rate: number;
  reason_json: Record<string, unknown>;
  reason_summary: string;
  reason_summary_zh: string;
  model_id: string;
  trades_followed: number;
  follow_pnl?: number;
  result_updated_at?: string;
  created_at: string;
};

// ── Nova Session (thinking timeline) ──

export type NovaSession = {
  id: number;
  session_date: string;
  round: number;
  phase: "analyzing" | "final" | "verified";
  candidates_json: Record<string, unknown>;
  observations_json: Record<string, unknown>;
  decision_json: Record<string, unknown>;
  nl_summary: string;
  nl_summary_zh: string;
  picked_wallet_id?: number;
  model_id: string;
  input_tokens: number;
  output_tokens: number;
  latency_ms: number;
  confidence_score?: number;
  focus_metrics?: Record<string, unknown>;
  hesitation_points?: Record<string, unknown>;
  created_at: string;
};


// ── Nova Insight (AI Brain Visualization) ──

export type NovaStatus = {
  is_active: boolean;
  current_round: number;
  total_rounds: number;
  phase: "analyzing" | "final" | "verified";
  confidence_score?: number;
  focus_metrics: string[];
  candidate_count: number;
  last_updated: string;
  next_round_at?: string;
  session_date: string;
};

export type ThinkingRound = {
  session: NovaSession;
  confidence_change?: number;
  is_breakthrough: boolean;
  is_hesitation: boolean;
};

export type CandidateScore = {
  wallet_id: number;
  address: string;
  pseudonym?: string;
  nova_score: number;
  win_rate: number;
  stability: number;
  activity: number;
  nova_comment: string;
  nova_comment_zh: string;
  rank: number;
};

export type DecisionExplanation = {
  pick_id: number;
  wallet_id: number;
  weight_distribution: Record<string, number>;
  metric_comparison: MetricComparison[];
  key_reasons: string[];
  key_reasons_zh: string[];
};

export type MetricComparison = {
  metric: string;
  metric_zh: string;
  pick_value: number;
  average_value: number;
  nova_standard: string;
  passed: boolean;
};


// ── Nova Memory (Learning & History) ──

export type LearningRecord = {
  validation_date: string;
  pick_wallet_id: number;
  wallet_address: string;
  follow_pnl?: number;
  is_success: boolean;
  lesson_learned: string;
  lesson_learned_zh: string;
  strategy_adjustment?: Record<string, unknown>;
};

export type StrategyChange = {
  date: string;
  description: string;
  description_zh: string;
};

export type MemorySummary = {
  total_validations: number;
  success_count: number;
  failure_count: number;
  success_rate: number;
  weekly_success_rate: number;
  recent_lessons: string[];
  recent_lessons_zh: string[];
  strategy_evolution: StrategyChange[];
};

export type NovaMemory = {
  summary: MemorySummary;
  history: LearningRecord[];
};
