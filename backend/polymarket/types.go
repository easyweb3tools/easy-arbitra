package polymarket

import (
	"encoding/json"
	"fmt"
	"time"
)

// Profile represents a Polymarket user profile from Gamma API.
type Profile struct {
	ProxyWallet  string `json:"proxyWallet"`
	Name         string `json:"name"`
	Pseudonym    string `json:"pseudonym"`
	ProfileImage string `json:"profileImage"`
	Bio          string `json:"bio"`
}

// Trade represents a single trade from Data API.
type Trade struct {
	ID          string  `json:"id"`
	ProxyWallet string  `json:"proxyWallet"`
	Side        string  `json:"side"`        // "BUY" or "SELL"
	Asset       string  `json:"asset"`       // token ID
	ConditionID string  `json:"conditionId"` // market condition ID
	Size        float64 `json:"size,string"`
	Price       float64 `json:"price,string"`
	Timestamp   int64   `json:"timestamp"`
	Title       string  `json:"title"`
	Outcome     string  `json:"outcome"` // "Yes" or "No"
}

func (t Trade) Time() time.Time {
	return time.Unix(t.Timestamp, 0)
}

// Market represents a Polymarket market from Gamma API.
type Market struct {
	ID          string  `json:"id"`
	Question    string  `json:"question"`
	ConditionID string  `json:"condition_id"`
	Slug        string  `json:"slug"`
	VolumeNum   float64 `json:"volumeNum"`
	StartDate   string  `json:"startDateIso"`
	EndDate     string  `json:"endDateIso"`
	Active      bool    `json:"active"`
	Closed      bool    `json:"closed"`
}

// Event represents a Polymarket event from Gamma API.
type Event struct {
	ID      string   `json:"id"`
	Slug    string   `json:"slug"`
	Title   string   `json:"title"`
	Markets []Market `json:"markets"`
}

// Tag represents a sport/category tag from Gamma API.
type Tag struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Slug  string `json:"slug"`
}

func (t *Tag) UnmarshalJSON(data []byte) error {
	type rawTag struct {
		ID    json.RawMessage `json:"id"`
		Label string          `json:"label"`
		Slug  string          `json:"slug"`
	}

	var raw rawTag
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.Label = raw.Label
	t.Slug = raw.Slug

	if len(raw.ID) == 0 || string(raw.ID) == "null" {
		return nil
	}

	if err := json.Unmarshal(raw.ID, &t.ID); err == nil {
		return nil
	}

	var numericID json.Number
	if err := json.Unmarshal(raw.ID, &numericID); err == nil {
		t.ID = numericID.String()
		return nil
	}

	return fmt.Errorf("unsupported tag id: %s", string(raw.ID))
}

// EnrichedTrade is a trade enriched with market metadata.
type EnrichedTrade struct {
	ConditionID     string  `json:"condition_id"`
	MarketQuestion  string  `json:"market_question"`
	TradeTime       string  `json:"trade_time"`
	Side            string  `json:"side"`
	Size            float64 `json:"size"`
	Price           float64 `json:"price"`
	Outcome         string  `json:"outcome"`
	MarketVolume    float64 `json:"market_volume"`
	MarketStartTime string  `json:"market_start_time"`
}
