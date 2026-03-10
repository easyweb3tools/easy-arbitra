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

		// Step 1: Get NBA tag ID from sports metadata
		tags, err := client.GetSportsTags()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get sports tags: %v", err)), nil
		}

		var nbaTagID string
		for _, tag := range tags {
			slug := strings.ToLower(tag.Slug)
			label := strings.ToLower(tag.Label)
			if strings.Contains(slug, "nba") || strings.Contains(label, "nba") {
				nbaTagID = tag.ID
				break
			}
		}

		if nbaTagID == "" {
			// Fallback: try with known slugs
			nbaTagID = "nba"
		}

		// Step 2: Get NBA events (paginate to cover historical data)
		nbaConditionIDs := make(map[string]bool)
		for offset := 0; ; offset += 100 {
			events, err := client.GetEvents(nbaTagID, 100, offset)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to get NBA events: %v", err)), nil
			}
			for _, event := range events {
				for _, market := range event.Markets {
					if market.ConditionID != "" {
						nbaConditionIDs[market.ConditionID] = true
					}
				}
			}
			if len(events) < 100 {
				break // last page
			}
			time.Sleep(200 * time.Millisecond)
		}

		// Step 3: Get user trades
		trades, err := client.GetTrades(wallet, 500, 0)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get trades: %v", err)), nil
		}

		// Step 4: Filter to NBA trades
		// First pass: match against known NBA conditionIDs
		var nbaTrades []polymarket.Trade
		var unmatchedConditionIDs []string
		unmatchedSet := make(map[string]bool)
		for _, trade := range trades {
			if nbaConditionIDs[trade.ConditionID] {
				nbaTrades = append(nbaTrades, trade)
			} else if !unmatchedSet[trade.ConditionID] {
				unmatchedSet[trade.ConditionID] = true
				unmatchedConditionIDs = append(unmatchedConditionIDs, trade.ConditionID)
			}
		}

		// Second pass: check unmatched trades against market metadata
		// to catch NBA markets not returned by the events endpoint
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
				if strings.Contains(q, "nba") || strings.Contains(q, "basketball") {
					nbaConditionIDs[m.ConditionID] = true
				}
			}
			if end < len(unmatchedConditionIDs) {
				time.Sleep(200 * time.Millisecond)
			}
		}

		// Re-filter with expanded NBA set
		nbaTrades = nil
		for _, trade := range trades {
			if nbaConditionIDs[trade.ConditionID] {
				nbaTrades = append(nbaTrades, trade)
			}
		}

		// Step 5: Get market metadata for enrichment
		conditionIDSet := make(map[string]bool)
		for _, t := range nbaTrades {
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
		enriched := make([]polymarket.EnrichedTrade, 0, len(nbaTrades))
		for _, t := range nbaTrades {
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
			Sport:       "nba",
			TotalTrades: len(enriched),
			Trades:      enriched,
		}

		data, _ := json.Marshal(result)
		return mcp.NewToolResultText(string(data)), nil
	}
}
