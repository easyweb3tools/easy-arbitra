package worker

import (
	"context"
	"time"

	"easy-arbitra/backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TraderStatsRefresher struct {
	logger    *zap.Logger
	db        *gorm.DB
	statsRepo *repository.TraderStatsRepository
}

func NewTraderStatsRefresher(logger *zap.Logger, db *gorm.DB) *TraderStatsRefresher {
	return &TraderStatsRefresher{
		logger:    logger,
		db:        db,
		statsRepo: repository.NewTraderStatsRepository(db),
	}
}

func (w *TraderStatsRefresher) Name() string {
	return "trader_stats_refresher"
}

func (w *TraderStatsRefresher) RunOnce(ctx context.Context) error {
	start := time.Now()

	// Get the last sync time from ingest_cursor
	var lastSyncTime time.Time
	var cursorValue string
	err := w.db.WithContext(ctx).
		Table("ingest_cursor").
		Select("cursor_value").
		Where("source = ? AND stream = ?", "trader_stats", "last_refresh").
		Scan(&cursorValue).Error

	if err == nil && cursorValue != "" {
		// Parse the cursor value as timestamp
		lastSyncTime, err = time.Parse(time.RFC3339, cursorValue)
		if err != nil {
			w.logger.Warn("failed to parse last sync time, doing full refresh", zap.Error(err))
			lastSyncTime = time.Time{}
		}
	}

	// If no cursor or very old (> 7 days), do full refresh
	if lastSyncTime.IsZero() || time.Since(lastSyncTime) > 7*24*time.Hour {
		w.logger.Info("performing full refresh of trader_stats")
		if err := w.statsRepo.RefreshFull(ctx); err != nil {
			return err
		}
	} else {
		// Incremental refresh
		w.logger.Info("performing incremental refresh of trader_stats", zap.Time("since", lastSyncTime))
		if err := w.statsRepo.RefreshIncremental(ctx, lastSyncTime); err != nil {
			return err
		}
	}

	// Update cursor
	now := time.Now().UTC()
	upsertCursor := `
INSERT INTO ingest_cursor (source, stream, cursor_value, updated_at)
VALUES (?, ?, ?, ?)
ON CONFLICT (source, stream) DO UPDATE SET
    cursor_value = EXCLUDED.cursor_value,
    updated_at = EXCLUDED.updated_at`

	if err := w.db.WithContext(ctx).Exec(upsertCursor, "trader_stats", "last_refresh", now.Format(time.RFC3339), now).Error; err != nil {
		w.logger.Warn("failed to update cursor", zap.Error(err))
	}

	// Get stats count
	count, _ := w.statsRepo.Count(ctx)

	w.logger.Info("trader_stats refresh completed",
		zap.Int64("wallets", count),
		zap.Duration("duration", time.Since(start)))

	return nil
}
