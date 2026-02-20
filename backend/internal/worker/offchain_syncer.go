package worker

import (
	"context"

	"easy-arbitra/backend/internal/client"
)

type OffchainEventSyncer struct {
	client *client.OffchainClient
	limit  int
}

func NewOffchainEventSyncer(client *client.OffchainClient, limit int) *OffchainEventSyncer {
	return &OffchainEventSyncer{client: client, limit: limit}
}

func (s *OffchainEventSyncer) Name() string { return "offchain_event_syncer" }

func (s *OffchainEventSyncer) RunOnce(ctx context.Context) error {
	_, err := s.client.FetchEvents(ctx, s.limit)
	return err
}
