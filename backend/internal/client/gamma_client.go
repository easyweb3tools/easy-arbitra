package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type GammaMarket struct {
	ConditionID string  `json:"conditionId"`
	Slug        string  `json:"slug"`
	Question    string  `json:"question"`
	Category    string  `json:"category"`
	Active      bool    `json:"active"`
	Volume      float64 `json:"volume"`
	Liquidity   float64 `json:"liquidity"`
}

type gammaMarketWire struct {
	ConditionID string          `json:"conditionId"`
	Slug        string          `json:"slug"`
	Question    string          `json:"question"`
	Category    string          `json:"category"`
	Active      bool            `json:"active"`
	Volume      json.RawMessage `json:"volume"`
	Liquidity   json.RawMessage `json:"liquidity"`
}

type GammaClient struct {
	http *HTTPClient
}

func NewGammaClient(baseURL string, timeout time.Duration) *GammaClient {
	return &GammaClient{http: NewHTTPClient(baseURL, timeout)}
}

func (c *GammaClient) FetchMarkets(ctx context.Context, limit int) ([]GammaMarket, error) {
	if limit <= 0 {
		limit = 100
	}
	var raw json.RawMessage
	if err := c.http.GetJSON(ctx, fmt.Sprintf("/markets?limit=%d", limit), &raw); err != nil {
		return nil, err
	}

	rows, err := decodeGammaMarkets(raw)
	if err != nil {
		return nil, fmt.Errorf("decode gamma markets: %w", err)
	}
	return rows, nil
}

func decodeGammaMarkets(raw json.RawMessage) ([]GammaMarket, error) {
	var arr []gammaMarketWire
	if err := json.Unmarshal(raw, &arr); err == nil {
		return normalizeGammaMarkets(arr), nil
	}

	var wrapped struct {
		Data []gammaMarketWire `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Data != nil {
		return normalizeGammaMarkets(wrapped.Data), nil
	}

	var wrappedMarkets struct {
		Markets []gammaMarketWire `json:"markets"`
	}
	if err := json.Unmarshal(raw, &wrappedMarkets); err == nil && wrappedMarkets.Markets != nil {
		return normalizeGammaMarkets(wrappedMarkets.Markets), nil
	}
	return nil, fmt.Errorf("unsupported gamma response format")
}

func normalizeGammaMarkets(rows []gammaMarketWire) []GammaMarket {
	out := make([]GammaMarket, 0, len(rows))
	for _, row := range rows {
		out = append(out, GammaMarket{
			ConditionID: row.ConditionID,
			Slug:        row.Slug,
			Question:    row.Question,
			Category:    row.Category,
			Active:      row.Active,
			Volume:      parseJSONFloat(row.Volume),
			Liquidity:   parseJSONFloat(row.Liquidity),
		})
	}
	return out
}

func parseJSONFloat(raw json.RawMessage) float64 {
	if len(raw) == 0 {
		return 0
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return f
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return v
		}
	}
	return 0
}
