package worker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"easy-arbitra/backend/internal/client"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
)

func ingestTrade(
	ctx context.Context,
	trade client.TradeDTO,
	walletRepo *repository.WalletRepository,
	marketRepo *repository.MarketRepository,
	tokenRepo *repository.TokenRepository,
	tradeRepo *repository.TradeRepository,
) error {
	if strings.TrimSpace(trade.TransactionHash) == "" || strings.TrimSpace(trade.TokenID) == "" || strings.TrimSpace(trade.Market) == "" {
		return nil
	}

	market, err := marketRepo.EnsureByConditionID(ctx, trade.Market)
	if err != nil {
		return err
	}
	token, err := tokenRepo.EnsureToken(ctx, market.ID, trade.TokenID, 1)
	if err != nil {
		return err
	}

	var makerID *int64
	if strings.HasPrefix(strings.ToLower(trade.MakerAddress), "0x") {
		maker, err := walletRepo.EnsureByAddress(ctx, trade.MakerAddress)
		if err == nil {
			makerID = &maker.ID
		}
	}
	var takerID *int64
	if strings.HasPrefix(strings.ToLower(trade.TakerAddress), "0x") {
		taker, err := walletRepo.EnsureByAddress(ctx, trade.TakerAddress)
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
	uniqSeed := fmt.Sprintf(
		"%s|%d|%s|%s|%s|%s|%.10f|%.10f|%d",
		strings.ToLower(trade.TransactionHash),
		trade.Timestamp,
		strings.ToLower(strings.TrimSpace(trade.TokenID)),
		strings.ToLower(strings.TrimSpace(trade.MakerAddress)),
		strings.ToLower(strings.TrimSpace(trade.TakerAddress)),
		strings.ToLower(strings.TrimSpace(trade.Side)),
		trade.Price,
		trade.Size,
		side,
	)
	sum := sha256.Sum256([]byte(uniqSeed))
	uniqKey := hex.EncodeToString(sum[:])
	return tradeRepo.Upsert(ctx, model.TradeFill{
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
}
