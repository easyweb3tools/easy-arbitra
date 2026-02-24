package repository

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"easy-arbitra/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IngestCursorRepository struct{ db *gorm.DB }

type IngestRunRepository struct{ db *gorm.DB }

func NewIngestCursorRepository(db *gorm.DB) *IngestCursorRepository {
	return &IngestCursorRepository{db: db}
}
func NewIngestRunRepository(db *gorm.DB) *IngestRunRepository { return &IngestRunRepository{db: db} }

func (r *IngestCursorRepository) Get(ctx context.Context, source string, stream string) (*model.IngestCursor, error) {
	var row model.IngestCursor
	err := r.db.WithContext(ctx).Where("source = ? AND stream = ?", strings.TrimSpace(source), strings.TrimSpace(stream)).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *IngestCursorRepository) Upsert(ctx context.Context, source string, stream string, cursorValue string) error {
	row := model.IngestCursor{
		Source:      strings.TrimSpace(source),
		Stream:      strings.TrimSpace(stream),
		CursorValue: strings.TrimSpace(cursorValue),
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "source"}, {Name: "stream"}},
		DoUpdates: clause.Assignments(map[string]any{
			"cursor_value": row.CursorValue,
			"updated_at":   gorm.Expr("NOW()"),
		}),
	}).Create(&row).Error
}

func (r *IngestRunRepository) Start(ctx context.Context, jobName string) (int64, error) {
	row := model.IngestRun{
		JobName: strings.TrimSpace(jobName),
		Status:  "running",
		Stats:   []byte(`{}`),
	}
	if err := r.db.WithContext(ctx).Create(&row).Error; err != nil {
		return 0, err
	}
	return row.ID, nil
}

func (r *IngestRunRepository) Finish(ctx context.Context, id int64, status string, stats map[string]any, errText string) error {
	status = strings.TrimSpace(status)
	if status == "" {
		status = "done"
	}
	payload := []byte(`{}`)
	if stats != nil {
		if b, err := json.Marshal(stats); err == nil {
			payload = b
		}
	}
	updates := map[string]any{
		"status":   status,
		"stats":    payload,
		"ended_at": time.Now().UTC(),
	}
	if strings.TrimSpace(errText) != "" {
		updates["error_text"] = strings.TrimSpace(errText)
	}
	return r.db.WithContext(ctx).Model(&model.IngestRun{}).Where("id = ?", id).Updates(updates).Error
}
