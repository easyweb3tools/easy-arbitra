package polymarket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetTrades fetches trades for a user address with pagination.
func (c *Client) GetTrades(user string, limit, offset int) ([]Trade, error) {
	u := fmt.Sprintf("%s/trades?user=%s&limit=%d&offset=%d",
		c.DataBase, url.QueryEscape(user), limit, offset)
	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, fmt.Errorf("trades request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("trades API returned %d: %s", resp.StatusCode, string(body))
	}

	var trades []Trade
	if err := json.NewDecoder(resp.Body).Decode(&trades); err != nil {
		return nil, fmt.Errorf("decode trades: %w", err)
	}
	return trades, nil
}

// GetRecentTrades fetches recent global trades with pagination.
func (c *Client) GetRecentTrades(limit, offset int) ([]Trade, error) {
	u := fmt.Sprintf("%s/trades?limit=%d&offset=%d",
		c.DataBase, limit, offset)
	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, fmt.Errorf("recent trades request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("recent trades API returned %d: %s", resp.StatusCode, string(body))
	}

	var trades []Trade
	if err := json.NewDecoder(resp.Body).Decode(&trades); err != nil {
		return nil, fmt.Errorf("decode recent trades: %w", err)
	}
	return trades, nil
}
