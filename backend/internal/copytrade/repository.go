package copytrade

import (
	"context"
	"time"

	"easy-arbitra/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// ── Config CRUD ──

func (r *Repository) UpsertConfig(ctx context.Context, cfg *model.CopyTradingConfig) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_fingerprint"}, {Name: "wallet_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"enabled":           cfg.Enabled,
			"max_position_usdc": cfg.MaxPositionUSDC,
			"risk_preference":   cfg.RiskPreference,
			"updated_at":        gorm.Expr("NOW()"),
		}),
	}).Create(cfg).Error
}

func (r *Repository) GetConfig(ctx context.Context, userFP string, walletID int64) (*model.CopyTradingConfig, error) {
	var row model.CopyTradingConfig
	err := r.db.WithContext(ctx).Where("user_fingerprint = ? AND wallet_id = ?", userFP, walletID).First(&row).Error
	return &row, err
}

func (r *Repository) GetConfigByID(ctx context.Context, id int64) (*model.CopyTradingConfig, error) {
	var row model.CopyTradingConfig
	err := r.db.WithContext(ctx).First(&row, id).Error
	return &row, err
}

func (r *Repository) ListConfigsByUser(ctx context.Context, userFP string) ([]model.CopyTradingConfig, error) {
	var rows []model.CopyTradingConfig
	err := r.db.WithContext(ctx).Where("user_fingerprint = ?", userFP).Order("created_at desc").Find(&rows).Error
	return rows, err
}

func (r *Repository) ListEnabledConfigs(ctx context.Context) ([]model.CopyTradingConfig, error) {
	var rows []model.CopyTradingConfig
	err := r.db.WithContext(ctx).Where("enabled = true").Find(&rows).Error
	return rows, err
}

func (r *Repository) UpdateLastChecked(ctx context.Context, configID int64, t time.Time) error {
	return r.db.WithContext(ctx).Model(&model.CopyTradingConfig{}).Where("id = ?", configID).Update("last_checked_at", t).Error
}

func (r *Repository) DisableConfig(ctx context.Context, userFP string, walletID int64) error {
	return r.db.WithContext(ctx).Model(&model.CopyTradingConfig{}).
		Where("user_fingerprint = ? AND wallet_id = ?", userFP, walletID).
		Updates(map[string]any{"enabled": false, "updated_at": gorm.Expr("NOW()")}).Error
}

// ── Decision CRUD ──

func (r *Repository) CreateDecision(ctx context.Context, d *model.CopyTradeDecision) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *Repository) GetDecision(ctx context.Context, id int64) (*model.CopyTradeDecision, error) {
	var row model.CopyTradeDecision
	err := r.db.WithContext(ctx).First(&row, id).Error
	return &row, err
}

func (r *Repository) ListDecisionsByConfig(ctx context.Context, configID int64, limit, offset int) ([]model.CopyTradeDecision, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.CopyTradeDecision{}).Where("config_id = ?", configID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.CopyTradeDecision
	err := q.Order("created_at desc").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}

func (r *Repository) ListRecentDecisionsByUser(ctx context.Context, userFP string, limit int) ([]model.CopyTradeDecision, error) {
	var rows []model.CopyTradeDecision
	err := r.db.WithContext(ctx).Raw(`
SELECT d.*
FROM copy_trade_decision d
JOIN copy_trading_config c ON c.id = d.config_id
WHERE c.user_fingerprint = ?
ORDER BY d.created_at DESC
LIMIT ?`, userFP, limit).Scan(&rows).Error
	return rows, err
}

func (r *Repository) ListOpenPositionsByUser(ctx context.Context, userFP string) ([]model.CopyTradeDecision, error) {
	var rows []model.CopyTradeDecision
	err := r.db.WithContext(ctx).Raw(`
SELECT d.*
FROM copy_trade_decision d
JOIN copy_trading_config c ON c.id = d.config_id
WHERE c.user_fingerprint = ? AND d.decision = 'copy' AND d.status = 'executed' AND d.closed_at IS NULL
ORDER BY d.created_at DESC`, userFP).Scan(&rows).Error
	return rows, err
}

