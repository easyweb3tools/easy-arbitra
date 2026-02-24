package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/pkg/polyaddr"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct{ db *gorm.DB }

type MarketRepository struct{ db *gorm.DB }

type TokenRepository struct{ db *gorm.DB }

type TradeRepository struct{ db *gorm.DB }

type ScoreRepository struct{ db *gorm.DB }

type AIReportRepository struct{ db *gorm.DB }

type WatchlistRepository struct{ db *gorm.DB }

type PortfolioRepository struct{ db *gorm.DB }

type WalletListFilter struct {
	Tracked *bool
	Search  string
	SortBy  string
	Order   string
	Limit   int
	Offset  int
}

type PotentialWalletFilter struct {
	MinTrades      int64
	MinRealizedPnL float64
	StrategyType   string
	PoolTier       string
	HasAIReport    *bool
	SortBy         string
	Order          string
	Limit          int
	Offset         int
}

type MarketListFilter struct {
	Category string
	Status   *int16
	SortBy   string
	Order    string
	Limit    int
	Offset   int
}

type WalletPnLSummary struct {
	TradingPnL   float64 `gorm:"column:trading_pnl" json:"trading_pnl"`
	MakerRebates float64 `gorm:"column:maker_rebates" json:"maker_rebates"`
	FeesPaid     float64 `gorm:"column:fees_paid" json:"fees_paid"`
	TotalTrades  int64   `gorm:"column:total_trades" json:"total_trades"`
	Volume30D    float64 `gorm:"column:volume_30d" json:"volume_30d"`
}

type WalletTimingSummary struct {
	MeanDeltaMinutes float64 `json:"mean_delta_minutes"`
	StdDevMinutes    float64 `json:"stddev_minutes"`
	Samples          int64   `json:"samples"`
}

type LeaderboardRow struct {
	WalletID      int64   `json:"wallet_id"`
	Address       []byte  `json:"-"`
	Pseudonym     *string `json:"pseudonym,omitempty"`
	StrategyType  string  `json:"strategy_type"`
	SmartScore    int     `json:"smart_score"`
	InfoEdgeLevel string  `json:"info_edge_level"`
	ScoredAt      string  `json:"scored_at"`
}

type WalletTradeCountRow struct {
	WalletID   int64  `gorm:"column:wallet_id"`
	Address    []byte `gorm:"column:address"`
	TradeCount int64  `gorm:"column:trade_count"`
}

type PotentialWalletRow struct {
	WalletID       int64      `gorm:"column:wallet_id"`
	Address        []byte     `gorm:"column:address"`
	Pseudonym      *string    `gorm:"column:pseudonym"`
	IsTracked      bool       `gorm:"column:is_tracked"`
	TradeCount     int64      `gorm:"column:trade_count"`
	TradingPnL     float64    `gorm:"column:trading_pnl"`
	MakerRebates   float64    `gorm:"column:maker_rebates"`
	RealizedPnL    float64    `gorm:"column:realized_pnl"`
	SmartScore     int        `gorm:"column:smart_score"`
	InfoEdgeLevel  string     `gorm:"column:info_edge_level"`
	StrategyType   string     `gorm:"column:strategy_type"`
	PoolTier       string     `gorm:"column:pool_tier"`
	NLSummary      string     `gorm:"column:nl_summary"`
	LastAnalyzedAt *time.Time `gorm:"column:last_analyzed_at"`
}

type WalletAICandidateRow struct {
	WalletID       int64      `gorm:"column:wallet_id"`
	Address        []byte     `gorm:"column:address"`
	TradeCount     int64      `gorm:"column:trade_count"`
	TradingPnL     float64    `gorm:"column:trading_pnl"`
	MakerRebates   float64    `gorm:"column:maker_rebates"`
	RealizedPnL    float64    `gorm:"column:realized_pnl"`
	LastAnalyzedAt *time.Time `gorm:"column:last_analyzed_at"`
}

type OpsTopRealizedRow struct {
	WalletID        int64      `gorm:"column:wallet_id"`
	Address         []byte     `gorm:"column:address"`
	TradeCount      int64      `gorm:"column:trade_count"`
	RealizedPnL     float64    `gorm:"column:realized_pnl"`
	RealizedPnL24h  float64    `gorm:"column:realized_pnl_24h"`
	HasAIReport     bool       `gorm:"column:has_ai_report"`
	LatestSummary   string     `gorm:"column:latest_summary"`
	LatestModelID   string     `gorm:"column:latest_model_id"`
	LastAnalyzedRaw *time.Time `gorm:"column:last_analyzed_at"`
}

type OpsTopAIConfidenceRow struct {
	WalletID       int64      `gorm:"column:wallet_id"`
	Address        []byte     `gorm:"column:address"`
	TradeCount     int64      `gorm:"column:trade_count"`
	RealizedPnL    float64    `gorm:"column:realized_pnl"`
	SmartScore     int        `gorm:"column:smart_score"`
	InfoEdgeLevel  string     `gorm:"column:info_edge_level"`
	StrategyType   string     `gorm:"column:strategy_type"`
	NLSummary      string     `gorm:"column:nl_summary"`
	LastAnalyzedAt *time.Time `gorm:"column:last_analyzed_at"`
}

type WatchlistWalletRow struct {
	WatchlistID    int64      `gorm:"column:watchlist_id"`
	WatchlistedAt  time.Time  `gorm:"column:watchlisted_at"`
	WalletID       int64      `gorm:"column:wallet_id"`
	Address        []byte     `gorm:"column:address"`
	Pseudonym      *string    `gorm:"column:pseudonym"`
	IsTracked      bool       `gorm:"column:is_tracked"`
	TradeCount     int64      `gorm:"column:trade_count"`
	TradingPnL     float64    `gorm:"column:trading_pnl"`
	MakerRebates   float64    `gorm:"column:maker_rebates"`
	RealizedPnL    float64    `gorm:"column:realized_pnl"`
	SmartScore     int        `gorm:"column:smart_score"`
	InfoEdgeLevel  string     `gorm:"column:info_edge_level"`
	StrategyType   string     `gorm:"column:strategy_type"`
	PoolTier       string     `gorm:"column:pool_tier"`
	NLSummary      string     `gorm:"column:nl_summary"`
	LastAnalyzedAt *time.Time `gorm:"column:last_analyzed_at"`
}

