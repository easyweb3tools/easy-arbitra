package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type OffchainEvent struct {
	EventID     string
	Title       string
	EventType   string
	Source      string
	Time        time.Time
	ConditionID string
	Payload     json.RawMessage
}

type offchainEventWire struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Category    string `json:"category"`
	StartDate   string `json:"startDate"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	PublishedAt string `json:"published_at"`
	Markets     []struct {
		ConditionID string `json:"conditionId"`
	} `json:"markets"`
}

type OffchainClient struct {
	http *HTTPClient
}

func NewOffchainClient(baseURL string, timeout time.Duration) *OffchainClient {
	return &OffchainClient{http: NewHTTPClient(baseURL, timeout)}
}

func (c *OffchainClient) FetchEvents(ctx context.Context, limit int, offset int) ([]OffchainEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	var raw json.RawMessage
	if err := c.http.GetJSON(ctx, fmt.Sprintf("/events?limit=%d&offset=%d", limit, offset), &raw); err != nil {
		return nil, err
	}
	rows, err := decodeOffchainEvents(raw)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func decodeOffchainEvents(raw json.RawMessage) ([]OffchainEvent, error) {
	var arr []offchainEventWire
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil, fmt.Errorf("decode events: %w", err)
	}
	out := make([]OffchainEvent, 0, len(arr))
	for _, row := range arr {
		eventTime := parseEventTime(row.UpdatedAt)
		if eventTime.IsZero() {
			eventTime = parseEventTime(row.PublishedAt)
		}
		if eventTime.IsZero() {
			eventTime = parseEventTime(row.CreatedAt)
		}
		if eventTime.IsZero() {
			eventTime = parseEventTime(row.StartDate)
		}
		if eventTime.IsZero() {
			eventTime = time.Now().UTC()
		}

		conditionID := ""
		if len(row.Markets) > 0 {
			conditionID = row.Markets[0].ConditionID
		}
		payload, _ := json.Marshal(row)
		eventType := "gamma_event"
		if row.Category != "" {
			eventType = "gamma_" + row.Category
		}
		out = append(out, OffchainEvent{
			EventID:     row.ID,
			Title:       row.Title,
			EventType:   eventType,
			Source:      "gamma_api",
			Time:        eventTime,
			ConditionID: conditionID,
			Payload:     payload,
		})
	}
	return out, nil
}

func parseEventTime(raw string) time.Time {
	if raw == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999-07",
		"2006-01-02 15:04:05.999999Z07:00",
		"2006-01-02 15:04:05-07",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if ts, err := time.Parse(layout, raw); err == nil {
			return ts.UTC()
		}
	}
	return time.Time{}
}
