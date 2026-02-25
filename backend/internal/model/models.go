package model

import (
	"time"

	"gorm.io/datatypes"
)

type Wallet struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Address     []byte     `gorm:"type:bytea;uniqueIndex:idx_wallet_address_chain" json:"-"`
	ChainID     int        `gorm:"uniqueIndex:idx_wallet_address_chain;default:137" json:"chain_id"`
	Pseudonym   *string    `json:"pseudonym,omitempty"`
	IsTracked   bool       `gorm:"default:false" json:"is_tracked"`
	FirstSeenAt time.Time  `json:"first_seen_at"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Wallet) TableName() string { return "wallet" }

type Market struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	ConditionID     string     `gorm:"uniqueIndex;size:66" json:"condition_id"`
	Slug            string     `gorm:"size:255" json:"slug"`
	Title           string     `json:"title"`
	Category        string     `gorm:"size:50" json:"category"`
	Status          int16      `gorm:"default:0" json:"status"`
	HasFee          bool       `gorm:"default:false" json:"has_fee"`
	ResolutionTime  *time.Time `json:"resolution_time,omitempty"`
	ResolvedOutcome *int16     `json:"resolved_outcome,omitempty"`
	Volume          float64    `json:"volume"`
	Liquidity       float64    `json:"liquidity"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (Market) TableName() string { return "market" }

type Token struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	MarketID  int64     `gorm:"index" json:"market_id"`
	TokenID   string    `gorm:"uniqueIndex;size:80" json:"token_id"`
	Side      int16     `json:"side"`
	CreatedAt time.Time `json:"created_at"`
}

func (Token) TableName() string { return "token" }

type TradeFill struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TokenID       int64     `gorm:"index" json:"token_id"`
	MakerWalletID *int64    `gorm:"index" json:"maker_wallet_id,omitempty"`
	TakerWalletID *int64    `gorm:"index" json:"taker_wallet_id,omitempty"`
	Side          int16     `json:"side"`
	Price         float64   `json:"price"`
	Size          float64   `json:"size"`
	FeePaid       float64   `json:"fee_paid"`
	TxHash        []byte    `gorm:"type:bytea" json:"-"`
	BlockNumber   int64     `json:"block_number"`
	BlockTime     time.Time `gorm:"index" json:"block_time"`
	Source        int16     `json:"source"`
	UniqKey       string    `gorm:"uniqueIndex;size:130" json:"uniq_key"`
	CreatedAt     time.Time `json:"created_at"`
}

func (TradeFill) TableName() string { return "trade_fill" }

type OffchainEvent struct {
	ID            int64          `gorm:"primaryKey" json:"id"`
	MarketID      *int64         `gorm:"index" json:"market_id,omitempty"`
	SourceEventID string         `gorm:"column:source_event_id;size:120;default:'';uniqueIndex:idx_offchain_event_source_event" json:"source_event_id"`
	EventTime     time.Time      `gorm:"index" json:"event_time"`
	EventType     string         `gorm:"size:30" json:"event_type"`
	Source        string         `gorm:"column:source_name;size:100;uniqueIndex:idx_offchain_event_source_event" json:"source_name"`
	Title         string         `json:"title"`
	Payload       datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	CreatedAt     time.Time      `json:"created_at"`
}

func (OffchainEvent) TableName() string { return "offchain_event" }

type WalletFeaturesDaily struct {
	WalletID          int64     `gorm:"primaryKey" json:"wallet_id"`
	FeatureDate       time.Time `gorm:"primaryKey;type:date" json:"feature_date"`
	Pnl7d             float64   `json:"pnl_7d"`
	Pnl30d            float64   `json:"pnl_30d"`
	Pnl90d            float64   `json:"pnl_90d"`
	MakerRatio        float64   `json:"maker_ratio"`
	TradeCount30d     int       `json:"trade_count_30d"`
	ActiveDays30d     int       `json:"active_days_30d"`
	TxFrequencyPerDay float64   `json:"tx_frequency_per_day"`
	AvgEdge           float64   `json:"avg_edge"`
	CreatedAt         time.Time `json:"created_at"`
}

func (WalletFeaturesDaily) TableName() string { return "wallet_features_daily" }

