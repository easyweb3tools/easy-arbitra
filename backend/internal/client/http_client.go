package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	baseURL      string
	client       *http.Client
	maxRetries   int
	retryBackoff time.Duration
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	return &HTTPClient{
		baseURL:      strings.TrimRight(baseURL, "/"),
		client:       &http.Client{Timeout: timeout},
		maxRetries:   2,
		retryBackoff: 250 * time.Millisecond,
	}
}

func (c *HTTPClient) GetJSON(ctx context.Context, path string, out any) error {
	url := c.baseURL + path
	attempts := c.maxRetries + 1
	var lastErr error

	for i := 0; i < attempts; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("new request: %w", err)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("do request: %w", err)
		} else {
			bodyErr := decodeJSONResponse(resp, out)
			if bodyErr == nil {
				return nil
			}
			lastErr = bodyErr
			if resp.StatusCode < 500 {
				return lastErr
			}
		}

		if i < attempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.retryBackoff * time.Duration(i+1)):
			}
		}
	}

	return lastErr
}

func decodeJSONResponse(resp *http.Response, out any) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}
