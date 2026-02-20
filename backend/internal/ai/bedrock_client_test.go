package ai

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"easy-arbitra/backend/config"
	"go.uber.org/zap"
)

func TestMockAnalyzerIncludesLayerFields(t *testing.T) {
	logger := zap.NewNop()
	analyzer := NewAnalyzer(config.NovaConfig{AnalysisModel: "mock-model"}, logger)

	out, err := analyzer.AnalyzeWallet(context.Background(), WalletAnalysisInput{
		WalletID:      1,
		WalletAddress: "0xabc",
		StrategyType:  "quant",
		SmartScore:    77,
		InfoEdgeLevel: "processing_edge",
		TradingPnL:    12.3,
		MakerRebates:  1.2,
		FeesPaid:      0.5,
		TotalTrades:   10,
		Volume30D:     999,
		AsOf:          time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("AnalyzeWallet error: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(out.ReportJSON, &payload); err != nil {
		t.Fatalf("unmarshal report: %v", err)
	}
	if _, ok := payload["layer1_facts"]; !ok {
		t.Fatalf("expected layer1_facts in report")
	}
	if _, ok := payload["natural_language_summary"]; !ok {
		t.Fatalf("expected natural_language_summary in report")
	}
}
