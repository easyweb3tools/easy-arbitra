package worker

import (
	"context"
	"time"

	"easy-arbitra/backend/internal/copytrade"
	"easy-arbitra/backend/internal/repository"
)

type CopyTradeSyncer struct {
	copyTradeService *copytrade.Service
	copyTradeRepo    *copytrade.Repository
	tradeRepo        *repository.TradeRepository
	marketRepo       *repository.MarketRepository
}

func NewCopyTradeSyncer(
	copyTradeService *copytrade.Service,
	copyTradeRepo *copytrade.Repository,
	tradeRepo *repository.TradeRepository,
	marketRepo *repository.MarketRepository,
) *CopyTradeSyncer {
	return &CopyTradeSyncer{
		copyTradeService: copyTradeService,
		copyTradeRepo:    copyTradeRepo,
		tradeRepo:        tradeRepo,
		marketRepo:       marketRepo,
	}
}

func (s *CopyTradeSyncer) Name() string { return "copy_trade_syncer" }

func (s *CopyTradeSyncer) RunOnce(ctx context.Context) error {
	configs, err := s.copyTradeRepo.ListEnabledConfigs(ctx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	for _, cfg := range configs {
		since := cfg.LastCheckedAt
		if since.IsZero() {
			since = now.Add(-24 * time.Hour)
		}

		trades, _, err := s.tradeRepo.ListByWalletID(ctx, cfg.WalletID, 20, 0)
		if err != nil {
			continue
		}

		for _, trade := range trades {
			if !trade.BlockTime.After(since) {
				continue
			}

			market, _ := s.marketRepo.GetBySlug(ctx, trade.MarketSlug)

			_, _ = s.copyTradeService.ProcessNewTrade(ctx, cfg.ID, trade, market)
		}

		_ = s.copyTradeRepo.UpdateLastChecked(ctx, cfg.ID, now)
	}

	return nil
}
