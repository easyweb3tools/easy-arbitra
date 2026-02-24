package repository

import (
	"context"
	"strings"

	"easy-arbitra/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OffchainEventRepository struct{ db *gorm.DB }

func NewOffchainEventRepository(db *gorm.DB) *OffchainEventRepository {
	return &OffchainEventRepository{db: db}
}

func (r *OffchainEventRepository) UpsertMany(ctx context.Context, events []model.OffchainEvent) error {
	if len(events) == 0 {
		return nil
	}
	for i := range events {
		events[i].Source = strings.TrimSpace(events[i].Source)
		events[i].SourceEventID = strings.TrimSpace(events[i].SourceEventID)
		if events[i].Source == "" {
			events[i].Source = "unknown"
		}
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "source_name"}, {Name: "source_event_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"market_id":  gorm.Expr("EXCLUDED.market_id"),
			"event_time": gorm.Expr("EXCLUDED.event_time"),
			"event_type": gorm.Expr("EXCLUDED.event_type"),
			"title":      gorm.Expr("EXCLUDED.title"),
			"payload":    gorm.Expr("EXCLUDED.payload"),
			"created_at": gorm.Expr("NOW()"),
		}),
	}).Create(&events).Error
}
