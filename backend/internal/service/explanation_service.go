package service

import (
	"context"
	"errors"
	"time"

	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/pkg/polyaddr"
	"gorm.io/gorm"
)

type ExplanationService struct {
	walletRepo  *repository.WalletRepository
	featureRepo *repository.FeatureRepository
	scoreRepo   *repository.ScoreRepository
	tradeRepo   *repository.TradeRepository
	infoEdge    *InfoEdgeService
}

type WalletExplanation struct {
	WalletID    int64          `json:"wallet_id"`
	Address     string         `json:"address"`
	Layer1      map[string]any `json:"layer1"`
	Layer2      map[string]any `json:"layer2"`
	Layer3      map[string]any `json:"layer3"`
	Disclosures []string       `json:"disclosures"`
	GeneratedAt string         `json:"generated_at"`
}

func NewExplanationService(
	walletRepo *repository.WalletRepository,
	featureRepo *repository.FeatureRepository,
	scoreRepo *repository.ScoreRepository,
	tradeRepo *repository.TradeRepository,
	infoEdge *InfoEdgeService,
) *ExplanationService {
	return &ExplanationService{walletRepo: walletRepo, featureRepo: featureRepo, scoreRepo: scoreRepo, tradeRepo: tradeRepo, infoEdge: infoEdge}
}

func (s *ExplanationService) GetWalletExplanation(ctx context.Context, walletID int64) (*WalletExplanation, error) {
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	pnl, err := s.tradeRepo.AggregateByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	info, _ := s.infoEdge.Evaluate(ctx, walletID)
	feature, _ := s.featureRepo.LatestByWalletID(ctx, walletID)
	score, _ := s.scoreRepo.LatestByWalletID(ctx, walletID)

	layer2 := map[string]any{"strategy_type": "unknown", "smart_score": 0, "confidence": 0.0}
	if score != nil {
		layer2 = map[string]any{
			"strategy_type":   score.StrategyType,
			"smart_score":     score.SmartScore,
			"confidence":      score.StrategyConfidence,
			"info_edge_level": score.InfoEdgeLevel,
		}
	}

	layer3 := map[string]any{"mean_delta_minutes": 0.0, "samples": 0, "label": "insufficient_data"}
	if info != nil {
		layer3 = map[string]any{
			"mean_delta_minutes": info.MeanDeltaMinutes,
			"stddev_minutes":     info.StdDevMinutes,
			"samples":            info.Samples,
			"p_value":            info.PValue,
			"label":              info.Classification,
		}
	}

	layer1 := map[string]any{
		"trading_pnl":   pnl.TradingPnL,
		"maker_rebates": pnl.MakerRebates,
		"fees_paid":     pnl.FeesPaid,
		"total_trades":  pnl.TotalTrades,
		"volume_30d":    pnl.Volume30D,
	}
	if feature != nil {
		layer1["pnl_7d"] = feature.Pnl7d
		layer1["pnl_30d"] = feature.Pnl30d
		layer1["pnl_90d"] = feature.Pnl90d
		layer1["avg_edge"] = feature.AvgEdge
	}

	return &WalletExplanation{
		WalletID: walletID,
		Address:  polyaddr.BytesToHex(wallet.Address),
		Layer1:   layer1,
		Layer2:   layer2,
		Layer3:   layer3,
		Disclosures: []string{
			"Classification and score are probabilistic outputs.",
			"Timing analysis depends on offchain event completeness.",
		},
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
