package polymarket

import (
	"net/http"
	"time"
)

const (
	GammaBase = "https://gamma-api.polymarket.com"
	DataBase  = "https://data-api.polymarket.com"
)

// Client wraps HTTP calls to Polymarket APIs.
type Client struct {
	HTTP      *http.Client
	GammaBase string
	DataBase  string
}

// NewClient creates a Polymarket API client with a 10-second timeout.
func NewClient() *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
		GammaBase: GammaBase,
		DataBase:  DataBase,
	}
}
