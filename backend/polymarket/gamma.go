package polymarket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetPublicProfile fetches the public profile for a wallet address.
func (c *Client) GetPublicProfile(address string) (*Profile, error) {
	u := fmt.Sprintf("%s/public-profile?address=%s", c.GammaBase, url.QueryEscape(address))
	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, fmt.Errorf("profile request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("profile API returned %d: %s", resp.StatusCode, string(body))
	}

	var profile Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}
	return &profile, nil
}

// GetSportsTags fetches all sport tags from the Gamma API.
func (c *Client) GetSportsTags() ([]Tag, error) {
	u := fmt.Sprintf("%s/sports", c.GammaBase)
	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, fmt.Errorf("sports request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sports API returned %d: %s", resp.StatusCode, string(body))
	}

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, fmt.Errorf("decode sports: %w", err)
	}
	return tags, nil
}

// GetEvents fetches events for a given tag ID with pagination.
// Does not filter by active/closed status so historical events are included.
func (c *Client) GetEvents(tagID string, limit, offset int) ([]Event, error) {
	u := fmt.Sprintf("%s/events?tag=%s&limit=%d&offset=%d",
		c.GammaBase, url.QueryEscape(tagID), limit, offset)
	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, fmt.Errorf("events request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("events API returned %d: %s", resp.StatusCode, string(body))
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("decode events: %w", err)
	}
	return events, nil
}

// GetMarkets fetches markets by condition IDs (comma-separated).
func (c *Client) GetMarkets(conditionIDs []string) ([]Market, error) {
	if len(conditionIDs) == 0 {
		return nil, nil
	}
	csv := strings.Join(conditionIDs, ",")
	u := fmt.Sprintf("%s/markets?condition_ids=%s", c.GammaBase, url.QueryEscape(csv))
	resp, err := c.HTTP.Get(u)
	if err != nil {
		return nil, fmt.Errorf("markets request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("markets API returned %d: %s", resp.StatusCode, string(body))
	}

	var markets []Market
	if err := json.NewDecoder(resp.Body).Decode(&markets); err != nil {
		return nil, fmt.Errorf("decode markets: %w", err)
	}
	return markets, nil
}