type WatchlistFeedRow struct {
	EventID        int64           `gorm:"column:event_id"`
	WalletID       int64           `gorm:"column:wallet_id"`
	Address        []byte          `gorm:"column:address"`
	Pseudonym      *string         `gorm:"column:pseudonym"`
	EventType      string          `gorm:"column:event_type"`
	EventPayload   json.RawMessage `gorm:"column:event_payload"`
	ActionRequired bool            `gorm:"column:action_required"`
	Suggestion     *string         `gorm:"column:suggestion"`
	SuggestionZh   *string         `gorm:"column:suggestion_zh"`
	EventTime      time.Time       `gorm:"column:event_time"`
}

type WalletEventRow struct {
	EventID        int64           `gorm:"column:event_id"`
	EventType      string          `gorm:"column:event_type"`
	EventPayload   json.RawMessage `gorm:"column:event_payload"`
	ActionRequired bool            `gorm:"column:action_required"`
	Suggestion     *string         `gorm:"column:suggestion"`
	SuggestionZh   *string         `gorm:"column:suggestion_zh"`
	EventTime      time.Time       `gorm:"column:event_time"`
}

func NewWalletRepository(db *gorm.DB) *WalletRepository     { return &WalletRepository{db: db} }
func NewMarketRepository(db *gorm.DB) *MarketRepository     { return &MarketRepository{db: db} }
func NewTokenRepository(db *gorm.DB) *TokenRepository       { return &TokenRepository{db: db} }
func NewTradeRepository(db *gorm.DB) *TradeRepository       { return &TradeRepository{db: db} }
func NewScoreRepository(db *gorm.DB) *ScoreRepository       { return &ScoreRepository{db: db} }
func NewAIReportRepository(db *gorm.DB) *AIReportRepository { return &AIReportRepository{db: db} }
func NewWatchlistRepository(db *gorm.DB) *WatchlistRepository {
	return &WatchlistRepository{db: db}
}
func NewPortfolioRepository(db *gorm.DB) *PortfolioRepository { return &PortfolioRepository{db: db} }

func (r *WalletRepository) List(ctx context.Context, f WalletListFilter) ([]model.Wallet, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Wallet{})
	if f.Tracked != nil {
		q = q.Where("is_tracked = ?", *f.Tracked)
	}
	if f.Search != "" {
		q = q.Where("pseudonym ILIKE ?", "%"+f.Search+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := buildOrderClause(f.SortBy, f.Order, map[string]struct{}{
		"id": {}, "created_at": {}, "updated_at": {}, "first_seen_at": {},
	})

	var rows []model.Wallet
	err := q.Order(orderClause).Limit(f.Limit).Offset(f.Offset).Find(&rows).Error
	return rows, total, err
}

func (r *WalletRepository) GetByID(ctx context.Context, id int64) (*model.Wallet, error) {
	var row model.Wallet
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *WalletRepository) CountTracked(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Wallet{}).Where("is_tracked = ?", true).Count(&total).Error
	return total, err
}

func (r *WalletRepository) ListIDs(ctx context.Context) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).Model(&model.Wallet{}).Pluck("id", &ids).Error
	return ids, err
}

func (r *WalletRepository) ListBackfillCandidates(ctx context.Context, minTrades int64, maxTrades int64, limit int) ([]WalletTradeCountRow, error) {
	if limit <= 0 {
		limit = 10
	}
	rows := make([]WalletTradeCountRow, 0, limit)
	err := r.db.WithContext(ctx).Raw(`
WITH counts AS (
    SELECT taker_wallet_id AS wallet_id, COUNT(*) AS c
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
    GROUP BY taker_wallet_id
    UNION ALL
    SELECT maker_wallet_id AS wallet_id, COUNT(*) AS c
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
    GROUP BY maker_wallet_id
),
merged AS (
    SELECT wallet_id, SUM(c) AS trade_count
    FROM counts
    GROUP BY wallet_id
)
SELECT w.id AS wallet_id, w.address AS address, m.trade_count
FROM merged m
JOIN wallet w ON w.id = m.wallet_id
WHERE m.trade_count BETWEEN ? AND ?
ORDER BY m.trade_count DESC, w.id ASC
LIMIT ?`, minTrades, maxTrades, limit).Scan(&rows).Error
	return rows, err
}

func (r *WalletRepository) ListAIAnalyzeCandidates(
	ctx context.Context,
	minTrades int64,
	minRealizedPnL float64,
	cooldown time.Duration,
	limit int,
) ([]WalletAICandidateRow, error) {
	if limit <= 0 {
		limit = 10
	}
	if minTrades <= 0 {
		minTrades = 100
	}
	cooldownSeconds := int64(cooldown / time.Second)
	rows := make([]WalletAICandidateRow, 0, limit)
	err := r.db.WithContext(ctx).Raw(`
WITH taker AS (
    SELECT
        taker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE
            WHEN side = 0 THEN (price * size) - fee_paid
            WHEN side = 1 THEN -((price * size) + fee_paid)
            ELSE 0
        END), 0) AS trading_pnl
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
    GROUP BY taker_wallet_id
),
maker AS (
    SELECT
        maker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
    GROUP BY maker_wallet_id
),
merged AS (
    SELECT
        x.wallet_id,
        SUM(x.trade_count) AS trade_count,
        SUM(x.trading_pnl) AS trading_pnl,
        SUM(x.maker_rebates) AS maker_rebates
    FROM (
        SELECT wallet_id, trade_count, trading_pnl, 0::numeric AS maker_rebates FROM taker
        UNION ALL
        SELECT wallet_id, trade_count, 0::numeric AS trading_pnl, maker_rebates FROM maker
    ) x
    GROUP BY x.wallet_id
),
latest AS (
    SELECT wallet_id, MAX(created_at) AS last_analyzed_at
    FROM ai_analysis_report
    GROUP BY wallet_id
)
SELECT
    w.id AS wallet_id,
    w.address AS address,
    m.trade_count,
    m.trading_pnl,
    m.maker_rebates,
    (m.trading_pnl + m.maker_rebates) AS realized_pnl,
    l.last_analyzed_at
FROM merged m
JOIN wallet w ON w.id = m.wallet_id
LEFT JOIN latest l ON l.wallet_id = w.id
WHERE
    m.trade_count >= ?
    AND (m.trading_pnl + m.maker_rebates) > ?
    AND (? <= 0 OR l.last_analyzed_at IS NULL OR l.last_analyzed_at < NOW() - make_interval(secs => ?))
ORDER BY realized_pnl DESC, m.trade_count DESC, w.id ASC
LIMIT ?`, minTrades, minRealizedPnL, cooldownSeconds, cooldownSeconds, limit).Scan(&rows).Error
	return rows, err
}

