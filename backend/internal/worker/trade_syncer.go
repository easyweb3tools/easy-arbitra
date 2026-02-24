package worker

import (
	"context"
	"strconv"
	"time"

	"easy-arbitra/backend/internal/client"
	"easy-arbitra/backend/internal/repository"
	"gorm.io/gorm"
)

const (
	tradeCursorSource = "data_api"
	tradeCursorStream = "trades_latest_ts"
)

type TradeSyncer struct {
	dataClient     *client.DataAPIClient
	walletRepo     *repository.WalletRepository
	marketRepo     *repository.MarketRepository
	tokenRepo      *repository.TokenRepository
	tradeRepo      *repository.TradeRepository
	cursorRepo     *repository.IngestCursorRepository
	limit          int
	maxPages       int
	cursorLookback time.Duration
}

func NewTradeSyncer(
	dataClient *client.DataAPIClient,
	walletRepo *repository.WalletRepository,
	marketRepo *repository.MarketRepository,
	tokenRepo *repository.TokenRepository,
	tradeRepo *repository.TradeRepository,
	cursorRepo *repository.IngestCursorRepository,
	limit int,
	maxPages int,
	cursorLookback time.Duration,
) *TradeSyncer {
	if limit <= 0 {
		limit = 200
	}
	if maxPages <= 0 {
		maxPages = 20
	}
	if cursorLookback < 0 {
		cursorLookback = 0
	}
	return &TradeSyncer{
		dataClient:     dataClient,
		walletRepo:     walletRepo,
		marketRepo:     marketRepo,
		tokenRepo:      tokenRepo,
		tradeRepo:      tradeRepo,
		cursorRepo:     cursorRepo,
		limit:          limit,
		maxPages:       maxPages,
		cursorLookback: cursorLookback,
	}
}

func (s *TradeSyncer) Name() string { return "trade_syncer" }

func (s *TradeSyncer) RunOnce(ctx context.Context) error {
	cursorTs := int64(0)
	if s.cursorRepo != nil {
		cursor, err := s.cursorRepo.Get(ctx, tradeCursorSource, tradeCursorStream)
		if err == nil {
			if parsed, parseErr := strconv.ParseInt(cursor.CursorValue, 10, 64); parseErr == nil {
				cursorTs = parsed
			}
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}

	lookbackSeconds := int64(s.cursorLookback / time.Second)
	cutoffTs := cursorTs - lookbackSeconds
	if cutoffTs < 0 {
		cutoffTs = 0
	}

	maxSeenTs := cursorTs
	offset := 0
	for page := 0; page < s.maxPages; page++ {
		trades, err := s.dataClient.FetchTradesPage(ctx, s.limit, offset)
		if err != nil {
			return err
		}
		if len(trades) == 0 {
			break
		}

		allOlderThanCutoff := true
		for _, trade := range trades {
			ts := trade.Timestamp
			if ts > maxSeenTs {
				maxSeenTs = ts
			}
			if ts >= cutoffTs {
				allOlderThanCutoff = false
				if err := ingestTrade(ctx, trade, s.walletRepo, s.marketRepo, s.tokenRepo, s.tradeRepo); err != nil {
					return err
				}
			}
		}
		if allOlderThanCutoff {
			break
		}
		if len(trades) < s.limit {
			break
		}
		offset += s.limit
	}

	if s.cursorRepo != nil && maxSeenTs > 0 && maxSeenTs >= cursorTs {
		if err := s.cursorRepo.Upsert(ctx, tradeCursorSource, tradeCursorStream, strconv.FormatInt(maxSeenTs, 10)); err != nil {
			return err
		}
	}
	return nil
}
