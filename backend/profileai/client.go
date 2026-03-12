package profileai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var allowedLabels = []string{
	"Quick Scout",
	"Heavy Hitter",
	"Favorite Backer",
	"Contrarian Hunter",
	"Early Whale",
	"Late Whale",
	"Early Bird",
	"Steady Player",
}

type Client struct {
	baseURL    string
	model      string
	apiKey     string
	timeout    time.Duration
	httpClient *http.Client
}

type Input struct {
	Wallet                  string
	DisplayName             string
	SourceRank              int
	WinRate                 float64
	PnlUSD                  float64
	NbaTrades               int
	RecentMarkets           int
	EntryTimingHours        float64
	SizeRatioPct            float64
	Conviction              float64
	DeterministicStyleLabel string
	PresentationScore       float64
}

type Result struct {
	StyleLabel string
	Summary    string
	Source     string
	Model      string
}

func NewFromEnv() *Client {
	baseURL := strings.TrimSpace(getenv("AI_BASE_URL", ""))
	model := strings.TrimSpace(getenv("AI_MODEL", ""))
	apiKey := strings.TrimSpace(getenv("AI_API_KEY", ""))
	timeout := parseTimeout(getenv("AI_TIMEOUT_MS", ""), 45*time.Second)

	return &Client{
		baseURL: baseURL,
		model:   model,
		apiKey:  apiKey,
		timeout: timeout,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Configured() bool {
	return c != nil && c.baseURL != "" && c.model != "" && c.apiKey != ""
}

func (c *Client) Classify(ctx context.Context, input Input) (Result, error) {
	if !c.Configured() {
		return Result{
			StyleLabel: input.DeterministicStyleLabel,
			Summary:    "AI tagging is not configured; using deterministic style label.",
			Source:     "fallback",
			Model:      "",
		}, nil
	}

	body, err := json.Marshal(map[string]any{
		"model":       c.model,
		"temperature": 0.2,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": buildSystemPrompt(),
			},
			{
				"role":    "user",
				"content": buildUserPrompt(input),
			},
		},
	})
	if err != nil {
		return Result{}, fmt.Errorf("marshal ai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chatCompletionsURL(c.baseURL), bytes.NewReader(body))
	if err != nil {
		return Result{}, fmt.Errorf("build ai request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("call ai provider: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var raw bytes.Buffer
		raw.ReadFrom(resp.Body)
		return Result{}, fmt.Errorf("ai provider returned %d: %s", resp.StatusCode, raw.String())
	}

	var payload struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return Result{}, fmt.Errorf("decode ai response: %w", err)
	}
	if len(payload.Choices) == 0 {
		return Result{}, fmt.Errorf("ai response missing choices")
	}

	var parsed struct {
		StyleLabel string `json:"style_label"`
		Summary    string `json:"style_summary"`
	}
	content := extractJSONObject(payload.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return Result{}, fmt.Errorf("parse ai profile json: %w", err)
	}
	if parsed.StyleLabel == "" {
		parsed.StyleLabel = input.DeterministicStyleLabel
	}
	if !isAllowedLabel(parsed.StyleLabel) {
		parsed.StyleLabel = input.DeterministicStyleLabel
	}
	if parsed.Summary == "" {
		parsed.Summary = "AI returned an empty summary; using deterministic style label."
	}

	return Result{
		StyleLabel: parsed.StyleLabel,
		Summary:    parsed.Summary,
		Source:     "ai",
		Model:      c.model,
	}, nil
}

func buildSystemPrompt() string {
	return fmt.Sprintf(`You are classifying public Polymarket NBA trader behavior.
Return JSON only with keys "style_label" and "style_summary".
Choose style_label from this exact set: %s.
style_summary must be one sentence, under 28 words, grounded only in the metrics provided.`, strings.Join(allowedLabels, ", "))
}

func buildUserPrompt(input Input) string {
	data, _ := json.Marshal(input)
	return string(data)
}

func chatCompletionsURL(baseURL string) string {
	normalized := strings.TrimRight(baseURL, "/")
	if strings.HasSuffix(normalized, "/chat/completions") {
		return normalized
	}
	return normalized + "/chat/completions"
}

func extractJSONObject(content string) string {
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start >= 0 && end > start {
		return content[start : end+1]
	}
	return content
}

func isAllowedLabel(label string) bool {
	for _, allowed := range allowedLabels {
		if label == allowed {
			return true
		}
	}
	return false
}

func parseTimeout(value string, fallback time.Duration) time.Duration {
	if value == "" {
		return fallback
	}
	ms, err := time.ParseDuration(value + "ms")
	if err == nil && ms > 0 {
		return ms
	}
	return fallback
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(strings.Trim(os.Getenv(key), `"`)); value != "" {
		return value
	}
	return fallback
}