func (r *WalletRepository) CountPotentialWallets(ctx context.Context, f PotentialWalletFilter) (int64, error) {
	if f.MinTrades <= 0 {
		f.MinTrades = 100
	}
	hasAI := ""
	if f.HasAIReport != nil {
		hasAI = fmt.Sprintf("%t", *f.HasAIReport)
	}
	var total int64
	err := r.db.WithContext(ctx).Raw(`
WITH taker AS (
    SELECT
        taker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE
            WHEN side = 0 THEN (price * size) - fee_paid
            WHEN side = 1 THEN -((price * size) + fee_paid)
            ELSE 0
        END), 0) AS trading_pnl
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
    GROUP BY taker_wallet_id
),
maker AS (
    SELECT
        maker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
    GROUP BY maker_wallet_id
),
merged AS (
    SELECT
        x.wallet_id,
        SUM(x.trade_count) AS trade_count,
        SUM(x.trading_pnl) AS trading_pnl,
        SUM(x.maker_rebates) AS maker_rebates
    FROM (
        SELECT wallet_id, trade_count, trading_pnl, 0::numeric AS maker_rebates FROM taker
        UNION ALL
        SELECT wallet_id, trade_count, 0::numeric AS trading_pnl, maker_rebates FROM maker
    ) x
    GROUP BY x.wallet_id
),
latest_score AS (
    SELECT DISTINCT ON (wallet_id)
        wallet_id, strategy_type, pool_tier
    FROM wallet_score
    ORDER BY wallet_id, scored_at DESC
),
latest_ai AS (
    SELECT wallet_id, MAX(created_at) AS last_analyzed_at
    FROM ai_analysis_report
    GROUP BY wallet_id
)
SELECT COUNT(*)
FROM merged m
JOIN wallet w ON w.id = m.wallet_id
LEFT JOIN latest_score s ON s.wallet_id = w.id
LEFT JOIN latest_ai a ON a.wallet_id = w.id
WHERE m.trade_count >= ?
  AND (m.trading_pnl + m.maker_rebates) > ?
  AND (? = '' OR COALESCE(s.strategy_type, 'unknown') = ?)
  AND (? = '' OR COALESCE(s.pool_tier, 'observation') = ?)
  AND (
    ? = '' OR
    (? = 'true' AND a.last_analyzed_at IS NOT NULL) OR
    (? = 'false' AND a.last_analyzed_at IS NULL)
  )`,
		f.MinTrades,
		f.MinRealizedPnL,
		strings.TrimSpace(f.StrategyType),
		strings.TrimSpace(f.StrategyType),
		strings.TrimSpace(f.PoolTier),
		strings.TrimSpace(f.PoolTier),
		hasAI,
		hasAI,
		hasAI,
	).Scan(&total).Error
	return total, err
}

