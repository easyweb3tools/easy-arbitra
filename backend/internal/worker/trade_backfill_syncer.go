package worker

import (
	"context"
	"sync"

	"easy-arbitra/backend/internal/client"
	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/pkg/polyaddr"
)

type TradeBackfillSyncer struct {
	dataClient        *client.DataAPIClient
	walletRepo        *repository.WalletRepository
	marketRepo        *repository.MarketRepository
	tokenRepo         *repository.TokenRepository
	tradeRepo         *repository.TradeRepository
	walletsPerSync    int
	pagesPerWallet    int
	pageSize          int
	concurrency       int
	targetMinTrades   int64
	minCandidateTrade int64
}

func NewTradeBackfillSyncer(
	dataClient *client.DataAPIClient,
	walletRepo *repository.WalletRepository,
	marketRepo *repository.MarketRepository,
	tokenRepo *repository.TokenRepository,
	tradeRepo *repository.TradeRepository,
	walletsPerSync int,
	pagesPerWallet int,
	pageSize int,
	concurrency int,
	targetMinTrades int64,
) *TradeBackfillSyncer {
	if walletsPerSync <= 0 {
		walletsPerSync = 20
	}
	if pagesPerWallet <= 0 {
		pagesPerWallet = 3
	}
	if pageSize <= 0 {
		pageSize = 200
	}
	if concurrency <= 0 {
		concurrency = 8
	}
	if targetMinTrades < 2 {
		targetMinTrades = 100
	}
	return &TradeBackfillSyncer{
		dataClient:        dataClient,
		walletRepo:        walletRepo,
		marketRepo:        marketRepo,
		tokenRepo:         tokenRepo,
		tradeRepo:         tradeRepo,
		walletsPerSync:    walletsPerSync,
		pagesPerWallet:    pagesPerWallet,
		pageSize:          pageSize,
		concurrency:       concurrency,
		targetMinTrades:   targetMinTrades,
		minCandidateTrade: 1,
	}
}

func (s *TradeBackfillSyncer) Name() string { return "trade_backfill_syncer" }

func (s *TradeBackfillSyncer) RunOnce(ctx context.Context) error {
	candidates, err := s.walletRepo.ListBackfillCandidates(ctx, s.minCandidateTrade, s.targetMinTrades-1, s.walletsPerSync)
	if err != nil {
		return err
	}
	if len(candidates) == 0 {
		return nil
	}

	workerN := s.concurrency
	if workerN > len(candidates) {
		workerN = len(candidates)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan repository.WalletTradeCountRow)
	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	workerFn := func() {
		defer wg.Done()
		for candidate := range jobs {
			if err := s.backfillWallet(ctx, candidate); err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})
				return
			}
		}
	}
	for i := 0; i < workerN; i++ {
		wg.Add(1)
		go workerFn()
	}

	for _, candidate := range candidates {
		select {
		case <-ctx.Done():
			break
		case jobs <- candidate:
		}
	}
	close(jobs)
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}
	return nil
}

func (s *TradeBackfillSyncer) backfillWallet(ctx context.Context, candidate repository.WalletTradeCountRow) error {
	user := polyaddr.BytesToHex(candidate.Address)
	offset := 0
	for page := 0; page < s.pagesPerWallet; page++ {
		trades, err := s.dataClient.FetchTradesByUser(ctx, user, s.pageSize, offset)
		if err != nil {
			return err
		}
		if len(trades) == 0 {
			break
		}
		for _, trade := range trades {
			if err := ingestTrade(ctx, trade, s.walletRepo, s.marketRepo, s.tokenRepo, s.tradeRepo); err != nil {
				return err
			}
		}
		if len(trades) < s.pageSize {
			break
		}
		offset += s.pageSize
	}
	return nil
}
