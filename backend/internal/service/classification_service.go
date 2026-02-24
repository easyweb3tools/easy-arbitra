package service

import (
	"context"
	"fmt"
	"time"

	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"gorm.io/datatypes"
)

type ClassificationService struct {
	featureRepo *repository.FeatureRepository
	scoreRepo   *repository.ScoreRepository
}

func NewClassificationService(featureRepo *repository.FeatureRepository, scoreRepo *repository.ScoreRepository) *ClassificationService {
	return &ClassificationService{featureRepo: featureRepo, scoreRepo: scoreRepo}
}

func (s *ClassificationService) ClassifyWallet(ctx context.Context, walletID int64) error {
	f, err := s.featureRepo.LatestByWalletID(ctx, walletID)
	if err != nil {
		return err
	}

	strategy := "quant"
	infoEdge := "quant"
	confidence := 0.6
	score := 60

	if f.TradeCount30d < 30 {
		strategy = "lucky"
		infoEdge = "luck"
		confidence = 0.45
		score = 40
	} else if f.MakerRatio > 0.7 {
		strategy = "market_maker"
		infoEdge = "processing_edge"
		confidence = 0.72
		score = 70
	} else if f.Pnl30d < 0 {
		strategy = "noise"
		infoEdge = "luck"
		confidence = 0.58
		score = 35
	} else if f.ActiveDays30d < 5 {
		strategy = "event_trader"
		infoEdge = "processing_edge"
		confidence = 0.62
		score = 56
	}
	if f.Pnl30d > 1000 {
		score += 10
	}
	if f.AvgEdge > 0.02 {
		score += 8
		infoEdge = "processing_edge"
	}
	if f.TxFrequencyPerDay > 15 {
		score += 4
	}
	if score > 95 {
		score = 95
	}
	if score < 0 {
		score = 0
	}

	poolTier := "observation"
	if score >= 80 && f.Pnl30d > 0 && f.Pnl7d >= 0 {
		poolTier = "star"
	} else if score >= 65 && f.Pnl30d > 0 {
		poolTier = "strategy"
	}
	riskLevel := "high"
	if score >= 80 {
		riskLevel = "low"
	} else if score >= 60 {
		riskLevel = "medium"
	}
	suitableFor := "aggressive"
	if strategy == "market_maker" || strategy == "quant" {
		suitableFor = "conservative"
	} else if strategy == "event_trader" || strategy == "arbitrage" {
		suitableFor = "event_driven"
	}
	suggestedPosition := "1-3%"
	if riskLevel == "low" {
		suggestedPosition = "5-10%"
	} else if riskLevel == "medium" {
		suggestedPosition = "3-5%"
	}
	momentum := "stable"
	if f.Pnl7d > 0 {
		momentum = "heating"
	} else if f.Pnl7d < 0 {
		momentum = "cooling"
	}
	now := time.Now().UTC()

	detail := datatypes.JSON([]byte(fmt.Sprintf(
		`{"pnl_7d":%.2f,"pnl_30d":%.2f,"pnl_90d":%.2f,"maker_ratio":%.4f,"trade_count_30d":%d,"active_days_30d":%d,"tx_frequency_per_day":%.3f,"avg_edge":%.6f}`,
		f.Pnl7d, f.Pnl30d, f.Pnl90d, f.MakerRatio, f.TradeCount30d, f.ActiveDays30d, f.TxFrequencyPerDay, f.AvgEdge,
	)))
	return s.scoreRepo.UpsertLatest(ctx, model.WalletScore{
		WalletID:           walletID,
		ScoredAt:           now.Truncate(time.Hour),
		StrategyType:       strategy,
		StrategyConfidence: confidence,
		InfoEdgeLevel:      infoEdge,
		InfoEdgeConfidence: confidence,
		SmartScore:         score,
		PoolTier:           poolTier,
		PoolTierUpdatedAt:  &now,
		SuitableFor:        suitableFor,
		RiskLevel:          riskLevel,
		SuggestedPosition:  suggestedPosition,
		Momentum:           momentum,
		ScoringDetail:      detail,
	})
}
