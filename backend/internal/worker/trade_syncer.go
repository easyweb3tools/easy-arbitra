package worker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"easy-arbitra/backend/internal/client"
	"easy-arbitra/backend/internal/model"
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
		if strings.TrimSpace(trade.TransactionHash) == "" || strings.TrimSpace(trade.TokenID) == "" || strings.TrimSpace(trade.Market) == "" {
			continue
		}

		market, err := s.marketRepo.EnsureByConditionID(ctx, trade.Market)
		if err != nil {
			return err
		}
		token, err := s.tokenRepo.EnsureToken(ctx, market.ID, trade.TokenID, 1)
		if err != nil {
			return err
		}

		var makerID *int64
		if strings.HasPrefix(strings.ToLower(trade.MakerAddress), "0x") {
			maker, err := s.walletRepo.EnsureByAddress(ctx, trade.MakerAddress)
			if err == nil {
				makerID = &maker.ID
			}
		}
		var takerID *int64
		if strings.HasPrefix(strings.ToLower(trade.TakerAddress), "0x") {
			taker, err := s.walletRepo.EnsureByAddress(ctx, trade.TakerAddress)
			if err == nil {
				takerID = &taker.ID
			}
		}

		side := int16(1)
		if strings.EqualFold(trade.Side, "sell") {
			side = 0
		}
		blockTime := time.Now().UTC()
		if trade.Timestamp > 0 {
			blockTime = time.Unix(trade.Timestamp, 0).UTC()
		}
		uniqKey := fmt.Sprintf("%s-%d", strings.ToLower(trade.TransactionHash), trade.Timestamp)
		err = s.tradeRepo.Upsert(ctx, model.TradeFill{
			TokenID:       token.ID,
			MakerWalletID: makerID,
			TakerWalletID: takerID,
			Side:          side,
			Price:         trade.Price,
			Size:          trade.Size,
			FeePaid:       trade.FeePaid,
			BlockTime:     blockTime,
			Source:        1,
			UniqKey:       uniqKey,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
