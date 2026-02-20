package client

import (
	"context"
	"time"
)

type OffchainEvent struct {
	Title     string
	EventType string
	Source    string
	Time      time.Time
}

type OffchainClient struct{}

func NewOffchainClient() *OffchainClient {
	return &OffchainClient{}
}

func (c *OffchainClient) FetchEvents(ctx context.Context, limit int) ([]OffchainEvent, error) {
	_ = ctx
	if limit <= 0 {
		limit = 20
	}
	_ = limit
	return []OffchainEvent{}, nil
}
