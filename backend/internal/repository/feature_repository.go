package repository

import (
	"context"
	"time"

	"easy-arbitra/backend/internal/model"
	"gorm.io/gorm"
)

type FeatureRepository struct{ db *gorm.DB }

func NewFeatureRepository(db *gorm.DB) *FeatureRepository { return &FeatureRepository{db: db} }

func (r *FeatureRepository) BuildDaily(ctx context.Context, featureDate time.Time) error {
	return r.db.WithContext(ctx).Exec(`
INSERT INTO wallet_features_daily (
  wallet_id, feature_date, pnl_7d, pnl_30d, pnl_90d, maker_ratio, trade_count_30d, active_days_30d, tx_frequency_per_day, avg_edge, created_at
)
SELECT
  stats.wallet_id,
  ?::date AS feature_date,
  stats.pnl_7d,
  stats.pnl_30d,
  stats.pnl_90d,
  CASE WHEN stats.trade_count_30d > 0 THEN stats.maker_trades_30d::numeric / stats.trade_count_30d::numeric ELSE 0 END AS maker_ratio,
  stats.trade_count_30d,
  stats.active_days_30d,
  CASE WHEN stats.active_days_30d > 0 THEN stats.trade_count_30d::numeric / stats.active_days_30d::numeric ELSE 0 END AS tx_frequency_per_day,
  stats.avg_edge,
  NOW()
FROM (
  SELECT
    wallet_id,
    COALESCE(SUM(pnl_component) FILTER (WHERE block_time > NOW() - INTERVAL '7 days'), 0) AS pnl_7d,
    COUNT(*) FILTER (WHERE block_time > NOW() - INTERVAL '30 days') AS trade_count_30d,
    COUNT(DISTINCT DATE(block_time)) FILTER (WHERE block_time > NOW() - INTERVAL '30 days') AS active_days_30d,
    COUNT(*) FILTER (WHERE is_maker = 1 AND block_time > NOW() - INTERVAL '30 days') AS maker_trades_30d,
    COALESCE(SUM(pnl_component) FILTER (WHERE block_time > NOW() - INTERVAL '30 days'), 0) AS pnl_30d,
    COALESCE(SUM(pnl_component) FILTER (WHERE block_time > NOW() - INTERVAL '90 days'), 0) AS pnl_90d,
    COALESCE(AVG(edge_component) FILTER (WHERE block_time > NOW() - INTERVAL '30 days'), 0) AS avg_edge
  FROM (
    SELECT
      taker_wallet_id AS wallet_id,
      block_time,
      0 AS is_maker,
      CASE WHEN side = 0 THEN (price * size) - fee_paid ELSE -((price * size) + fee_paid) END AS pnl_component,
      CASE WHEN side = 0 THEN (1 - price) ELSE (price - 0.5) END AS edge_component
    FROM trade_fill
    WHERE taker_wallet_id IS NOT NULL

    UNION ALL

    SELECT
      maker_wallet_id AS wallet_id,
      block_time,
      1 AS is_maker,
      CASE WHEN fee_paid < 0 THEN ABS(fee_paid) ELSE 0 END AS pnl_component,
      0 AS edge_component
    FROM trade_fill
    WHERE maker_wallet_id IS NOT NULL
  ) u
  GROUP BY wallet_id
) stats
ON CONFLICT (wallet_id, feature_date)
DO UPDATE SET
  pnl_7d = EXCLUDED.pnl_7d,
  pnl_30d = EXCLUDED.pnl_30d,
  pnl_90d = EXCLUDED.pnl_90d,
  maker_ratio = EXCLUDED.maker_ratio,
  trade_count_30d = EXCLUDED.trade_count_30d,
  active_days_30d = EXCLUDED.active_days_30d,
  tx_frequency_per_day = EXCLUDED.tx_frequency_per_day,
  avg_edge = EXCLUDED.avg_edge,
  created_at = NOW()`, featureDate.UTC().Format("2006-01-02")).Error
}

func (r *FeatureRepository) LatestByWalletID(ctx context.Context, walletID int64) (*model.WalletFeaturesDaily, error) {
	var row model.WalletFeaturesDaily
	err := r.db.WithContext(ctx).Where("wallet_id = ?", walletID).Order("feature_date desc").First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *FeatureRepository) ListByWalletID(ctx context.Context, walletID int64, limit int) ([]model.WalletFeaturesDaily, error) {
	if limit <= 0 {
		limit = 90
	}
	if limit > 365 {
		limit = 365
	}
	var rows []model.WalletFeaturesDaily
	err := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Order("feature_date desc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *ScoreRepository) UpsertLatest(ctx context.Context, row model.WalletScore) error {
	return r.db.WithContext(ctx).Create(&row).Error
}
