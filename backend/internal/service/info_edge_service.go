package service

import (
	"context"
	"math"

	"easy-arbitra/backend/internal/repository"
)

type InfoEdgeService struct {
	tradeRepo *repository.TradeRepository
}

type InfoEdgeReport struct {
	WalletID         int64   `json:"wallet_id"`
	MeanDeltaMinutes float64 `json:"mean_delta_minutes"`
	StdDevMinutes    float64 `json:"stddev_minutes"`
	Samples          int64   `json:"samples"`
	ZScore           float64 `json:"z_score"`
	PValue           float64 `json:"p_value"`
	Classification   string  `json:"classification"`
}

func NewInfoEdgeService(tradeRepo *repository.TradeRepository) *InfoEdgeService {
	return &InfoEdgeService{tradeRepo: tradeRepo}
}

func (s *InfoEdgeService) Evaluate(ctx context.Context, walletID int64) (*InfoEdgeReport, error) {
	timing, err := s.tradeRepo.TimingSummaryByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	report := &InfoEdgeReport{
		WalletID:         walletID,
		MeanDeltaMinutes: timing.MeanDeltaMinutes,
		StdDevMinutes:    timing.StdDevMinutes,
		Samples:          timing.Samples,
		Classification:   "insufficient_data",
	}
	if timing.Samples < 5 || timing.StdDevMinutes <= 0 {
		return report, nil
	}

	stderr := timing.StdDevMinutes / math.Sqrt(float64(timing.Samples))
	if stderr <= 0 {
		return report, nil
	}

	report.ZScore = timing.MeanDeltaMinutes / stderr
	cdf := 0.5 * (1 + math.Erf(math.Abs(report.ZScore)/math.Sqrt2))
	report.PValue = 2 * (1 - cdf)

	switch {
	case report.PValue < 0.05 && report.MeanDeltaMinutes <= -30:
		report.Classification = "processing_edge"
	case report.PValue < 0.10 && report.MeanDeltaMinutes < 0:
		report.Classification = "mild_edge"
	default:
		report.Classification = "no_edge"
	}

	return report, nil
}
