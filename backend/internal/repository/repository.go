package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/pkg/polyaddr"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct{ db *gorm.DB }

type MarketRepository struct{ db *gorm.DB }

type TokenRepository struct{ db *gorm.DB }

type TradeRepository struct{ db *gorm.DB }

type ScoreRepository struct{ db *gorm.DB }

type AIReportRepository struct{ db *gorm.DB }

type WalletListFilter struct {
	Tracked *bool
	Search  string
	SortBy  string
	Order   string
	Limit   int
	Offset  int
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
	TradingPnL   float64 `json:"trading_pnl"`
	MakerRebates float64 `json:"maker_rebates"`
	FeesPaid     float64 `json:"fees_paid"`
	TotalTrades  int64   `json:"total_trades"`
	Volume30D    float64 `json:"volume_30d"`
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

func NewWalletRepository(db *gorm.DB) *WalletRepository     { return &WalletRepository{db: db} }
func NewMarketRepository(db *gorm.DB) *MarketRepository     { return &MarketRepository{db: db} }
func NewTokenRepository(db *gorm.DB) *TokenRepository       { return &TokenRepository{db: db} }
func NewTradeRepository(db *gorm.DB) *TradeRepository       { return &TradeRepository{db: db} }
func NewScoreRepository(db *gorm.DB) *ScoreRepository       { return &ScoreRepository{db: db} }
func NewAIReportRepository(db *gorm.DB) *AIReportRepository { return &AIReportRepository{db: db} }

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
