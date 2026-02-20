package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"easy-arbitra/backend/config"
	"go.uber.org/zap"
)

type WalletAnalysisInput struct {
	WalletID      int64
	WalletAddress string
	StrategyType  string
	SmartScore    int
	InfoEdgeLevel string
	TradingPnL    float64
	MakerRebates  float64
	FeesPaid      float64
	TotalTrades   int64
	Volume30D     float64
	AsOf          time.Time
}

type WalletAnalysisOutput struct {
	ModelID      string
	ReportJSON   []byte
	NLSummary    string
	RiskWarnings []string
	InputTokens  int
	OutputTokens int
	LatencyMS    int
}

type Analyzer interface {
	AnalyzeWallet(ctx context.Context, in WalletAnalysisInput) (*WalletAnalysisOutput, error)
}

type MockAnalyzer struct {
	cfg    config.NovaConfig
	logger *zap.Logger
}

func NewAnalyzer(cfg config.NovaConfig, logger *zap.Logger) Analyzer {
	return &MockAnalyzer{cfg: cfg, logger: logger}
}

func (a *MockAnalyzer) AnalyzeWallet(_ context.Context, in WalletAnalysisInput) (*WalletAnalysisOutput, error) {
	report := map[string]any{
		"wallet_id":      in.WalletID,
		"wallet_address": in.WalletAddress,
		"as_of":          in.AsOf.UTC().Format(time.RFC3339),
		"layer1_facts": map[string]any{
			"trading_pnl":   in.TradingPnL,
			"maker_rebates": in.MakerRebates,
			"fees_paid":     in.FeesPaid,
			"total_trades":  in.TotalTrades,
			"volume_30d":    in.Volume30D,
		},
		"layer2_strategy": map[string]any{
			"primary_type": in.StrategyType,
			"confidence":   0.68,
			"evidence_points": []string{
				"strategy derived from wallet_features_daily",
				"smart score from latest wallet_score row",
			},
		},
		"layer3_info_edge": map[string]any{
			"level":                in.InfoEdgeLevel,
			"confidence":           0.63,
			"mean_delta_t_minutes": -15,
			"early_entry_rate":     0.21,
		},
		"smart_score":              in.SmartScore,
		"natural_language_summary": fmt.Sprintf("Wallet %s shows %s behavior with score %d. 30d volume %.2f across %d trades.", in.WalletAddress, in.StrategyType, in.SmartScore, in.Volume30D, in.TotalTrades),
	}
	payload, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("marshal report: %w", err)
	}

	summary := fmt.Sprintf("Wallet %s currently maps to %s strategy with smart score %d and %.2f 30d volume.", in.WalletAddress, in.StrategyType, in.SmartScore, in.Volume30D)
	return &WalletAnalysisOutput{
		ModelID:      a.cfg.AnalysisModel,
		ReportJSON:   payload,
		NLSummary:    summary,
		RiskWarnings: []string{"Probabilistic estimate only", "Not investment advice"},
		InputTokens:  400,
		OutputTokens: 240,
		LatencyMS:    120,
	}, nil
}