func (r *WalletRepository) ListPotentialWallets(ctx context.Context, f PotentialWalletFilter) ([]PotentialWalletRow, error) {
	if f.MinTrades <= 0 {
		f.MinTrades = 100
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Limit > 200 {
		f.Limit = 200
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	hasAI := ""
	if f.HasAIReport != nil {
		hasAI = fmt.Sprintf("%t", *f.HasAIReport)
	}
	orderClause := buildOrderClause(f.SortBy, f.Order, map[string]struct{}{
		"trade_count":      {},
		"realized_pnl":     {},
		"smart_score":      {},
		"last_analyzed_at": {},
		"wallet_id":        {},
	})
	query := fmt.Sprintf(`
WITH taker AS (
    SELECT
        taker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE
            WHEN side = 0 THEN (price * size) - fee_paid
            WHEN side = 1 THEN -((price * size) + fee_paid)
            ELSE 0
        END), 0) AS trading_pnl
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
    GROUP BY taker_wallet_id
),
maker AS (
    SELECT
        maker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
    GROUP BY maker_wallet_id
),
merged AS (
    SELECT
        x.wallet_id,
        SUM(x.trade_count) AS trade_count,
        SUM(x.trading_pnl) AS trading_pnl,
        SUM(x.maker_rebates) AS maker_rebates
    FROM (
        SELECT wallet_id, trade_count, trading_pnl, 0::numeric AS maker_rebates FROM taker
        UNION ALL
        SELECT wallet_id, trade_count, 0::numeric AS trading_pnl, maker_rebates FROM maker
    ) x
    GROUP BY x.wallet_id
),
latest_score AS (
    SELECT DISTINCT ON (wallet_id)
        wallet_id, smart_score, info_edge_level, strategy_type, pool_tier
    FROM wallet_score
    ORDER BY wallet_id, scored_at DESC
),
latest_ai AS (
    SELECT DISTINCT ON (wallet_id) wallet_id, created_at AS last_analyzed_at, nl_summary
    FROM ai_analysis_report
    ORDER BY wallet_id, created_at DESC
)
SELECT
    w.id AS wallet_id,
    w.address AS address,
    w.pseudonym AS pseudonym,
    w.is_tracked AS is_tracked,
    m.trade_count AS trade_count,
    m.trading_pnl AS trading_pnl,
    m.maker_rebates AS maker_rebates,
    (m.trading_pnl + m.maker_rebates) AS realized_pnl,
    COALESCE(s.smart_score, 0) AS smart_score,
    COALESCE(s.info_edge_level, 'unknown') AS info_edge_level,
    COALESCE(s.strategy_type, 'unknown') AS strategy_type,
    COALESCE(s.pool_tier, 'observation') AS pool_tier,
    COALESCE(a.nl_summary, '') AS nl_summary,
    a.last_analyzed_at
FROM merged m
JOIN wallet w ON w.id = m.wallet_id
LEFT JOIN latest_score s ON s.wallet_id = w.id
LEFT JOIN latest_ai a ON a.wallet_id = w.id
WHERE m.trade_count >= ?
  AND (m.trading_pnl + m.maker_rebates) > ?
  AND (? = '' OR COALESCE(s.strategy_type, 'unknown') = ?)
  AND (? = '' OR COALESCE(s.pool_tier, 'observation') = ?)
  AND (
    ? = '' OR
    (? = 'true' AND a.last_analyzed_at IS NOT NULL) OR
    (? = 'false' AND a.last_analyzed_at IS NULL)
  )
ORDER BY %s, w.id ASC
LIMIT ? OFFSET ?`, orderClause)
	rows := make([]PotentialWalletRow, 0, f.Limit)
	err := r.db.WithContext(ctx).Raw(query,
		f.MinTrades,
		f.MinRealizedPnL,
		strings.TrimSpace(f.StrategyType),
		strings.TrimSpace(f.StrategyType),
		strings.TrimSpace(f.PoolTier),
		strings.TrimSpace(f.PoolTier),
		hasAI,
		hasAI,
		hasAI,
		f.Limit,
		f.Offset,
	).Scan(&rows).Error
	return rows, err
}

func (r *WalletRepository) CountNewPotentialWallets24h(
	ctx context.Context,
	minTrades int64,
	minRealizedPnL float64,
) (int64, error) {
	if minTrades <= 0 {
		minTrades = 100
	}
	var total int64
	err := r.db.WithContext(ctx).Raw(`
WITH taker AS (
    SELECT
        taker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE
            WHEN side = 0 THEN (price * size) - fee_paid
            WHEN side = 1 THEN -((price * size) + fee_paid)
            ELSE 0
        END), 0) AS trading_pnl
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
    GROUP BY taker_wallet_id
),
maker AS (
    SELECT
        maker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
    GROUP BY maker_wallet_id
),
merged AS (
    SELECT
        x.wallet_id,
        SUM(x.trade_count) AS trade_count,
        SUM(x.trading_pnl) AS trading_pnl,
        SUM(x.maker_rebates) AS maker_rebates
    FROM (
        SELECT wallet_id, trade_count, trading_pnl, 0::numeric AS maker_rebates FROM taker
        UNION ALL
        SELECT wallet_id, trade_count, 0::numeric AS trading_pnl, maker_rebates FROM maker
    ) x
    GROUP BY x.wallet_id
)
SELECT COUNT(*)
FROM merged m
JOIN wallet w ON w.id = m.wallet_id
WHERE w.first_seen_at > NOW() - INTERVAL '24 hours'
  AND m.trade_count >= ?
  AND (m.trading_pnl + m.maker_rebates) > ?`, minTrades, minRealizedPnL).Scan(&total).Error
	return total, err
}

func (r *WalletRepository) ListOpsTopRealizedWallets(
	ctx context.Context,
	limit int,
	minTrades int64,
	minRealizedPnL float64,
) ([]OpsTopRealizedRow, error) {
	if limit <= 0 {
		limit = 5
	}
	if minTrades <= 0 {
		minTrades = 100
	}
	rows := make([]OpsTopRealizedRow, 0, limit)
	err := r.db.WithContext(ctx).Raw(`
WITH total_taker AS (
    SELECT
        taker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE
            WHEN side = 0 THEN (price * size) - fee_paid
            WHEN side = 1 THEN -((price * size) + fee_paid)
            ELSE 0
        END), 0) AS trading_pnl
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
    GROUP BY taker_wallet_id
),
total_maker AS (
    SELECT
        maker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
    GROUP BY maker_wallet_id
),
total_merged AS (
    SELECT
        x.wallet_id,
        SUM(x.trade_count) AS trade_count,
        SUM(x.trading_pnl) AS trading_pnl,
        SUM(x.maker_rebates) AS maker_rebates
    FROM (
        SELECT wallet_id, trade_count, trading_pnl, 0::numeric AS maker_rebates FROM total_taker
        UNION ALL
        SELECT wallet_id, trade_count, 0::numeric AS trading_pnl, maker_rebates FROM total_maker
    ) x
    GROUP BY x.wallet_id
),
day_taker AS (
    SELECT
        taker_wallet_id AS wallet_id,
        COALESCE(SUM(CASE
            WHEN side = 0 THEN (price * size) - fee_paid
            WHEN side = 1 THEN -((price * size) + fee_paid)
            ELSE 0
        END), 0) AS trading_pnl_24h
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
      AND block_time > NOW() - INTERVAL '24 hours'
    GROUP BY taker_wallet_id
),
day_maker AS (
    SELECT
        maker_wallet_id AS wallet_id,
        COALESCE(SUM(CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates_24h
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
      AND block_time > NOW() - INTERVAL '24 hours'
    GROUP BY maker_wallet_id
),
day_merged AS (
    SELECT
        x.wallet_id,
        SUM(x.trading_pnl_24h) AS trading_pnl_24h,
        SUM(x.maker_rebates_24h) AS maker_rebates_24h
    FROM (
        SELECT wallet_id, trading_pnl_24h, 0::numeric AS maker_rebates_24h FROM day_taker
        UNION ALL
        SELECT wallet_id, 0::numeric AS trading_pnl_24h, maker_rebates_24h FROM day_maker
    ) x
    GROUP BY x.wallet_id
),
latest_ai AS (
    SELECT DISTINCT ON (wallet_id)
        wallet_id,
        model_id,
        nl_summary,
        created_at
    FROM ai_analysis_report
    ORDER BY wallet_id, created_at DESC
)
SELECT
    w.id AS wallet_id,
    w.address AS address,
    t.trade_count,
    (t.trading_pnl + t.maker_rebates) AS realized_pnl,
    (COALESCE(d.trading_pnl_24h, 0) + COALESCE(d.maker_rebates_24h, 0)) AS realized_pnl_24h,
    (a.wallet_id IS NOT NULL) AS has_ai_report,
    COALESCE(a.nl_summary, '') AS latest_summary,
    COALESCE(a.model_id, '') AS latest_model_id,
    a.created_at AS last_analyzed_at
FROM total_merged t
JOIN wallet w ON w.id = t.wallet_id
LEFT JOIN day_merged d ON d.wallet_id = t.wallet_id
LEFT JOIN latest_ai a ON a.wallet_id = t.wallet_id
WHERE t.trade_count >= ?
  AND (t.trading_pnl + t.maker_rebates) > ?
ORDER BY realized_pnl_24h DESC, realized_pnl DESC, t.trade_count DESC, w.id ASC
LIMIT ?`, minTrades, minRealizedPnL, limit).Scan(&rows).Error
	return rows, err
}

func (r *WalletRepository) ListOpsTopAIConfidenceWallets(
	ctx context.Context,
	limit int,
	minTrades int64,
	minRealizedPnL float64,
) ([]OpsTopAIConfidenceRow, error) {
	if limit <= 0 {
		limit = 5
	}
	if minTrades <= 0 {
		minTrades = 100
	}
	rows := make([]OpsTopAIConfidenceRow, 0, limit)
	err := r.db.WithContext(ctx).Raw(`
WITH taker AS (
    SELECT
        taker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE
            WHEN side = 0 THEN (price * size) - fee_paid
            WHEN side = 1 THEN -((price * size) + fee_paid)
            ELSE 0
        END), 0) AS trading_pnl
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
    GROUP BY taker_wallet_id
),
maker AS (
    SELECT
        maker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
    GROUP BY maker_wallet_id
),
merged AS (
    SELECT
        x.wallet_id,
        SUM(x.trade_count) AS trade_count,
        SUM(x.trading_pnl) AS trading_pnl,
        SUM(x.maker_rebates) AS maker_rebates
    FROM (
        SELECT wallet_id, trade_count, trading_pnl, 0::numeric AS maker_rebates FROM taker
        UNION ALL
        SELECT wallet_id, trade_count, 0::numeric AS trading_pnl, maker_rebates FROM maker
    ) x
    GROUP BY x.wallet_id
),
latest_score AS (
    SELECT DISTINCT ON (wallet_id)
        wallet_id, smart_score, info_edge_level, strategy_type
    FROM wallet_score
    ORDER BY wallet_id, scored_at DESC
),
latest_ai AS (
    SELECT DISTINCT ON (wallet_id)
        wallet_id, nl_summary, created_at
    FROM ai_analysis_report
    ORDER BY wallet_id, created_at DESC
)
SELECT
    w.id AS wallet_id,
    w.address AS address,
    m.trade_count,
    (m.trading_pnl + m.maker_rebates) AS realized_pnl,
    COALESCE(s.smart_score, 0) AS smart_score,
    COALESCE(s.info_edge_level, 'unknown') AS info_edge_level,
    COALESCE(s.strategy_type, 'unknown') AS strategy_type,
    COALESCE(a.nl_summary, '') AS nl_summary,
    a.created_at AS last_analyzed_at
FROM merged m
JOIN wallet w ON w.id = m.wallet_id
JOIN latest_score s ON s.wallet_id = w.id
JOIN latest_ai a ON a.wallet_id = w.id
WHERE m.trade_count >= ?
  AND (m.trading_pnl + m.maker_rebates) > ?
ORDER BY s.smart_score DESC, m.trade_count DESC, (m.trading_pnl + m.maker_rebates) DESC, w.id ASC
LIMIT ?`, minTrades, minRealizedPnL, limit).Scan(&rows).Error
	return rows, err
}

func (r *WalletRepository) EnsureByAddress(ctx context.Context, addressHex string) (*model.Wallet, error) {
	address, err := polyaddr.HexToBytes(addressHex)
	if err != nil {
		return nil, err
	}
	wallet := model.Wallet{Address: address, ChainID: 137, FirstSeenAt: time.Now().UTC()}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}, {Name: "chain_id"}},
		DoNothing: true,
	}).Create(&wallet).Error; err != nil {
		return nil, err
	}
	var row model.Wallet
	if err := r.db.WithContext(ctx).Where("address = ? AND chain_id = ?", address, 137).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *MarketRepository) List(ctx context.Context, f MarketListFilter) ([]model.Market, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Market{})
	if f.Category != "" {
		q = q.Where("category = ?", f.Category)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := buildOrderClause(f.SortBy, f.Order, map[string]struct{}{
		"id": {}, "created_at": {}, "updated_at": {}, "volume": {}, "liquidity": {},
	})

	var rows []model.Market
	err := q.Order(orderClause).Limit(f.Limit).Offset(f.Offset).Find(&rows).Error
	return rows, total, err
}

