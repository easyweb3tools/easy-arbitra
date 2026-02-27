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

type AnomalyAlert struct {
	ID           int64          `gorm:"column:id;primaryKey" json:"id"`
	WalletID     int64          `gorm:"column:wallet_id;index" json:"wallet_id"`
	MarketID     *int64         `gorm:"column:market_id;index" json:"market_id,omitempty"`
	AlertType    string         `gorm:"column:alert_type;size:30;index" json:"alert_type"`
	Severity     int16          `gorm:"column:severity;index" json:"severity"`
	Evidence     datatypes.JSON `gorm:"column:evidence;type:jsonb" json:"evidence"`
	Description  string         `gorm:"column:description" json:"description"`
	Acknowledged bool           `gorm:"column:acknowledged;default:false" json:"acknowledged"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (AnomalyAlert) TableName() string { return "anomaly_alert" }

type Watchlist struct {
	ID              int64     `gorm:"column:id;primaryKey" json:"id"`
	WalletID        int64     `gorm:"column:wallet_id;uniqueIndex:idx_watchlist_wallet_user" json:"wallet_id"`
	UserFingerprint string    `gorm:"column:user_fingerprint;size:120;uniqueIndex:idx_watchlist_wallet_user" json:"user_fingerprint"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Watchlist) TableName() string { return "watchlist" }

type WalletUpdateEvent struct {
	ID             int64          `gorm:"column:id;primaryKey" json:"id"`
	WalletID       int64          `gorm:"column:wallet_id;index" json:"wallet_id"`
	EventType      string         `gorm:"column:event_type;size:40;index" json:"event_type"`
	Payload        datatypes.JSON `gorm:"column:payload;type:jsonb" json:"payload"`
	ActionRequired bool           `gorm:"column:action_required;default:false" json:"action_required"`
	Suggestion     *string        `gorm:"column:suggestion" json:"suggestion,omitempty"`
	SuggestionZh   *string        `gorm:"column:suggestion_zh" json:"suggestion_zh,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (WalletUpdateEvent) TableName() string { return "wallet_update_event" }

type Portfolio struct {
	ID          int64          `gorm:"column:id;primaryKey" json:"id"`
	Name        string         `gorm:"column:name;size:100" json:"name"`
	NameZh      *string        `gorm:"column:name_zh;size:100" json:"name_zh,omitempty"`
	Description *string        `gorm:"column:description" json:"description,omitempty"`
	RiskLevel   string         `gorm:"column:risk_level;size:10" json:"risk_level"`
	WalletIDs   datatypes.JSON `gorm:"column:wallet_ids;type:jsonb;not null" json:"wallet_ids"`
	IsActive    bool           `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
}

func (Portfolio) TableName() string { return "portfolio" }

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

type CopyTradingConfig struct {
	ID              int64     `gorm:"column:id;primaryKey" json:"id"`
	UserFingerprint string    `gorm:"column:user_fingerprint;size:120;uniqueIndex:idx_ctc_user_wallet" json:"user_fingerprint"`
	WalletID        int64     `gorm:"column:wallet_id;uniqueIndex:idx_ctc_user_wallet" json:"wallet_id"`
	Enabled         bool      `gorm:"column:enabled;default:true" json:"enabled"`
	MaxPositionUSDC float64   `gorm:"column:max_position_usdc;type:numeric(12,2);default:1000" json:"max_position_usdc"`
	RiskPreference  string    `gorm:"column:risk_preference;size:20;default:moderate" json:"risk_preference"`
	LastCheckedAt   time.Time `gorm:"column:last_checked_at" json:"last_checked_at"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (CopyTradingConfig) TableName() string { return "copy_trading_config" }

type CopyTradeDecision struct {
	ID            int64          `gorm:"column:id;primaryKey" json:"id"`
	ConfigID      int64          `gorm:"column:config_id;index" json:"config_id"`
	LeaderTradeID *int64         `gorm:"column:leader_trade_id" json:"leader_trade_id,omitempty"`
	MarketID      *int64         `gorm:"column:market_id" json:"market_id,omitempty"`
	MarketTitle   string         `gorm:"column:market_title" json:"market_title"`
	Decision      string         `gorm:"column:decision;size:10" json:"decision"`
	Confidence    float64        `gorm:"column:confidence;type:numeric(4,3)" json:"confidence"`
	Outcome       string         `gorm:"column:outcome;size:5" json:"outcome"`
	Action        string         `gorm:"column:action;size:5" json:"action"`
	Price         float64        `gorm:"column:price;type:numeric(12,6)" json:"price"`
	SizeUSDC      float64        `gorm:"column:size_usdc;type:numeric(12,2)" json:"size_usdc"`
	StopLossPrice *float64       `gorm:"column:stop_loss_price;type:numeric(12,6)" json:"stop_loss_price,omitempty"`
	Reasoning     string         `gorm:"column:reasoning" json:"reasoning"`
	ReasoningEn   string         `gorm:"column:reasoning_en" json:"reasoning_en"`
	RiskNotes     datatypes.JSON `gorm:"column:risk_notes;type:jsonb;default:'[]'" json:"risk_notes"`
	ModelID       string         `gorm:"column:model_id;size:50" json:"model_id"`
	InputTokens   int            `gorm:"column:input_tokens" json:"input_tokens"`
	OutputTokens  int            `gorm:"column:output_tokens" json:"output_tokens"`
	LatencyMS     int            `gorm:"column:latency_ms" json:"latency_ms"`
	Status        string         `gorm:"column:status;size:20;default:pending;index" json:"status"`
	ExecutedAt    *time.Time     `gorm:"column:executed_at" json:"executed_at,omitempty"`
	ClosedAt      *time.Time     `gorm:"column:closed_at" json:"closed_at,omitempty"`
	ClosePrice    *float64       `gorm:"column:close_price;type:numeric(12,6)" json:"close_price,omitempty"`
	RealizedPnL   *float64       `gorm:"column:realized_pnl;type:numeric(12,4)" json:"realized_pnl,omitempty"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
}

func (CopyTradeDecision) TableName() string { return "copy_trade_decision" }

type CopyTradeDailyPerf struct {
	ConfigID      int64     `gorm:"column:config_id;primaryKey" json:"config_id"`
	PerfDate      time.Time `gorm:"column:perf_date;primaryKey;type:date" json:"perf_date"`
	TotalCopies   int       `gorm:"column:total_copies" json:"total_copies"`
	Profitable    int       `gorm:"column:profitable" json:"profitable"`
	TotalPnL      float64   `gorm:"column:total_pnl;type:numeric(12,4)" json:"total_pnl"`
	TotalExposure float64   `gorm:"column:total_exposure;type:numeric(12,2)" json:"total_exposure"`
	Skipped       int       `gorm:"column:skipped" json:"skipped"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

func (CopyTradeDailyPerf) TableName() string { return "copy_trade_daily_perf" }

type User struct {
	ID           int64     `gorm:"column:id;primaryKey" json:"id"`
	Email        string    `gorm:"column:email;size:255;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;size:255" json:"-"`
	Name         string    `gorm:"column:name;size:100" json:"name"`
	AvatarURL    string    `gorm:"column:avatar_url;size:500" json:"avatar_url,omitempty"`
	Provider     string    `gorm:"column:provider;size:20;default:email" json:"provider"`
	ProviderID   string    `gorm:"column:provider_id;size:255" json:"-"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string { return "user_account" }
