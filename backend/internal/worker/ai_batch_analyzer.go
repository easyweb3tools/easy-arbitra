package worker

import (
	"context"
	"errors"
	"strings"
	"time"

	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/internal/service"
)

type AIBatchAnalyzer struct {
	aiService      *service.AIService
	walletRepo     *repository.WalletRepository
	batchSize      int
	minTrades      int64
	minRealizedPnL float64
	cooldown       time.Duration
	requestSpacing time.Duration
}

func NewAIBatchAnalyzer(
	aiService *service.AIService,
	walletRepo *repository.WalletRepository,
	batchSize int,
	minTrades int64,
	minRealizedPnL float64,
	cooldown time.Duration,
	requestSpacing time.Duration,
) *AIBatchAnalyzer {
	if batchSize <= 0 {
		batchSize = 3
	}
	if minTrades < 2 {
		minTrades = 100
	}
	if requestSpacing < 0 {
		requestSpacing = 0
	}
	return &AIBatchAnalyzer{
		aiService:      aiService,
		walletRepo:     walletRepo,
		batchSize:      batchSize,
		minTrades:      minTrades,
		minRealizedPnL: minRealizedPnL,
		cooldown:       cooldown,
		requestSpacing: requestSpacing,
	}
}

func (w *AIBatchAnalyzer) Name() string { return "ai_batch_analyzer" }

func (w *AIBatchAnalyzer) RunOnce(ctx context.Context) error {
	candidates, err := w.walletRepo.ListAIAnalyzeCandidates(ctx, w.minTrades, w.minRealizedPnL, w.cooldown, w.batchSize)
	if err != nil {
		return err
	}
	var firstErr error
	for i, candidate := range candidates {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		_, err := w.aiService.AnalyzeByWalletID(ctx, candidate.WalletID, false)
		if err == nil {
			if i < len(candidates)-1 {
				if err := waitWithContext(ctx, w.requestSpacing); err != nil {
					return err
				}
			}
			continue
		}
		if errors.Is(err, service.ErrInsufficientTrades) || errors.Is(err, service.ErrNonPositivePnL) || errors.Is(err, service.ErrNotFound) {
			continue
		}
		if isRateLimitError(err) {
			return nil
		}
		if firstErr == nil {
			firstErr = err
		}
		if i < len(candidates)-1 {
			if err := waitWithContext(ctx, w.requestSpacing); err != nil {
				return err
			}
		}
	}
	return firstErr
}

func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "http 429") || strings.Contains(msg, "rate limit")
}

func waitWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
