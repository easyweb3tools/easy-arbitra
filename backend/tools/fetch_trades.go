package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/brucexwang/easy-arbitra/backend/polymarket"
	"github.com/mark3labs/mcp-go/mcp"
)

type FetchTradesResult struct {
	Wallet      string                     `json:"wallet"`
	Sport       string                     `json:"sport"`
	TotalTrades int                        `json:"total_trades"`
	Trades      []polymarket.EnrichedTrade `json:"trades"`
}

func FetchSportsTrades(client *polymarket.Client) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()
		wallet, ok := args["wallet"].(string)
		if !ok || wallet == "" {
			return mcp.NewToolResultError("wallet parameter is required"), nil
		}

		// Parse optional sport param (default: nba)
		sport := "nba"
		if s, ok := args["sport"].(string); ok && s != "" {
			sport = strings.ToLower(s)
		}

		// Parse optional limit param (default: 500)
		tradeLimit := 500
		if l, ok := args["limit"].(float64); ok && l > 0 {
			tradeLimit = int(l)
		}

		result, err := FetchSportsTradesData(ctx, client, wallet, sport, tradeLimit)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}

func FetchSportsTradesData(
	ctx context.Context,
	client *polymarket.Client,
	wallet, sport string,
	tradeLimit int,
) (FetchTradesResult, error) {
	LogToolf(ctx, "Fetching recent trades for %s (limit=%d)", wallet, tradeLimit)
	trades, err := client.GetTrades(wallet, tradeLimit, 0)
	if err != nil {
		return FetchTradesResult{}, fmt.Errorf("failed to get trades: %v", err)
	}
	LogToolf(ctx, "Fetched %d raw trades", len(trades))

	if len(trades) == 0 {
		return FetchTradesResult{
			Wallet:      wallet,
			Sport:       sport,
			TotalTrades: 0,
			Trades:      []polymarket.EnrichedTrade{},
		}, nil
	}

	conditionIDSet := make(map[string]bool)
	for _, trade := range trades {
		if trade.ConditionID != "" {
			conditionIDSet[trade.ConditionID] = true
		}
	}

	conditionIDList := make([]string, 0, len(conditionIDSet))
	for id := range conditionIDSet {
		conditionIDList = append(conditionIDList, id)
	}
	LogToolf(ctx, "Resolving metadata for %d unique markets", len(conditionIDList))

	marketMap := make(map[string]polymarket.Market)
	sportConditionIDs := make(map[string]bool)
	for i := 0; i < len(conditionIDList); i += 20 {
		end := i + 20
		if end > len(conditionIDList) {
			end = len(conditionIDList)
		}
		batch := conditionIDList[i:end]
		LogToolf(ctx, "Fetching market metadata batch %d-%d", i+1, end)

		markets, err := client.GetMarkets(batch)
		if err != nil {
			LogToolf(ctx, "Skipping market batch %d-%d after error: %v", i+1, end, err)
			continue
		}

		for _, m := range markets {
			marketMap[m.ConditionID] = m
			if isSportMarket(m, sport) {
				sportConditionIDs[m.ConditionID] = true
			}
		}
	}

	LogToolf(ctx, "Filtering %d trades down to %s positions", len(trades), strings.ToUpper(sport))
	var sportTrades []polymarket.Trade
	for _, trade := range trades {
		if sportConditionIDs[trade.ConditionID] || isSportText(trade.Title, sport) {
			sportTrades = append(sportTrades, trade)
		}
	}

	LogToolf(ctx, "Building enriched response for %d %s trades", len(sportTrades), strings.ToUpper(sport))
	enriched := make([]polymarket.EnrichedTrade, 0, len(sportTrades))
	for _, t := range sportTrades {
		et := polymarket.EnrichedTrade{
			ConditionID:     t.ConditionID,
			MarketQuestion:  "",
			TradeTime:       t.Time().Format(time.RFC3339),
			Side:            t.Side,
			Size:            t.Size,
			Price:           t.Price,
			Outcome:         t.Outcome,
			MarketVolume:    0,
			MarketStartTime: "",
		}
		if m, ok := marketMap[t.ConditionID]; ok {
			et.MarketQuestion = m.Question
			et.MarketVolume = m.VolumeNum
			et.MarketStartTime = m.StartDate
		}
		enriched = append(enriched, et)
	}
	LogToolf(ctx, "Fetch trades complete")

	return FetchTradesResult{
		Wallet:      wallet,
		Sport:       sport,
		TotalTrades: len(enriched),
		Trades:      enriched,
	}, nil
}

func isSportMarket(m polymarket.Market, sport string) bool {
	return isSportText(m.Question, sport) || isSportText(m.Slug, sport)
}

func isSportText(text, sport string) bool {
	value := strings.ToLower(text)
	if strings.Contains(value, sport) {
		return true
	}

	return sport == "nba" && strings.Contains(value, "basketball")
}