func (r *MarketRepository) GetByID(ctx context.Context, id int64) (*model.Market, error) {
	var row model.Market
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *MarketRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Market{}).Count(&total).Error
	return total, err
}

func (r *MarketRepository) UpsertMany(ctx context.Context, markets []model.Market) error {
	if len(markets) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "condition_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"slug":       gorm.Expr("EXCLUDED.slug"),
			"title":      gorm.Expr("EXCLUDED.title"),
			"category":   gorm.Expr("EXCLUDED.category"),
			"status":     gorm.Expr("EXCLUDED.status"),
			"has_fee":    gorm.Expr("EXCLUDED.has_fee"),
			"volume":     gorm.Expr("EXCLUDED.volume"),
			"liquidity":  gorm.Expr("EXCLUDED.liquidity"),
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).Create(&markets).Error
}

func (r *MarketRepository) EnsureByConditionID(ctx context.Context, conditionID string) (*model.Market, error) {
	var row model.Market
	err := r.db.WithContext(ctx).Where("condition_id = ?", conditionID).First(&row).Error
	if err == nil {
		return &row, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	row = model.Market{
		ConditionID: conditionID,
		Title:       fmt.Sprintf("market-%s", conditionID),
		Status:      0,
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *TokenRepository) EnsureToken(ctx context.Context, marketID int64, tokenID string, side int16) (*model.Token, error) {
	var row model.Token
	err := r.db.WithContext(ctx).Where("token_id = ?", tokenID).First(&row).Error
	if err == nil {
		return &row, nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	row = model.Token{MarketID: marketID, TokenID: tokenID, Side: side}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *TradeRepository) Upsert(ctx context.Context, row model.TradeFill) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "uniq_key"}},
		DoUpdates: clause.Assignments(map[string]any{
			"price":        row.Price,
			"size":         row.Size,
			"fee_paid":     row.FeePaid,
			"block_time":   row.BlockTime,
			"block_number": row.BlockNumber,
		}),
	}).Create(&row).Error
}

