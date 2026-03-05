package repository

import (
	"context"
	"time"

	"easy-arbitra/backend/internal/model"

	"gorm.io/gorm"
)

type NovaLearningRepository struct {
	db *gorm.DB
}

func NewNovaLearningRepository(db *gorm.DB) *NovaLearningRepository {
	return &NovaLearningRepository{db: db}
}

// Create creates a new learning log entry
func (r *NovaLearningRepository) Create(ctx context.Context, log *model.NovaLearningLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// ListRecent returns recent learning logs
func (r *NovaLearningRepository) ListRecent(ctx context.Context, limit int) ([]model.NovaLearningLog, error) {
	if limit <= 0 {
		limit = 30
	}
	if limit > 90 {
		limit = 90
	}
	var logs []model.NovaLearningLog
	err := r.db.WithContext(ctx).
		Order("validation_date DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetByDate returns learning log for a specific date
func (r *NovaLearningRepository) GetByDate(ctx context.Context, date time.Time) (*model.NovaLearningLog, error) {
	var log model.NovaLearningLog
	err := r.db.WithContext(ctx).
		Where("validation_date = ?", date.Truncate(24*time.Hour)).
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetSuccessRate returns success rate for a time period
func (r *NovaLearningRepository) GetSuccessRate(ctx context.Context, days int) (float64, error) {
	if days <= 0 {
		days = 7
	}

	var result struct {
		Total   int64
		Success int64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN is_success THEN 1 ELSE 0 END) as success
		FROM nova_learning_log
		WHERE validation_date >= NOW() - INTERVAL '? days'
	`, days).Scan(&result).Error

	if err != nil || result.Total == 0 {
		return 0, err
	}

	return float64(result.Success) / float64(result.Total) * 100, nil
}

// GetStrategyAdjustments returns recent strategy adjustments
func (r *NovaLearningRepository) GetStrategyAdjustments(ctx context.Context, limit int) ([]model.NovaLearningLog, error) {
	if limit <= 0 {
		limit = 10
	}
	var logs []model.NovaLearningLog
	err := r.db.WithContext(ctx).
		Where("strategy_adjustment IS NOT NULL").
		Order("validation_date DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