func (r *Repository) CloseDecision(ctx context.Context, id int64, closePrice float64, realizedPnL float64) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&model.CopyTradeDecision{}).Where("id = ?", id).Updates(map[string]any{
		"status":       "stopped",
		"closed_at":    now,
		"close_price":  closePrice,
		"realized_pnl": realizedPnL,
	}).Error
}

func (r *Repository) CountExecutedByConfig(ctx context.Context, configID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CopyTradeDecision{}).
		Where("config_id = ? AND decision = 'copy' AND status = 'executed' AND closed_at IS NULL", configID).
		Count(&count).Error
	return count, err
}

func (r *Repository) SumExposureByConfig(ctx context.Context, configID int64) (float64, error) {
	var sum *float64
	err := r.db.WithContext(ctx).Model(&model.CopyTradeDecision{}).
		Select("COALESCE(SUM(size_usdc), 0)").
		Where("config_id = ? AND decision = 'copy' AND status = 'executed' AND closed_at IS NULL", configID).
		Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	if sum == nil {
		return 0, nil
	}
	return *sum, nil
}

func (r *Repository) HasRecentCopyInMarket(ctx context.Context, configID int64, marketID int64, lookback time.Duration) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CopyTradeDecision{}).
		Where("config_id = ? AND market_id = ? AND decision = 'copy' AND created_at > ?", configID, marketID, time.Now().UTC().Add(-lookback)).
		Count(&count).Error
	return count > 0, err
}

// ── Performance ──

func (r *Repository) UpsertDailyPerf(ctx context.Context, p *model.CopyTradeDailyPerf) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "config_id"}, {Name: "perf_date"}},
		DoUpdates: clause.Assignments(map[string]any{
			"total_copies":   gorm.Expr("copy_trade_daily_perf.total_copies + EXCLUDED.total_copies"),
			"profitable":     gorm.Expr("copy_trade_daily_perf.profitable + EXCLUDED.profitable"),
			"total_pnl":      gorm.Expr("copy_trade_daily_perf.total_pnl + EXCLUDED.total_pnl"),
			"total_exposure":  gorm.Expr("copy_trade_daily_perf.total_exposure + EXCLUDED.total_exposure"),
			"skipped":        gorm.Expr("copy_trade_daily_perf.skipped + EXCLUDED.skipped"),
		}),
	}).Create(p).Error
}

func (r *Repository) ListDailyPerfByConfig(ctx context.Context, configID int64, limit int) ([]model.CopyTradeDailyPerf, error) {
	var rows []model.CopyTradeDailyPerf
	err := r.db.WithContext(ctx).Where("config_id = ?", configID).Order("perf_date desc").Limit(limit).Find(&rows).Error
	return rows, err
}

// ── Monitoring queries ──

type HourlyStat struct {
	Hour           time.Time `json:"hour"`
	Runs           int       `json:"runs"`
	WalletsChecked int       `json:"wallets_checked"`
	NewTrades      int       `json:"new_trades"`
	DecisionsCopy  int       `json:"decisions_copy"`
	DecisionsSkip  int       `json:"decisions_skip"`
	Errors         int       `json:"errors"`
}

type CopyableWallet struct {
	ID            int64   `json:"id"`
	Pseudonym     *string `json:"pseudonym,omitempty"`
	Address       string  `json:"address"`
	SmartScore    int     `json:"smart_score"`
	PoolTier      string  `json:"pool_tier"`
	StrategyType  string  `json:"strategy_type"`
	RiskLevel     string  `json:"risk_level"`
	Momentum      string  `json:"momentum"`
	Pnl30D        float64 `json:"pnl_30d"`
	TradeCount30D int     `json:"trade_count_30d"`
}