func (r *TradeRepository) AggregateByWalletID(ctx context.Context, walletID int64) (*WalletPnLSummary, error) {
	var out WalletPnLSummary
	err := r.db.WithContext(ctx).Raw(`
SELECT
  COALESCE(SUM(CASE
      WHEN taker_wallet_id = ? AND side = 0 THEN (price * size) - fee_paid
      WHEN taker_wallet_id = ? AND side = 1 THEN -((price * size) + fee_paid)
      ELSE 0
  END), 0) AS trading_pnl,
  COALESCE(SUM(CASE WHEN maker_wallet_id = ? AND fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates,
  COALESCE(SUM(CASE WHEN taker_wallet_id = ? THEN fee_paid ELSE 0 END), 0) AS fees_paid,
  COALESCE(COUNT(*), 0) AS total_trades,
  COALESCE(SUM(CASE WHEN block_time > NOW() - INTERVAL '30 days' THEN price * size ELSE 0 END), 0) AS volume_30d
FROM trade_fill
WHERE taker_wallet_id = ? OR maker_wallet_id = ?`, walletID, walletID, walletID, walletID, walletID, walletID).Scan(&out).Error
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *TradeRepository) TimingSummaryByWalletID(ctx context.Context, walletID int64) (*WalletTimingSummary, error) {
	var out WalletTimingSummary
	err := r.db.WithContext(ctx).Raw(`
WITH wallet_trades AS (
    SELECT tf.block_time, t.market_id
    FROM trade_fill tf
    JOIN token t ON t.id = tf.token_id
    WHERE tf.taker_wallet_id = ? OR tf.maker_wallet_id = ?
),
paired AS (
    SELECT
      wt.block_time,
      oe.event_time
    FROM wallet_trades wt
    JOIN LATERAL (
      SELECT event_time
      FROM offchain_event
      WHERE market_id = wt.market_id AND event_time <= wt.block_time
      ORDER BY event_time DESC
      LIMIT 1
    ) oe ON true
)
SELECT
  COALESCE(AVG(EXTRACT(EPOCH FROM (block_time - event_time)) / 60.0), 0) AS mean_delta_minutes,
  COALESCE(STDDEV_POP(EXTRACT(EPOCH FROM (block_time - event_time)) / 60.0), 0) AS stddev_minutes,
  COALESCE(COUNT(*), 0) AS samples
FROM paired`, walletID, walletID).Scan(&out).Error
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ScoreRepository) LatestByWalletID(ctx context.Context, walletID int64) (*model.WalletScore, error) {
	var row model.WalletScore
	err := r.db.WithContext(ctx).Where("wallet_id = ?", walletID).Order("scored_at desc").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *ScoreRepository) Leaderboard(ctx context.Context, limit int, offset int, sortBy string, order string) ([]LeaderboardRow, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	allowed := map[string]struct{}{"smart_score": {}, "scored_at": {}}
	orderClause := buildOrderClause(sortBy, order, allowed)
	query := fmt.Sprintf(`
SELECT ws.wallet_id, w.address, w.pseudonym, ws.strategy_type, ws.smart_score, ws.info_edge_level, to_char(ws.scored_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS scored_at
FROM wallet_score ws
JOIN wallet w ON w.id = ws.wallet_id
JOIN (
    SELECT wallet_id, MAX(scored_at) AS max_scored_at
    FROM wallet_score
    GROUP BY wallet_id
) latest ON latest.wallet_id = ws.wallet_id AND latest.max_scored_at = ws.scored_at
ORDER BY %s
LIMIT ? OFFSET ?`, orderClause)

	var rows []LeaderboardRow
	if err := r.db.WithContext(ctx).Raw(query, limit, offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := r.db.WithContext(ctx).Raw(`
SELECT COUNT(*) FROM (
    SELECT wallet_id
    FROM wallet_score
    GROUP BY wallet_id
) x`).Scan(&total).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *AIReportRepository) LatestByWalletID(ctx context.Context, walletID int64) (*model.AIAnalysisReport, error) {
	var row model.AIAnalysisReport
	err := r.db.WithContext(ctx).Where("wallet_id = ?", walletID).Order("created_at desc").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *AIReportRepository) Create(ctx context.Context, report *model.AIAnalysisReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *AIReportRepository) ListByWalletID(ctx context.Context, walletID int64, limit int) ([]model.AIAnalysisReport, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	var rows []model.AIAnalysisReport
	err := r.db.WithContext(ctx).Where("wallet_id = ?", walletID).Order("created_at desc").Limit(limit).Find(&rows).Error
	return rows, err
}

func (r *WatchlistRepository) Add(ctx context.Context, walletID int64, userFingerprint string) error {
	entry := model.Watchlist{
		WalletID:        walletID,
		UserFingerprint: strings.TrimSpace(userFingerprint),
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "wallet_id"}, {Name: "user_fingerprint"}},
		DoNothing: true,
	}).Create(&entry).Error
}

func (r *WatchlistRepository) Remove(ctx context.Context, walletID int64, userFingerprint string) error {
	return r.db.WithContext(ctx).
		Where("wallet_id = ? AND user_fingerprint = ?", walletID, strings.TrimSpace(userFingerprint)).
		Delete(&model.Watchlist{}).Error
}

func (r *WatchlistRepository) CountByUser(ctx context.Context, userFingerprint string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Watchlist{}).
		Where("user_fingerprint = ?", strings.TrimSpace(userFingerprint)).
		Count(&total).Error
	return total, err
}

func (r *WatchlistRepository) CountByWallet(ctx context.Context, walletID int64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Watchlist{}).
		Where("wallet_id = ?", walletID).
		Count(&total).Error
	return total, err
}

func (r *WatchlistRepository) CountByWalletSince(ctx context.Context, walletID int64, since time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Watchlist{}).
		Where("wallet_id = ? AND created_at >= ?", walletID, since.UTC()).
		Count(&total).Error
	return total, err
}

