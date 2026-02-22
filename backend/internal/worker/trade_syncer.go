package worker

import (
	"context"

	"easy-arbitra/backend/internal/client"
	"easy-arbitra/backend/internal/repository"
)

type TradeSyncer struct {
	dataClient *client.DataAPIClient
	walletRepo *repository.WalletRepository
	marketRepo *repository.MarketRepository
	tokenRepo  *repository.TokenRepository
	tradeRepo  *repository.TradeRepository
	limit      int
}

func NewTradeSyncer(
	dataClient *client.DataAPIClient,
	walletRepo *repository.WalletRepository,
	marketRepo *repository.MarketRepository,
	tokenRepo *repository.TokenRepository,
	tradeRepo *repository.TradeRepository,
	limit int,
) *TradeSyncer {
	return &TradeSyncer{
		dataClient: dataClient,
		walletRepo: walletRepo,
		marketRepo: marketRepo,
		tokenRepo:  tokenRepo,
		tradeRepo:  tradeRepo,
		limit:      limit,
	}
}

func (s *TradeSyncer) Name() string { return "trade_syncer" }

func (s *TradeSyncer) RunOnce(ctx context.Context) error {
	trades, err := s.dataClient.FetchTrades(ctx, s.limit)
	if err != nil {
		return err
	}
	for _, trade := range trades {
		if err := ingestTrade(ctx, trade, s.walletRepo, s.marketRepo, s.tokenRepo, s.tradeRepo); err != nil {
			return err
		}
	}
	return nil
}
