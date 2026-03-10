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
	Wallet      string                      `json:"wallet"`
	Sport       string                      `json:"sport"`
	TotalTrades int                         `json:"total_trades"`
	Trades      []polymarket.EnrichedTrade   `json:"trades"`
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

		// Step 1: Get sport tag ID from sports metadata
		tags, err := client.GetSportsTags()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get sports tags: %v", err)), nil
		}

		var sportTagID string
		for _, tag := range tags {
			slug := strings.ToLower(tag.Slug)
			label := strings.ToLower(tag.Label)
			if strings.Contains(slug, sport) || strings.Contains(label, sport) {
				sportTagID = tag.ID
				break
			}
		}

		if sportTagID == "" {
			sportTagID = sport
		}

		// Step 2: Get sport events (paginate to cover historical data)
		sportConditionIDs := make(map[string]bool)
		for offset := 0; ; offset += 100 {
			events, err := client.GetEvents(sportTagID, 100, offset)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get %s events: %v", sport, err)), nil
			}
			for _, event := range events {
				for _, market := range event.Markets {
					if market.ConditionID != "" {
						sportConditionIDs[market.ConditionID] = true
					}
				}
			}
			if len(events) < 100 {
				break // last page
			}
			time.Sleep(200 * time.Millisecond)
		}

		// Step 3: Get user trades
		trades, err := client.GetTrades(wallet, tradeLimit, 0)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get trades: %v", err)), nil
		}

		// Step 4: Filter to sport trades
		// First pass: match against known sport conditionIDs
		var sportTrades []polymarket.Trade
		var unmatchedConditionIDs []string
		unmatchedSet := make(map[string]bool)
		for _, trade := range trades {
			if sportConditionIDs[trade.ConditionID] {
				sportTrades = append(sportTrades, trade)
			} else if !unmatchedSet[trade.ConditionID] {
				unmatchedSet[trade.ConditionID] = true
				unmatchedConditionIDs = append(unmatchedConditionIDs, trade.ConditionID)
			}
		}

		// Second pass: check unmatched trades against market metadata
		// to catch sport markets not returned by the events endpoint
		for i := 0; i < len(unmatchedConditionIDs); i += 20 {
			end := i + 20
			if end > len(unmatchedConditionIDs) {
				end = len(unmatchedConditionIDs)
			}
			markets, err := client.GetMarkets(unmatchedConditionIDs[i:end])
			if err != nil {
				continue
			}
			for _, m := range markets {
				q := strings.ToLower(m.Question)
				if strings.Contains(q, sport) || (sport == "nba" && strings.Contains(q, "basketball")) {
					sportConditionIDs[m.ConditionID] = true
				}
			}
			if end < len(unmatchedConditionIDs) {
				time.Sleep(200 * time.Millisecond)
			}
		}

		// Re-filter with expanded sport set
		sportTrades = nil
		for _, trade := range trades {
			if sportConditionIDs[trade.ConditionID] {
				sportTrades = append(sportTrades, trade)
			}
		}

		// Step 5: Get market metadata for enrichment
		conditionIDSet := make(map[string]bool)
		for _, t := range sportTrades {
			conditionIDSet[t.ConditionID] = true
		}
		conditionIDList := make([]string, 0, len(conditionIDSet))
		for id := range conditionIDSet {
			conditionIDList = append(conditionIDList, id)
		}

		marketMap := make(map[string]polymarket.Market)
		// Batch in groups of 20 to avoid URL length limits
		for i := 0; i < len(conditionIDList); i += 20 {
			end := i + 20
			if end > len(conditionIDList) {
				end = len(conditionIDList)
			}
			batch := conditionIDList[i:end]
			markets, err := client.GetMarkets(batch)
			if err != nil {
				continue // skip failed batches
			}
			for _, m := range markets {
				marketMap[m.ConditionID] = m
			}
			// Rate limit
			if end < len(conditionIDList) {
				time.Sleep(200 * time.Millisecond)
			}
		}

		// Step 6: Build enriched trades
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

		result := FetchTradesResult{
			Wallet:      wallet,
			Sport:       sport,
			TotalTrades: len(enriched),
			Trades:      enriched,
		}

		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}
