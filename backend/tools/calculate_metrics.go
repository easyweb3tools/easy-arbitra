package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/brucexwang/easy-arbitra/backend/metrics"
	"github.com/brucexwang/easy-arbitra/backend/polymarket"
	"github.com/mark3labs/mcp-go/mcp"
)

type MetricsResult struct {
	Wallet     string         `json:"wallet"`
	Metrics    StyleMetrics   `json:"metrics"`
	SampleSize int            `json:"sample_size"`
	Warning    string         `json:"warning,omitempty"`
}

type StyleMetrics struct {
	EntryTimingHours float64 `json:"entry_timing_hours"`
	SizeRatioPct     float64 `json:"size_ratio_pct"`
	Conviction       float64 `json:"conviction"`
}

func CalculateStyleMetrics() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		wallet, ok := args["wallet"].(string)
		if !ok || wallet == "" {
			return mcp.NewToolResultError("wallet parameter is required"), nil
		}

		tradesJSON, ok := args["trades_json"].(string)
		if !ok || tradesJSON == "" {
			return mcp.NewToolResultError("trades_json parameter is required"), nil
		}

		// Accept both the full FetchTradesResult object and a plain []EnrichedTrade array
		var trades []polymarket.EnrichedTrade
		var wrapped FetchTradesResult
		if err := json.Unmarshal([]byte(tradesJSON), &wrapped); err == nil && wrapped.Wallet != "" {
			trades = wrapped.Trades
		} else if err := json.Unmarshal([]byte(tradesJSON), &trades); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to parse trades_json: %v", err)), nil
		}

		if len(trades) == 0 {
			result := MetricsResult{
				Wallet:     wallet,
				Metrics:    StyleMetrics{},
				SampleSize: 0,
			}
			data, _ := json.Marshal(result)
			return mcp.NewToolResultText(string(data)), nil
		}

		result := MetricsResult{
			Wallet: wallet,
			Metrics: StyleMetrics{
				EntryTimingHours: metrics.EntryTimingHours(trades),
				SizeRatioPct:     metrics.SizeRatioPct(trades),
				Conviction:       metrics.Conviction(trades),
			},
			SampleSize: len(trades),
		}

		if len(trades) < 3 {
			result.Warning = fmt.Sprintf("Small sample size (%d trades). Metrics may not be representative.", len(trades))
		}

		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}