func (r *WatchlistRepository) CountActionRequiredByUser(ctx context.Context, userFingerprint string) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(`
SELECT COUNT(*)
FROM (
    SELECT e.id AS event_id
    FROM watchlist wl
    JOIN wallet_update_event e ON e.wallet_id = wl.wallet_id
    WHERE wl.user_fingerprint = ?
      AND e.action_required = TRUE
    UNION ALL
    SELECT a.id AS event_id
    FROM watchlist wl
    JOIN anomaly_alert a ON a.wallet_id = wl.wallet_id
    WHERE wl.user_fingerprint = ?
      AND a.severity >= 2
) x`, strings.TrimSpace(userFingerprint), strings.TrimSpace(userFingerprint)).Scan(&total).Error
	return total, err
}

func (r *WatchlistRepository) ListByUser(ctx context.Context, userFingerprint string, limit int, offset int) ([]WatchlistWalletRow, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	rows := make([]WatchlistWalletRow, 0, limit)
	err := r.db.WithContext(ctx).Raw(`
WITH taker AS (
    SELECT
        taker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE
            WHEN side = 0 THEN (price * size) - fee_paid
            WHEN side = 1 THEN -((price * size) + fee_paid)
            ELSE 0
        END), 0) AS trading_pnl
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL
    GROUP BY taker_wallet_id
),
maker AS (
    SELECT
        maker_wallet_id AS wallet_id,
        COUNT(*) AS trade_count,
        COALESCE(SUM(CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END), 0) AS maker_rebates
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
    GROUP BY maker_wallet_id
),
merged AS (
    SELECT
        x.wallet_id,
        SUM(x.trade_count) AS trade_count,
        SUM(x.trading_pnl) AS trading_pnl,
        SUM(x.maker_rebates) AS maker_rebates
    FROM (
        SELECT wallet_id, trade_count, trading_pnl, 0::numeric AS maker_rebates FROM taker
        UNION ALL
        SELECT wallet_id, trade_count, 0::numeric AS trading_pnl, maker_rebates FROM maker
    ) x
    GROUP BY x.wallet_id
),
latest_score AS (
    SELECT DISTINCT ON (wallet_id)
        wallet_id, smart_score, info_edge_level, strategy_type, pool_tier
    FROM wallet_score
    ORDER BY wallet_id, scored_at DESC
),
latest_ai AS (
    SELECT DISTINCT ON (wallet_id)
        wallet_id, created_at AS last_analyzed_at, nl_summary
    FROM ai_analysis_report
    ORDER BY wallet_id, created_at DESC
)
SELECT
    wl.id AS watchlist_id,
    wl.created_at AS watchlisted_at,
    w.id AS wallet_id,
    w.address AS address,
    w.pseudonym AS pseudonym,
    w.is_tracked AS is_tracked,
    COALESCE(m.trade_count, 0) AS trade_count,
    COALESCE(m.trading_pnl, 0) AS trading_pnl,
    COALESCE(m.maker_rebates, 0) AS maker_rebates,
    (COALESCE(m.trading_pnl, 0) + COALESCE(m.maker_rebates, 0)) AS realized_pnl,
    COALESCE(s.smart_score, 0) AS smart_score,
    COALESCE(s.info_edge_level, 'unknown') AS info_edge_level,
    COALESCE(s.strategy_type, 'unknown') AS strategy_type,
    COALESCE(s.pool_tier, 'observation') AS pool_tier,
    COALESCE(a.nl_summary, '') AS nl_summary,
    a.last_analyzed_at
FROM watchlist wl
JOIN wallet w ON w.id = wl.wallet_id
LEFT JOIN merged m ON m.wallet_id = w.id
LEFT JOIN latest_score s ON s.wallet_id = w.id
LEFT JOIN latest_ai a ON a.wallet_id = w.id
WHERE wl.user_fingerprint = ?
ORDER BY wl.created_at DESC, wl.id DESC
LIMIT ? OFFSET ?`, strings.TrimSpace(userFingerprint), limit, offset).Scan(&rows).Error
	return rows, err
}

func (r *WatchlistRepository) CountFeedByUser(ctx context.Context, userFingerprint string, eventType string) (int64, error) {
	var total int64
	query := `
SELECT COUNT(*)
FROM (
    SELECT e.id AS event_id, e.wallet_id, e.event_type, e.created_at
    FROM watchlist wl
    JOIN wallet_update_event e ON e.wallet_id = wl.wallet_id
    WHERE wl.user_fingerprint = ?
    UNION ALL
    SELECT a.id AS event_id, a.wallet_id,
      CASE WHEN a.alert_type = 'pnl_spike' THEN 'pnl_jump' ELSE 'anomaly' END AS event_type,
      a.created_at
    FROM watchlist wl
    JOIN anomaly_alert a ON a.wallet_id = wl.wallet_id
    WHERE wl.user_fingerprint = ?
) x
WHERE (? = '' OR x.event_type = ?)
`
	err := r.db.WithContext(ctx).Raw(query,
		strings.TrimSpace(userFingerprint),
		strings.TrimSpace(userFingerprint),
		strings.TrimSpace(eventType),
		strings.TrimSpace(eventType),
	).Scan(&total).Error
	return total, err
}

