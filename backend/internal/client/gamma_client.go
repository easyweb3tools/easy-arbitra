package client

import (
	"context"
	"encoding/json"
	"fmt"
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
	var arr []GammaMarket
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}

	var wrapped struct {
		Data []GammaMarket `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Data != nil {
		return wrapped.Data, nil
	}

	var wrappedMarkets struct {
		Markets []GammaMarket `json:"markets"`
	}
	if err := json.Unmarshal(raw, &wrappedMarkets); err == nil && wrappedMarkets.Markets != nil {
		return wrappedMarkets.Markets, nil
	}
	return nil, fmt.Errorf("unsupported gamma response format")
}
