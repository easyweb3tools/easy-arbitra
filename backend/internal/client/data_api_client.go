package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type TradeDTO struct {
	TransactionHash string  `json:"transactionHash"`
	Timestamp       int64   `json:"timestamp"`
	Market          string  `json:"market"`
	TokenID         string  `json:"asset"`
	MakerAddress    string  `json:"makerAddress"`
	TakerAddress    string  `json:"takerAddress"`
	Price           float64 `json:"price"`
	Size            float64 `json:"size"`
	Side            string  `json:"side"`
	FeePaid         float64 `json:"fee"`
}

type DataAPIClient struct {
	http *HTTPClient
}

func NewDataAPIClient(baseURL string, timeout time.Duration) *DataAPIClient {
	return &DataAPIClient{http: NewHTTPClient(baseURL, timeout)}
}

func (c *DataAPIClient) FetchTrades(ctx context.Context, limit int) ([]TradeDTO, error) {
	if limit <= 0 {
		limit = 200
	}
	var raw json.RawMessage
	if err := c.http.GetJSON(ctx, fmt.Sprintf("/trades?limit=%d", limit), &raw); err != nil {
		return nil, err
	}
	rows, err := decodeTrades(raw)
	if err != nil {
		return nil, fmt.Errorf("decode trades: %w", err)
	}
	return rows, nil
}

func decodeTrades(raw json.RawMessage) ([]TradeDTO, error) {
	var arr []TradeDTO
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}

	var wrapped struct {
		Data []TradeDTO `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Data != nil {
		return wrapped.Data, nil
	}

	var wrappedTrades struct {
		Trades []TradeDTO `json:"trades"`
	}
	if err := json.Unmarshal(raw, &wrappedTrades); err == nil && wrappedTrades.Trades != nil {
		return wrappedTrades.Trades, nil
	}

	return nil, fmt.Errorf("unsupported data-api response format")
}
