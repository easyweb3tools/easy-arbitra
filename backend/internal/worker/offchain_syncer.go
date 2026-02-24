package worker

import (
	"context"
	"strconv"
	"strings"
	"time"

	"easy-arbitra/backend/internal/client"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	offchainCursorSource = "gamma_api"
	offchainCursorStream = "events_latest_ts"
)

type OffchainEventSyncer struct {
	client         *client.OffchainClient
	marketRepo     *repository.MarketRepository
	offchainRepo   *repository.OffchainEventRepository
	cursorRepo     *repository.IngestCursorRepository
	limit          int
	maxPages       int
	cursorLookback time.Duration
}

func NewOffchainEventSyncer(
	client *client.OffchainClient,
	marketRepo *repository.MarketRepository,
	offchainRepo *repository.OffchainEventRepository,
	cursorRepo *repository.IngestCursorRepository,
	limit int,
	maxPages int,
	cursorLookback time.Duration,
) *OffchainEventSyncer {
	if limit <= 0 {
		limit = 50
	}
	if maxPages <= 0 {
		maxPages = 5
	}
	if cursorLookback < 0 {
		cursorLookback = 0
	}
	return &OffchainEventSyncer{
		client:         client,
		marketRepo:     marketRepo,
		offchainRepo:   offchainRepo,
		cursorRepo:     cursorRepo,
		limit:          limit,
		maxPages:       maxPages,
		cursorLookback: cursorLookback,
	}
}

func (s *OffchainEventSyncer) Name() string { return "offchain_event_syncer" }

func (s *OffchainEventSyncer) RunOnce(ctx context.Context) error {
	cursorTs := int64(0)
	if s.cursorRepo != nil {
		cursor, err := s.cursorRepo.Get(ctx, offchainCursorSource, offchainCursorStream)
		if err == nil {
			if parsed, parseErr := strconv.ParseInt(strings.TrimSpace(cursor.CursorValue), 10, 64); parseErr == nil {
				cursorTs = parsed
			}
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}
	cutoffTs := cursorTs - int64(s.cursorLookback/time.Second)
	if cutoffTs < 0 {
		cutoffTs = 0
	}

	maxSeenTs := cursorTs
	offset := 0
	for page := 0; page < s.maxPages; page++ {
		events, err := s.client.FetchEvents(ctx, s.limit, offset)
		if err != nil {
			return err
		}
		if len(events) == 0 {
			break
		}

		rows := make([]model.OffchainEvent, 0, len(events))
		allOlderThanCutoff := true
		for _, event := range events {
			ts := event.Time.UTC().Unix()
			if ts > maxSeenTs {
				maxSeenTs = ts
			}
			if ts < cutoffTs {
				continue
			}
			allOlderThanCutoff = false

			var marketID *int64
			if event.ConditionID != "" {
				m, err := s.marketRepo.EnsureByConditionID(ctx, event.ConditionID)
				if err == nil && m != nil {
					marketID = &m.ID
				}
			}

			payload := datatypes.JSON(event.Payload)
			if len(payload) == 0 {
				payload = datatypes.JSON([]byte(`{}`))
			}
			rows = append(rows, model.OffchainEvent{
				MarketID:      marketID,
				SourceEventID: event.EventID,
				EventTime:     event.Time.UTC(),
				EventType:     truncateString(strings.TrimSpace(event.EventType), 30),
				Source:        strings.TrimSpace(event.Source),
				Title:         event.Title,
				Payload:       payload,
			})
		}

		if len(rows) > 0 {
			if err := s.offchainRepo.UpsertMany(ctx, rows); err != nil {
				return err
			}
		}
		if allOlderThanCutoff || len(events) < s.limit {
			break
		}
		offset += s.limit
	}

	if s.cursorRepo != nil && maxSeenTs > 0 && maxSeenTs >= cursorTs {
		if err := s.cursorRepo.Upsert(ctx, offchainCursorSource, offchainCursorStream, strconv.FormatInt(maxSeenTs, 10)); err != nil {
			return err
		}
	}
	return nil
}

func truncateString(in string, max int) string {
	if max <= 0 || len(in) <= max {
		return in
	}
	return in[:max]
}
