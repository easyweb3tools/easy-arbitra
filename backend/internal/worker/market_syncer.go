package worker

import (
	"context"

	"easy-arbitra/backend/internal/client"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
)

type MarketSyncer struct {
	client *client.GammaClient
	repo   *repository.MarketRepository
	limit  int
}

func NewMarketSyncer(client *client.GammaClient, repo *repository.MarketRepository, limit int) *MarketSyncer {
	return &MarketSyncer{client: client, repo: repo, limit: limit}
}

func (s *MarketSyncer) Name() string { return "market_syncer" }

func (s *MarketSyncer) RunOnce(ctx context.Context) error {
	rows, err := s.client.FetchMarkets(ctx, s.limit)
	if err != nil {
		return err
	}
	markets := make([]model.Market, 0, len(rows))
	for _, row := range rows {
		status := int16(1)
		if row.Active {
			status = 0
		}
		markets = append(markets, model.Market{
			ConditionID: row.ConditionID,
			Slug:        row.Slug,
			Title:       row.Question,
			Category:    row.Category,
			Status:      status,
			Volume:      row.Volume,
			Liquidity:   row.Liquidity,
		})
	}
	return s.repo.UpsertMany(ctx, markets)
}