func (r *Repository) GetSyncerRunHistory(ctx context.Context, limit int) ([]model.IngestRun, error) {
	var rows []model.IngestRun
	err := r.db.WithContext(ctx).
		Where("job_name = ?", "copy_trade_syncer").
		Order("started_at DESC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *Repository) GetHourlySyncStats(ctx context.Context, hours int) ([]HourlyStat, error) {
	var rows []HourlyStat
	err := r.db.WithContext(ctx).Raw(`
SELECT date_trunc('hour', started_at) AS hour,
       COUNT(*)                                          AS runs,
       COALESCE(SUM((stats->>'wallets_checked')::int),0) AS wallets_checked,
       COALESCE(SUM((stats->>'new_trades')::int),0)      AS new_trades,
       COALESCE(SUM((stats->>'decisions_copy')::int),0)   AS decisions_copy,
       COALESCE(SUM((stats->>'decisions_skip')::int),0)   AS decisions_skip,
       COALESCE(SUM((stats->>'errors')::int),0)           AS errors
FROM ingest_run
WHERE job_name = 'copy_trade_syncer' AND started_at > NOW() - make_interval(hours => ?)
GROUP BY 1 ORDER BY 1 DESC`, hours).Scan(&rows).Error
	return rows, err
}

func (r *Repository) GetCopyableWallets(ctx context.Context, limit int) ([]CopyableWallet, error) {
	var rows []CopyableWallet
	err := r.db.WithContext(ctx).Raw(`
SELECT w.id, w.pseudonym, encode(w.address, 'hex') AS address,
       ws.smart_score, ws.pool_tier, ws.strategy_type, ws.risk_level, ws.momentum,
       wf.pnl_30d, wf.trade_count_30d
FROM wallet w
JOIN wallet_score ws ON ws.wallet_id = w.id
JOIN wallet_features_daily wf ON wf.wallet_id = w.id
WHERE ws.id = (SELECT id FROM wallet_score WHERE wallet_id = w.id ORDER BY scored_at DESC LIMIT 1)
  AND wf.feature_date = (SELECT MAX(feature_date) FROM wallet_features_daily WHERE wallet_id = w.id)
  AND ws.smart_score >= 70
  AND wf.pnl_30d > 0
  AND wf.trade_count_30d >= 10
ORDER BY ws.smart_score DESC, wf.pnl_30d DESC
LIMIT ?`, limit).Scan(&rows).Error
	return rows, err
}

// ── Dashboard aggregates ──

type DashboardStats struct {
	TotalPnL     float64 `gorm:"column:total_pnl"`
	TotalCopies  int64   `gorm:"column:total_copies"`
	TotalSkipped int64   `gorm:"column:total_skipped"`
	Profitable   int64   `gorm:"column:profitable"`
	OpenCount    int64   `gorm:"column:open_count"`
}

func (r *Repository) GetDashboardStats(ctx context.Context, userFP string) (*DashboardStats, error) {
	var out DashboardStats
	err := r.db.WithContext(ctx).Raw(`
SELECT
  COALESCE(SUM(CASE WHEN d.realized_pnl IS NOT NULL THEN d.realized_pnl ELSE 0 END), 0) AS total_pnl,
  COUNT(*) FILTER (WHERE d.decision = 'copy') AS total_copies,
  COUNT(*) FILTER (WHERE d.decision = 'skip') AS total_skipped,
  COUNT(*) FILTER (WHERE d.decision = 'copy' AND d.realized_pnl > 0) AS profitable,
  COUNT(*) FILTER (WHERE d.decision = 'copy' AND d.status = 'executed' AND d.closed_at IS NULL) AS open_count
FROM copy_trade_decision d
JOIN copy_trading_config c ON c.id = d.config_id
WHERE c.user_fingerprint = ?`, userFP).Scan(&out).Error
	return &out, err
}
