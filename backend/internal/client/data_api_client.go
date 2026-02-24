package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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

type tradeWireDTO struct {
	TransactionHash string  `json:"transactionHash"`
	Timestamp       int64   `json:"timestamp"`
	Market          string  `json:"market"`
	ConditionID     string  `json:"conditionId"`
	TokenID         string  `json:"asset"`
	MakerAddress    string  `json:"makerAddress"`
	TakerAddress    string  `json:"takerAddress"`
	ProxyWallet     string  `json:"proxyWallet"`
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
	return c.FetchTradesPage(ctx, limit, 0)
}

func (c *DataAPIClient) FetchTradesPage(ctx context.Context, limit int, offset int) ([]TradeDTO, error) {
	if limit <= 0 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	var raw json.RawMessage
	if err := c.http.GetJSON(ctx, fmt.Sprintf("/trades?limit=%d&offset=%d", limit, offset), &raw); err != nil {
		return nil, err
	}
	rows, err := decodeTrades(raw)
	if err != nil {
		return nil, fmt.Errorf("decode trades: %w", err)
	}
	return rows, nil
}

func (c *DataAPIClient) FetchTradesByUser(ctx context.Context, user string, limit int, offset int) ([]TradeDTO, error) {
	if limit <= 0 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	qUser := url.QueryEscape(user)
	var raw json.RawMessage
	if err := c.http.GetJSON(ctx, fmt.Sprintf("/trades?user=%s&limit=%d&offset=%d", qUser, limit, offset), &raw); err != nil {
		return nil, err
	}
	rows, err := decodeTrades(raw)
	if err != nil {
		return nil, fmt.Errorf("decode user trades: %w", err)
	}
	return rows, nil
}

func decodeTrades(raw json.RawMessage) ([]TradeDTO, error) {
	var arr []tradeWireDTO
	if err := json.Unmarshal(raw, &arr); err == nil {
		return normalizeTrades(arr), nil
	}

	var wrapped struct {
		Data []tradeWireDTO `json:"data"`
	}
	if err := json.Unmarshal(raw, &wrapped); err == nil && wrapped.Data != nil {
		return normalizeTrades(wrapped.Data), nil
	}

	var wrappedTrades struct {
		Trades []tradeWireDTO `json:"trades"`
	}
	if err := json.Unmarshal(raw, &wrappedTrades); err == nil && wrappedTrades.Trades != nil {
		return normalizeTrades(wrappedTrades.Trades), nil
	}

	return nil, fmt.Errorf("unsupported data-api response format")
}

func normalizeTrades(rows []tradeWireDTO) []TradeDTO {
	out := make([]TradeDTO, 0, len(rows))
	for _, row := range rows {
		market := row.Market
		if row.ConditionID != "" {
			market = row.ConditionID
		}
		taker := row.TakerAddress
		if taker == "" {
			taker = row.ProxyWallet
		}
		out = append(out, TradeDTO{
			TransactionHash: row.TransactionHash,
			Timestamp:       row.Timestamp,
			Market:          market,
			TokenID:         row.TokenID,
			MakerAddress:    row.MakerAddress,
			TakerAddress:    taker,
			Price:           row.Price,
			Size:            row.Size,
			Side:            row.Side,
			FeePaid:         row.FeePaid,
		})
	}
	return out
}
