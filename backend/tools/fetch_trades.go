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

		// Parse optional limit param (default: 3000 scanned trades)
		tradeLimit := 3000
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
	const pageSize = 500
	const targetSportTrades = 40

	if tradeLimit < pageSize {
		tradeLimit = pageSize
	}

	LogToolf(ctx, "Scanning up to %d recent trades for %s signals", tradeLimit, strings.ToUpper(sport))

	allTrades := make([]polymarket.Trade, 0, minInt(tradeLimit, pageSize))
	sportTrades := make([]polymarket.Trade, 0, 64)
	scanned := 0

	for offset := 0; offset < tradeLimit; offset += pageSize {
		pageLimit := minInt(pageSize, tradeLimit-offset)
		LogToolf(ctx, "Fetching trades page offset=%d limit=%d", offset, pageLimit)

		pageTrades, err := client.GetTrades(wallet, pageLimit, offset)
		if err != nil {
			return FetchTradesResult{}, fmt.Errorf("failed to get trades: %v", err)
		}

		allTrades = append(allTrades, pageTrades...)
		scanned += len(pageTrades)
		for _, trade := range pageTrades {
			if isSportTrade(trade, sport) {
				sportTrades = append(sportTrades, trade)
			}
		}

		LogToolf(ctx, "Fetched %d trades on this page (%d %s matches so far)", len(pageTrades), len(sportTrades), strings.ToUpper(sport))

		if len(pageTrades) < pageLimit {
			break
		}
		if len(sportTrades) >= targetSportTrades {
			break
		}
	}

	LogToolf(ctx, "Scanned %d raw trades and found %d %s candidates", scanned, len(sportTrades), strings.ToUpper(sport))

	if len(allTrades) == 0 {
		return FetchTradesResult{
			Wallet:      wallet,
			Sport:       sport,
			TotalTrades: 0,
			Trades:      []polymarket.EnrichedTrade{},
		}, nil
	}

	conditionIDSet := make(map[string]bool)
	for _, trade := range sportTrades {
		if trade.ConditionID != "" {
			conditionIDSet[trade.ConditionID] = true
		}
	}

	conditionIDList := make([]string, 0, len(conditionIDSet))
	for id := range conditionIDSet {
		conditionIDList = append(conditionIDList, id)
	}
	LogToolf(ctx, "Resolving metadata for %d %s markets", len(conditionIDList), strings.ToUpper(sport))

	marketMap := make(map[string]polymarket.Market)
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

func isSportTrade(t polymarket.Trade, sport string) bool {
	return isSportText(t.Title, sport) || isSportText(t.Slug, sport)
}

func isSportText(text, sport string) bool {
	value := strings.ToLower(text)
	if strings.Contains(value, sport) {
		return true
	}

	return sport == "nba" && strings.Contains(value, "basketball")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
