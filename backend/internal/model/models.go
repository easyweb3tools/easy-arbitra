package model

import (
	"time"

	"gorm.io/datatypes"
)

type Wallet struct {
	ID          int64      `gorm:"column:id;primaryKey" json:"id"`
	Address     []byte     `gorm:"column:address;type:bytea;uniqueIndex:idx_wallet_address_chain" json:"-"`
	ChainID     int        `gorm:"column:chain_id;uniqueIndex:idx_wallet_address_chain;default:137" json:"chain_id"`
	Pseudonym   *string    `gorm:"column:pseudonym" json:"pseudonym,omitempty"`
	IsTracked   bool       `gorm:"column:is_tracked;default:false" json:"is_tracked"`
	FirstSeenAt time.Time  `gorm:"column:first_seen_at" json:"first_seen_at"`
	LastSeenAt  *time.Time `gorm:"column:last_seen_at" json:"last_seen_at,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Wallet) TableName() string { return "wallet" }

type Market struct {
	ID              int64      `gorm:"column:id;primaryKey" json:"id"`
	ConditionID     string     `gorm:"column:condition_id;uniqueIndex;size:66" json:"condition_id"`
	Slug            string     `gorm:"column:slug;size:255" json:"slug"`
	Title           string     `gorm:"column:title" json:"title"`
	Category        string     `gorm:"column:category;size:50" json:"category"`
	Status          int16      `gorm:"column:status;default:0" json:"status"`
	HasFee          bool       `gorm:"column:has_fee;default:false" json:"has_fee"`
	ResolutionTime  *time.Time `gorm:"column:resolution_time" json:"resolution_time,omitempty"`
	ResolvedOutcome *int16     `gorm:"column:resolved_outcome" json:"resolved_outcome,omitempty"`
	Volume          float64    `gorm:"column:volume" json:"volume"`
	Liquidity       float64    `gorm:"column:liquidity" json:"liquidity"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (Market) TableName() string { return "market" }

type Token struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	MarketID  int64     `gorm:"column:market_id;index" json:"market_id"`
	TokenID   string    `gorm:"column:token_id;uniqueIndex;size:80" json:"token_id"`
	Side      int16     `gorm:"column:side" json:"side"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Token) TableName() string { return "token" }

type TradeFill struct {
	ID            int64     `gorm:"column:id;primaryKey" json:"id"`
	TokenID       int64     `gorm:"column:token_id;index" json:"token_id"`
	MakerWalletID *int64    `gorm:"column:maker_wallet_id;index" json:"maker_wallet_id,omitempty"`
	TakerWalletID *int64    `gorm:"column:taker_wallet_id;index" json:"taker_wallet_id,omitempty"`
	Side          int16     `gorm:"column:side" json:"side"`
	Price         float64   `gorm:"column:price" json:"price"`
	Size          float64   `gorm:"column:size" json:"size"`
	FeePaid       float64   `gorm:"column:fee_paid" json:"fee_paid"`
	TxHash        []byte    `gorm:"column:tx_hash;type:bytea" json:"-"`
	BlockNumber   int64     `gorm:"column:block_number" json:"block_number"`
	BlockTime     time.Time `gorm:"column:block_time;index" json:"block_time"`
	Source        int16     `gorm:"column:source" json:"source"`
	UniqKey       string    `gorm:"column:uniq_key;uniqueIndex;size:130" json:"uniq_key"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

func (TradeFill) TableName() string { return "trade_fill" }

type OffchainEvent struct {
	ID            int64          `gorm:"column:id;primaryKey" json:"id"`
	MarketID      *int64         `gorm:"column:market_id;index" json:"market_id,omitempty"`
	SourceEventID string         `gorm:"column:source_event_id;size:120;default:'';uniqueIndex:idx_offchain_event_source_event" json:"source_event_id"`
	EventTime     time.Time      `gorm:"column:event_time;index" json:"event_time"`
	EventType     string         `gorm:"column:event_type;size:30" json:"event_type"`
	Source        string         `gorm:"column:source_name;size:100;uniqueIndex:idx_offchain_event_source_event" json:"source_name"`
	Title         string         `gorm:"column:title" json:"title"`
	Payload       datatypes.JSON `gorm:"column:payload;type:jsonb" json:"payload"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (OffchainEvent) TableName() string { return "offchain_event" }

type WalletFeaturesDaily struct {
	WalletID          int64     `gorm:"column:wallet_id;primaryKey" json:"wallet_id"`
	FeatureDate       time.Time `gorm:"column:feature_date;primaryKey;type:date" json:"feature_date"`
	Pnl7d             float64   `gorm:"column:pnl_7d" json:"pnl_7d"`
	Pnl30d            float64   `gorm:"column:pnl_30d" json:"pnl_30d"`
	Pnl90d            float64   `gorm:"column:pnl_90d" json:"pnl_90d"`
	MakerRatio        float64   `gorm:"column:maker_ratio" json:"maker_ratio"`
	TradeCount30d     int       `gorm:"column:trade_count_30d" json:"trade_count_30d"`
	ActiveDays30d     int       `gorm:"column:active_days_30d" json:"active_days_30d"`
	TxFrequencyPerDay float64   `gorm:"column:tx_frequency_per_day" json:"tx_frequency_per_day"`
	AvgEdge           float64   `gorm:"column:avg_edge" json:"avg_edge"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
}

func (WalletFeaturesDaily) TableName() string { return "wallet_features_daily" }

type WalletScore struct {
	ID                 int64          `gorm:"column:id;primaryKey" json:"id"`
	WalletID           int64          `gorm:"column:wallet_id;index" json:"wallet_id"`
	ScoredAt           time.Time      `gorm:"column:scored_at" json:"scored_at"`
	StrategyType       string         `gorm:"column:strategy_type;size:30" json:"strategy_type"`
	StrategyConfidence float64        `gorm:"column:strategy_confidence" json:"strategy_confidence"`
	InfoEdgeLevel      string         `gorm:"column:info_edge_level;size:20" json:"info_edge_level"`
	InfoEdgeConfidence float64        `gorm:"column:info_edge_confidence" json:"info_edge_confidence"`
	SmartScore         int            `gorm:"column:smart_score" json:"smart_score"`
	PoolTier           string         `gorm:"column:pool_tier;size:20;default:observation" json:"pool_tier"`
	PoolTierUpdatedAt  *time.Time     `gorm:"column:pool_tier_updated_at" json:"pool_tier_updated_at,omitempty"`
	SuitableFor        string         `gorm:"column:suitable_for;size:50" json:"suitable_for"`
	RiskLevel          string         `gorm:"column:risk_level;size:10" json:"risk_level"`
	SuggestedPosition  string         `gorm:"column:suggested_position;size:20" json:"suggested_position"`
	Momentum           string         `gorm:"column:momentum;size:20" json:"momentum"`
	ScoringDetail      datatypes.JSON `gorm:"column:scoring_detail;type:jsonb" json:"scoring_detail"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (WalletScore) TableName() string { return "wallet_score" }

type AIAnalysisReport struct {
	ID           int64          `gorm:"column:id;primaryKey" json:"id"`
	WalletID     int64          `gorm:"column:wallet_id;index" json:"wallet_id"`
	ModelID      string         `gorm:"column:model_id;size:50" json:"model_id"`
	Report       datatypes.JSON `gorm:"column:report;type:jsonb" json:"report"`
	NLSummary    string         `gorm:"column:nl_summary" json:"nl_summary"`
	RiskWarnings datatypes.JSON `gorm:"column:risk_warnings;type:jsonb" json:"risk_warnings"`
	InputTokens  int            `gorm:"column:input_tokens" json:"input_tokens"`
	OutputTokens int            `gorm:"column:output_tokens" json:"output_tokens"`
	LatencyMS    int            `gorm:"column:latency_ms" json:"latency_ms"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (AIAnalysisReport) TableName() string { return "ai_analysis_report" }

type IngestCursor struct {
	Source      string    `gorm:"column:source;primaryKey;size:100" json:"source"`
	Stream      string    `gorm:"column:stream;primaryKey;size:120" json:"stream"`
	CursorValue string    `gorm:"column:cursor_value;type:text;not null" json:"cursor_value"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (IngestCursor) TableName() string { return "ingest_cursor" }

type IngestRun struct {
	ID        int64          `gorm:"column:id;primaryKey" json:"id"`
	JobName   string         `gorm:"column:job_name;type:text;not null" json:"job_name"`
	StartedAt time.Time      `gorm:"column:started_at" json:"started_at"`
	EndedAt   *time.Time     `gorm:"column:ended_at" json:"ended_at,omitempty"`
	Status    string         `gorm:"column:status;type:text;not null" json:"status"`
	Stats     datatypes.JSON `gorm:"column:stats;type:jsonb;not null" json:"stats"`
	ErrorText *string        `gorm:"column:error_text" json:"error_text,omitempty"`
}

func (IngestRun) TableName() string { return "ingest_run" }

// DailyPick represents the daily recommended best trader.
type DailyPick struct {
	ID              int64          `gorm:"column:id;primaryKey" json:"id"`
	PickDate        time.Time      `gorm:"column:pick_date;type:date;uniqueIndex" json:"pick_date"`
	WalletID        int64          `gorm:"column:wallet_id;index" json:"wallet_id"`
	SmartScore      int            `gorm:"column:smart_score" json:"smart_score"`
	RealizedPnL     float64        `gorm:"column:realized_pnl;type:numeric(16,6)" json:"realized_pnl"`
	TotalTrades     int64          `gorm:"column:total_trades" json:"total_trades"`
	WinRate         float64        `gorm:"column:win_rate;type:numeric(5,4)" json:"win_rate"`
	ReasonJSON      datatypes.JSON `gorm:"column:reason_json;type:jsonb" json:"reason_json"`
	ReasonSummary   string         `gorm:"column:reason_summary" json:"reason_summary"`
	ReasonSummaryZh string         `gorm:"column:reason_summary_zh" json:"reason_summary_zh"`
	ModelID         string         `gorm:"column:model_id;size:50" json:"model_id"`
	TradesFollowed  int            `gorm:"column:trades_followed" json:"trades_followed"`
	FollowPnL       *float64       `gorm:"column:follow_pnl;type:numeric(16,6)" json:"follow_pnl,omitempty"`
	ResultUpdatedAt *time.Time     `gorm:"column:result_updated_at" json:"result_updated_at,omitempty"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (DailyPick) TableName() string { return "daily_pick" }

// NovaSession represents a single round of Nova's analysis.
// Multiple sessions per day form Nova's "working memory".
type NovaSession struct {
	ID               int64          `gorm:"column:id;primaryKey" json:"id"`
	SessionDate      time.Time      `gorm:"column:session_date;type:date;index" json:"session_date"`
	Round            int            `gorm:"column:round" json:"round"`
	Phase            string         `gorm:"column:phase;size:20" json:"phase"` // "analyzing" | "final" | "verified"
	CandidatesJSON   datatypes.JSON `gorm:"column:candidates_json;type:jsonb" json:"candidates_json"`
	ObservationsJSON datatypes.JSON `gorm:"column:observations_json;type:jsonb" json:"observations_json"`
	DecisionJSON     datatypes.JSON `gorm:"column:decision_json;type:jsonb" json:"decision_json"`
	NLSummary        string         `gorm:"column:nl_summary" json:"nl_summary"`
	NLSummaryZh      string         `gorm:"column:nl_summary_zh" json:"nl_summary_zh"`
	PickedWalletID   *int64         `gorm:"column:picked_wallet_id" json:"picked_wallet_id,omitempty"`
	ModelID          string         `gorm:"column:model_id;size:50" json:"model_id"`
	InputTokens      int            `gorm:"column:input_tokens" json:"input_tokens"`
	OutputTokens     int            `gorm:"column:output_tokens" json:"output_tokens"`
	LatencyMS        int            `gorm:"column:latency_ms" json:"latency_ms"`
	ConfidenceScore  *float64       `gorm:"column:confidence_score;type:decimal(5,2)" json:"confidence_score,omitempty"`
	FocusMetrics     datatypes.JSON `gorm:"column:focus_metrics;type:jsonb" json:"focus_metrics,omitempty"`
	HesitationPoints datatypes.JSON `gorm:"column:hesitation_points;type:jsonb" json:"hesitation_points,omitempty"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (NovaSession) TableName() string { return "nova_session" }

// NovaLearningLog represents Nova's learning from validation results.
type NovaLearningLog struct {
	ID                 int64          `gorm:"column:id;primaryKey" json:"id"`
	ValidationDate     time.Time      `gorm:"column:validation_date;type:date;not null;index" json:"validation_date"`
	PickWalletID       int64          `gorm:"column:pick_wallet_id;not null" json:"pick_wallet_id"`
	FollowPnL          *float64       `gorm:"column:follow_pnl;type:decimal(16,6)" json:"follow_pnl,omitempty"`
	IsSuccess          bool           `gorm:"column:is_success" json:"is_success"`
	LessonLearned      string         `gorm:"column:lesson_learned" json:"lesson_learned"`
	LessonLearnedZh    string         `gorm:"column:lesson_learned_zh" json:"lesson_learned_zh"`
	StrategyAdjustment datatypes.JSON `gorm:"column:strategy_adjustment;type:jsonb" json:"strategy_adjustment,omitempty"`
	CreatedAt          time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (NovaLearningLog) TableName() string { return "nova_learning_log" }

// TraderStats represents pre-aggregated trading statistics for each wallet.
// This table is refreshed periodically by workers to avoid real-time aggregation of 12.4M trade_fill rows.
type TraderStats struct {
	WalletID     int64     `gorm:"column:wallet_id;primaryKey" json:"wallet_id"`
	TradeCount   int64     `gorm:"column:trade_count;not null;default:0;index:idx_trader_stats_trade_count,sort:desc" json:"trade_count"`
	TradingPnL   float64   `gorm:"column:trading_pnl;type:numeric(20,6);not null;default:0" json:"trading_pnl"`
	MakerRebates float64   `gorm:"column:maker_rebates;type:numeric(20,6);not null;default:0" json:"maker_rebates"`
	RealizedPnL  float64   `gorm:"column:realized_pnl;type:numeric(20,6);generated always as (trading_pnl + maker_rebates) stored;index:idx_trader_stats_realized_pnl,sort:desc" json:"realized_pnl"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (TraderStats) TableName() string { return "trader_stats" }
