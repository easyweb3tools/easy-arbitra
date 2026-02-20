package worker

import (
	"context"

	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/internal/service"
)

type ScoreCalculator struct {
	walletRepo *repository.WalletRepository
	classifier *service.ClassificationService
}

func NewScoreCalculator(walletRepo *repository.WalletRepository, classifier *service.ClassificationService) *ScoreCalculator {
	return &ScoreCalculator{walletRepo: walletRepo, classifier: classifier}
}

func (s *ScoreCalculator) Name() string { return "score_calculator" }

func (s *ScoreCalculator) RunOnce(ctx context.Context) error {
	ids, err := s.walletRepo.ListIDs(ctx)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.classifier.ClassifyWallet(ctx, id); err != nil {
			continue
		}
	}
	return nil
}
