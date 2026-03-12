package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/mark3labs/mcp-go/mcp"
)

type ReportPayload struct {
	WalletCard WalletCard `json:"wallet_card"`
	RadarChart RadarChart `json:"radar_chart"`
	Report     Report     `json:"report"`
}

type WalletCard struct {
	Address      string `json:"address"`
	DisplayName  string `json:"display_name"`
	ProfileImage string `json:"profile_image"`
	Sport        string `json:"sport"`
	TotalTrades  int    `json:"total_trades"`
}

type RadarChart struct {
	EntryTiming float64 `json:"entry_timing"`
	SizeRatio   float64 `json:"size_ratio"`
	Conviction  float64 `json:"conviction"`
}

type Report struct {
	StyleLabel     string `json:"style_label"`
	SummaryContext string `json:"summary_context"`
}

func BuildReportPayload() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		walletInfoJSON, ok := args["wallet_info"].(string)
		if !ok || walletInfoJSON == "" {
			return mcp.NewToolResultError("wallet_info parameter is required"), nil
		}

		metricsJSON, ok := args["metrics_json"].(string)
		if !ok || metricsJSON == "" {
			return mcp.NewToolResultError("metrics_json parameter is required"), nil
		}

		tradesSummaryJSON, _ := args["trades_summary"].(string)

		// Parse wallet info
		var walletInfo ResolveResult
		if err := json.Unmarshal([]byte(walletInfoJSON), &walletInfo); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse wallet_info: %v", err)), nil
		}

		// Parse metrics
		var metricsData MetricsResult
		if err := json.Unmarshal([]byte(metricsJSON), &metricsData); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse metrics_json: %v", err)), nil
		}

		// Normalize radar chart values (0-1)
		entryTiming := math.Max(0, 1-metricsData.Metrics.EntryTimingHours/24)
		sizeRatio := math.Min(1, metricsData.Metrics.SizeRatioPct/10)
		// Conviction is already 0-1 (average buy price)
		conviction := metricsData.Metrics.Conviction

		// Determine style label
		styleLabel := DetermineStyleLabel(entryTiming, sizeRatio, conviction)

		// Build summary context
		summaryContext := fmt.Sprintf(
			"Entry timing: %.1f hours avg | Position size: %.4f%% of market volume | Conviction: %.2f (avg buy price) | Sample: %d trades",
			metricsData.Metrics.EntryTimingHours,
			metricsData.Metrics.SizeRatioPct,
			metricsData.Metrics.Conviction,
			metricsData.SampleSize,
		)

		if tradesSummaryJSON != "" {
			summaryContext += " | Trades data available for detailed analysis"
		}

		payload := ReportPayload{
			WalletCard: WalletCard{
				Address:      walletInfo.WalletAddress,
				DisplayName:  walletInfo.DisplayName,
				ProfileImage: walletInfo.ProfileImage,
				Sport:        "NBA",
				TotalTrades:  metricsData.SampleSize,
			},
			RadarChart: RadarChart{
				EntryTiming: math.Round(entryTiming*100) / 100,
				SizeRatio:   math.Round(sizeRatio*100) / 100,
				Conviction:  math.Round(conviction*100) / 100,
			},
			Report: Report{
				StyleLabel:     styleLabel,
				SummaryContext: summaryContext,
			},
		}

		data, _ := json.Marshal(payload)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func DetermineStyleLabel(entryTiming, sizeRatio, conviction float64) string {
	if entryTiming > 0.7 && sizeRatio > 0.5 {
		return "Early Whale"
	}
	if entryTiming > 0.7 && sizeRatio <= 0.5 {
		return "Quick Scout"
	}
	if entryTiming <= 0.3 && sizeRatio > 0.5 {
		return "Late Whale"
	}
	if conviction > 0.75 {
		return "Favorite Backer"
	}
	if conviction < 0.35 && conviction > 0 {
		return "Contrarian Hunter"
	}
	if sizeRatio > 0.7 {
		return "Heavy Hitter"
	}
	if entryTiming > 0.5 {
		return "Early Bird"
	}
	return "Steady Player"
}