func (r *WatchlistRepository) ListFeedByUser(ctx context.Context, userFingerprint string, eventType string, limit int, offset int) ([]WatchlistFeedRow, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	rows := make([]WatchlistFeedRow, 0, limit)
	err := r.db.WithContext(ctx).Raw(`
SELECT
    x.event_id,
    x.wallet_id,
    w.address,
    w.pseudonym,
    x.event_type,
    x.event_payload,
    x.action_required,
    x.suggestion,
    x.suggestion_zh,
    x.event_time
FROM (
    SELECT
        e.id AS event_id,
        e.wallet_id,
        e.event_type,
        e.payload AS event_payload,
        e.action_required,
        e.suggestion,
        e.suggestion_zh,
        e.created_at AS event_time
    FROM watchlist wl
    JOIN wallet_update_event e ON e.wallet_id = wl.wallet_id
    WHERE wl.user_fingerprint = ?
    UNION ALL
    SELECT
        a.id AS event_id,
        a.wallet_id,
        CASE WHEN a.alert_type = 'pnl_spike' THEN 'pnl_jump' ELSE 'anomaly' END AS event_type,
        a.evidence AS event_payload,
        (a.severity >= 2) AS action_required,
        CASE
          WHEN a.alert_type = 'pnl_spike' THEN 'PnL changed sharply. Review this wallet before increasing position.'
          WHEN a.alert_type = 'pre_event_timing' THEN 'Timing pattern changed. Re-check if strategy still matches your preference.'
          ELSE 'Anomaly detected. Consider reducing exposure until signals stabilize.'
        END AS suggestion,
        CASE
          WHEN a.alert_type = 'pnl_spike' THEN '收益出现明显变化，建议先复核该钱包，再考虑加仓。'
          WHEN a.alert_type = 'pre_event_timing' THEN '时序模式发生变化，建议检查该策略是否仍符合你的偏好。'
          ELSE '检测到异常，建议在信号稳定前控制仓位。'
        END AS suggestion_zh,
        a.created_at AS event_time
    FROM watchlist wl
    JOIN anomaly_alert a ON a.wallet_id = wl.wallet_id
    WHERE wl.user_fingerprint = ?
) x
JOIN wallet w ON w.id = x.wallet_id
WHERE (? = '' OR x.event_type = ?)
ORDER BY x.event_time DESC, x.event_id DESC
LIMIT ? OFFSET ?`,
		strings.TrimSpace(userFingerprint),
		strings.TrimSpace(userFingerprint),
		strings.TrimSpace(eventType),
		strings.TrimSpace(eventType),
		limit,
		offset,
	).Scan(&rows).Error
	return rows, err
}

func (r *WatchlistRepository) CreateUpdateEvent(ctx context.Context, walletID int64, eventType string, payload json.RawMessage) error {
	return r.CreateUpdateEventWithAdvice(ctx, walletID, eventType, payload, false, "", "")
}

func (r *WatchlistRepository) CreateUpdateEventWithAdvice(
	ctx context.Context,
	walletID int64,
	eventType string,
	payload json.RawMessage,
	actionRequired bool,
	suggestion string,
	suggestionZh string,
) error {
	var suggestionPtr *string
	if strings.TrimSpace(suggestion) != "" {
		v := strings.TrimSpace(suggestion)
		suggestionPtr = &v
	}
	var suggestionZhPtr *string
	if strings.TrimSpace(suggestionZh) != "" {
		v := strings.TrimSpace(suggestionZh)
		suggestionZhPtr = &v
	}
	event := model.WalletUpdateEvent{
		WalletID:       walletID,
		EventType:      strings.TrimSpace(eventType),
		Payload:        datatypes.JSON(payload),
		ActionRequired: actionRequired,
		Suggestion:     suggestionPtr,
		SuggestionZh:   suggestionZhPtr,
	}
	return r.db.WithContext(ctx).Create(&event).Error
}

func (r *WatchlistRepository) AddBatch(ctx context.Context, walletIDs []int64, userFingerprint string) ([]int64, error) {
	cleanUser := strings.TrimSpace(userFingerprint)
	if cleanUser == "" || len(walletIDs) == 0 {
		return nil, nil
	}
	created := make([]int64, 0, len(walletIDs))
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, walletID := range walletIDs {
			if walletID <= 0 {
				continue
			}
			entry := model.Watchlist{WalletID: walletID, UserFingerprint: cleanUser}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "wallet_id"}, {Name: "user_fingerprint"}},
				DoNothing: true,
			}).Create(&entry).Error; err != nil {
				return err
			}
			created = append(created, walletID)
		}
		return nil
	})
	return created, err
}

func (r *WalletRepository) ListRecentEventsByWalletID(ctx context.Context, walletID int64, limit int) ([]WalletEventRow, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	rows := make([]WalletEventRow, 0, limit)
	err := r.db.WithContext(ctx).Raw(`
SELECT
    x.event_id,
    x.event_type,
    x.event_payload,
    x.action_required,
    x.suggestion,
    x.suggestion_zh,
    x.event_time
FROM (
    SELECT
        e.id AS event_id,
        e.event_type,
        e.payload AS event_payload,
        e.action_required,
        e.suggestion,
        e.suggestion_zh,
        e.created_at AS event_time
    FROM wallet_update_event e
    WHERE e.wallet_id = ?
    UNION ALL
    SELECT
        a.id AS event_id,
        CASE WHEN a.alert_type = 'pnl_spike' THEN 'pnl_jump' ELSE 'anomaly' END AS event_type,
        a.evidence AS event_payload,
        (a.severity >= 2) AS action_required,
        CASE
          WHEN a.alert_type = 'pnl_spike' THEN 'PnL changed sharply. Review before adding exposure.'
          WHEN a.alert_type = 'pre_event_timing' THEN 'Timing edge changed. Re-check strategy fit.'
          ELSE 'Anomaly detected. Consider risk control actions first.'
        END AS suggestion,
        CASE
          WHEN a.alert_type = 'pnl_spike' THEN '收益出现明显变化，建议先复核再加仓。'
          WHEN a.alert_type = 'pre_event_timing' THEN '时序优势变化，建议重新评估策略匹配。'
          ELSE '检测到异常，建议优先进行风险控制。'
        END AS suggestion_zh,
        a.created_at AS event_time
    FROM anomaly_alert a
    WHERE a.wallet_id = ?
) x
ORDER BY x.event_time DESC, x.event_id DESC
LIMIT ?`, walletID, walletID, limit).Scan(&rows).Error
	return rows, err
}

func (r *PortfolioRepository) ListActive(ctx context.Context) ([]model.Portfolio, error) {
	var rows []model.Portfolio
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("id ASC").
		Find(&rows).Error
	return rows, err
}

func buildOrderClause(sortBy string, order string, allow map[string]struct{}) string {
	sortBy = strings.TrimSpace(sortBy)
	if _, ok := allow[sortBy]; !ok {
		if _, ok := allow["updated_at"]; ok {
			sortBy = "updated_at"
		} else if _, ok := allow["smart_score"]; ok {
			sortBy = "smart_score"
		} else {
			sortBy = "id"
		}
	}
	order = strings.ToLower(strings.TrimSpace(order))
	if order != "asc" {
		order = "desc"
	}
	return fmt.Sprintf("%s %s", sortBy, order)
}