type WalletScore struct {
	ID                 int64          `gorm:"primaryKey" json:"id"`
	WalletID           int64          `gorm:"index" json:"wallet_id"`
	ScoredAt           time.Time      `json:"scored_at"`
	StrategyType       string         `gorm:"size:30" json:"strategy_type"`
	StrategyConfidence float64        `json:"strategy_confidence"`
	InfoEdgeLevel      string         `gorm:"size:20" json:"info_edge_level"`
	InfoEdgeConfidence float64        `json:"info_edge_confidence"`
	SmartScore         int            `json:"smart_score"`
	PoolTier           string         `gorm:"size:20;default:observation" json:"pool_tier"`
	PoolTierUpdatedAt  *time.Time     `json:"pool_tier_updated_at,omitempty"`
	SuitableFor        string         `gorm:"size:50" json:"suitable_for"`
	RiskLevel          string         `gorm:"size:10" json:"risk_level"`
	SuggestedPosition  string         `gorm:"size:20" json:"suggested_position"`
	Momentum           string         `gorm:"size:20" json:"momentum"`
	ScoringDetail      datatypes.JSON `gorm:"type:jsonb" json:"scoring_detail"`
	CreatedAt          time.Time      `json:"created_at"`
}

func (WalletScore) TableName() string { return "wallet_score" }

type AIAnalysisReport struct {
	ID           int64          `gorm:"primaryKey" json:"id"`
	WalletID     int64          `gorm:"index" json:"wallet_id"`
	ModelID      string         `gorm:"size:50" json:"model_id"`
	Report       datatypes.JSON `gorm:"type:jsonb" json:"report"`
	NLSummary    string         `json:"nl_summary"`
	RiskWarnings datatypes.JSON `gorm:"type:jsonb" json:"risk_warnings"`
	InputTokens  int            `json:"input_tokens"`
	OutputTokens int            `json:"output_tokens"`
	LatencyMS    int            `json:"latency_ms"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (AIAnalysisReport) TableName() string { return "ai_analysis_report" }

type AnomalyAlert struct {
	ID           int64          `gorm:"primaryKey" json:"id"`
	WalletID     int64          `gorm:"index" json:"wallet_id"`
	MarketID     *int64         `gorm:"index" json:"market_id,omitempty"`
	AlertType    string         `gorm:"size:30;index" json:"alert_type"`
	Severity     int16          `gorm:"index" json:"severity"`
	Evidence     datatypes.JSON `gorm:"type:jsonb" json:"evidence"`
	Description  string         `json:"description"`
	Acknowledged bool           `gorm:"default:false" json:"acknowledged"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (AnomalyAlert) TableName() string { return "anomaly_alert" }

type Watchlist struct {
	ID              int64     `gorm:"primaryKey" json:"id"`
	WalletID        int64     `gorm:"uniqueIndex:idx_watchlist_wallet_user" json:"wallet_id"`
	UserFingerprint string    `gorm:"size:120;uniqueIndex:idx_watchlist_wallet_user" json:"user_fingerprint"`
	CreatedAt       time.Time `json:"created_at"`
}

func (Watchlist) TableName() string { return "watchlist" }

type WalletUpdateEvent struct {
	ID             int64          `gorm:"primaryKey" json:"id"`
	WalletID       int64          `gorm:"index" json:"wallet_id"`
	EventType      string         `gorm:"size:40;index" json:"event_type"`
	Payload        datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	ActionRequired bool           `gorm:"default:false" json:"action_required"`
	Suggestion     *string        `json:"suggestion,omitempty"`
	SuggestionZh   *string        `json:"suggestion_zh,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
}

func (WalletUpdateEvent) TableName() string { return "wallet_update_event" }

type Portfolio struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100" json:"name"`
	NameZh      *string        `gorm:"size:100" json:"name_zh,omitempty"`
	Description *string        `json:"description,omitempty"`
	RiskLevel   string         `gorm:"size:10" json:"risk_level"`
	WalletIDs   datatypes.JSON `gorm:"type:jsonb;not null" json:"wallet_ids"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (Portfolio) TableName() string { return "portfolio" }

type IngestCursor struct {
	Source      string    `gorm:"primaryKey;size:100" json:"source"`
	Stream      string    `gorm:"primaryKey;size:120" json:"stream"`
	CursorValue string    `gorm:"type:text;not null" json:"cursor_value"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (IngestCursor) TableName() string { return "ingest_cursor" }

type IngestRun struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	JobName   string         `gorm:"type:text;not null" json:"job_name"`
	StartedAt time.Time      `json:"started_at"`
	EndedAt   *time.Time     `json:"ended_at,omitempty"`
	Status    string         `gorm:"type:text;not null" json:"status"`
	Stats     datatypes.JSON `gorm:"type:jsonb;not null" json:"stats"`
	ErrorText *string        `json:"error_text,omitempty"`
}

func (IngestRun) TableName() string { return "ingest_run" }
