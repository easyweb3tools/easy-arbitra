package repository

import (
	"context"

	"easy-arbitra/backend/internal/model"
	"gorm.io/gorm"
)

type AnomalyRepository struct{ db *gorm.DB }

type AnomalyListFilter struct {
	Severity     *int16
	AlertType    string
	Acknowledged *bool
	Limit        int
	Offset       int
}

func NewAnomalyRepository(db *gorm.DB) *AnomalyRepository { return &AnomalyRepository{db: db} }

func (r *AnomalyRepository) Create(ctx context.Context, alert *model.AnomalyAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *AnomalyRepository) List(ctx context.Context, f AnomalyListFilter) ([]model.AnomalyAlert, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.AnomalyAlert{})
	if f.Severity != nil {
		q = q.Where("severity = ?", *f.Severity)
	}
	if f.AlertType != "" {
		q = q.Where("alert_type = ?", f.AlertType)
	}
	if f.Acknowledged != nil {
		q = q.Where("acknowledged = ?", *f.Acknowledged)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
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
	var rows []model.AnomalyAlert
	err := q.Order("created_at desc").Limit(f.Limit).Offset(f.Offset).Find(&rows).Error
	return rows, total, err
}

func (r *AnomalyRepository) MarkAcknowledged(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.AnomalyAlert{}).Where("id = ?", id).Update("acknowledged", true).Error
}

func (r *AnomalyRepository) GetByID(ctx context.Context, id int64) (*model.AnomalyAlert, error) {
	var row model.AnomalyAlert
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *AnomalyRepository) ExistsRecent(ctx context.Context, walletID int64, alertType string, lookbackHours int) (bool, error) {
	if lookbackHours <= 0 {
		lookbackHours = 6
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AnomalyAlert{}).
		Where("wallet_id = ? AND alert_type = ? AND created_at > NOW() - (? || ' hours')::interval", walletID, alertType, lookbackHours).
		Count(&count).Error
	return count > 0, err
}
