package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"easy-arbitra/backend/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"go.uber.org/zap"
)

type WalletAnalysisInput struct {
	WalletID      int64
	WalletAddress string
	StrategyType  string
	SmartScore    int
	InfoEdgeLevel string
	TradingPnL    float64
	MakerRebates  float64
	FeesPaid      float64
	TotalTrades   int64
	Volume30D     float64
	AsOf          time.Time
}

type WalletAnalysisOutput struct {
	ModelID      string
	ReportJSON   []byte
	NLSummary    string
	RiskWarnings []string
	InputTokens  int
	OutputTokens int
	LatencyMS    int
}

type Analyzer interface {
	AnalyzeWallet(ctx context.Context, in WalletAnalysisInput) (*WalletAnalysisOutput, error)
}

type BedrockAnalyzer struct {
	cfg      config.NovaConfig
	client   *bedrockruntime.Client
	fallback *MockAnalyzer
}

type NovaDevAPIAnalyzer struct {
	cfg      config.NovaConfig
	client   *http.Client
	fallback *MockAnalyzer
}

type MockAnalyzer struct {
	cfg    config.NovaConfig
	logger *zap.Logger
}

func NewAnalyzer(cfg config.NovaConfig, logger *zap.Logger) Analyzer {
	fallback := &MockAnalyzer{cfg: cfg, logger: logger}
	if !cfg.Enabled {
		return fallback
	}

	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "", "devapi", "nova-dev-api", "nova_api":
		if strings.TrimSpace(cfg.APIKey) == "" {
			logger.Warn("nova dev api selected but NOVA_API_KEY is empty, fallback to mock")
			return fallback
		}
		baseURL := strings.TrimSpace(cfg.APIBaseURL)
		if baseURL == "" {
			baseURL = "https://api.nova.amazon.com/v1"
		}
		return &NovaDevAPIAnalyzer{
			cfg: cfg,
			client: &http.Client{
				Timeout: 30 * time.Second,
			},
			fallback: fallback,
		}
	case "bedrock":
		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.Region))
		if err != nil {
			logger.Warn("nova bedrock selected but aws config init failed, fallback to mock", zap.Error(err))
			return fallback
		}
		return &BedrockAnalyzer{
			cfg:      cfg,
			client:   bedrockruntime.NewFromConfig(awsCfg),
			fallback: fallback,
		}
	default:
		logger.Warn("unknown NOVA_PROVIDER, fallback to mock", zap.String("provider", cfg.Provider))
		return fallback
	}
}

func (a *NovaDevAPIAnalyzer) AnalyzeWallet(ctx context.Context, in WalletAnalysisInput) (*WalletAnalysisOutput, error) {
	start := time.Now()
	fallback, _ := a.fallback.AnalyzeWallet(ctx, in)
	systemPrompt, userPrompt := buildPrompts(in)

	maxTokens := a.cfg.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 2048
	}
	baseURL := strings.TrimRight(strings.TrimSpace(a.cfg.APIBaseURL), "/")
	if baseURL == "" {
		baseURL = "https://api.nova.amazon.com/v1"
	}

	reqBody := map[string]any{
		"model": a.cfg.AnalysisModel,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"response_format": map[string]string{"type": "json_object"},
		"max_tokens":      maxTokens,
		"temperature":     a.cfg.Temperature,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal nova dev api request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("new nova dev api request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nova dev api request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read nova dev api response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("nova dev api http %d: %s", resp.StatusCode, string(raw))
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content any `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode nova dev api response: %w", err)
	}

	content := ""
	if len(out.Choices) > 0 {
		switch v := out.Choices[0].Message.Content.(type) {
		case string:
			content = v
		case []any:
			parts := make([]string, 0, len(v))
			for _, item := range v {
				if s, ok := item.(string); ok {
					parts = append(parts, s)
				}
			}
			content = strings.Join(parts, "\n")
		}
	}

	reportJSON, summary, warnings := mergeModelOutput(content, fallback)
	return &WalletAnalysisOutput{
		ModelID:      a.cfg.AnalysisModel,
		ReportJSON:   reportJSON,
		NLSummary:    summary,
		RiskWarnings: warnings,
		InputTokens:  out.Usage.PromptTokens,
		OutputTokens: out.Usage.CompletionTokens,
		LatencyMS:    int(time.Since(start).Milliseconds()),
	}, nil
}

func (a *BedrockAnalyzer) AnalyzeWallet(ctx context.Context, in WalletAnalysisInput) (*WalletAnalysisOutput, error) {
	start := time.Now()
	fallback, _ := a.fallback.AnalyzeWallet(ctx, in)
	systemPrompt, userPrompt := buildPrompts(in)

	maxTokens := int32(a.cfg.MaxTokens)
	if maxTokens <= 0 {
		maxTokens = 2048
	}

	resp, err := a.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId: &a.cfg.AnalysisModel,
		System: []bedrocktypes.SystemContentBlock{
			&bedrocktypes.SystemContentBlockMemberText{Value: systemPrompt},
		},
		Messages: []bedrocktypes.Message{
			{
				Role: bedrocktypes.ConversationRoleUser,
				Content: []bedrocktypes.ContentBlock{
					&bedrocktypes.ContentBlockMemberText{Value: userPrompt},
				},
			},
		},
		InferenceConfig: &bedrocktypes.InferenceConfiguration{
			MaxTokens:   &maxTokens,
			Temperature: &a.cfg.Temperature,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("bedrock converse: %w", err)
	}

	text := extractConverseText(resp.Output)
	reportJSON, summary, warnings := mergeModelOutput(text, fallback)

	inputTokens := 0
	outputTokens := 0
	if resp.Usage != nil {
		if resp.Usage.InputTokens != nil {
			inputTokens = int(*resp.Usage.InputTokens)
		}
		if resp.Usage.OutputTokens != nil {
			outputTokens = int(*resp.Usage.OutputTokens)
		}
	}

	return &WalletAnalysisOutput{
		ModelID:      a.cfg.AnalysisModel,
		ReportJSON:   reportJSON,
		NLSummary:    summary,
		RiskWarnings: warnings,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		LatencyMS:    int(time.Since(start).Milliseconds()),
	}, nil
}

func buildPrompts(in WalletAnalysisInput) (string, string) {
	systemPrompt := "You are an on-chain wallet analyst. Return ONLY valid JSON."
	userPrompt := fmt.Sprintf(`Analyze this wallet data and return strict JSON object with keys:
- natural_language_summary (string)
- risk_warnings (string array)
- layer2_strategy (object)
- layer3_info_edge (object)

Input:
{
  "wallet_id": %d,
  "wallet_address": %q,
  "as_of": %q,
  "layer1_facts": {
    "trading_pnl": %.8f,
    "maker_rebates": %.8f,
    "fees_paid": %.8f,
    "total_trades": %d,
    "volume_30d": %.8f
  },
  "layer2_seed": {
    "strategy_type": %q,
    "smart_score": %d
  },
  "layer3_seed": {
    "info_edge_level": %q
  }
}`, in.WalletID, in.WalletAddress, in.AsOf.UTC().Format(time.RFC3339), in.TradingPnL, in.MakerRebates, in.FeesPaid, in.TotalTrades, in.Volume30D, in.StrategyType, in.SmartScore, in.InfoEdgeLevel)
	return systemPrompt, userPrompt
}

func mergeModelOutput(modelText string, fallback *WalletAnalysisOutput) ([]byte, string, []string) {
	report := make(map[string]any)
	if fallback != nil && len(fallback.ReportJSON) > 0 {
		_ = json.Unmarshal(fallback.ReportJSON, &report)
	}
	if strings.TrimSpace(modelText) != "" {
		var modelReport map[string]any
		if err := json.Unmarshal([]byte(modelText), &modelReport); err == nil {
			for k, v := range modelReport {
				report[k] = v
			}
		} else {
			report["model_output_text"] = modelText
		}
	}

	payload, err := json.Marshal(report)
	if err != nil {
		payload = []byte("{}")
	}

	summary := ""
	warnings := []string{}
	if fallback != nil {
		summary = fallback.NLSummary
		warnings = fallback.RiskWarnings
	}
	if v, ok := report["natural_language_summary"].(string); ok && strings.TrimSpace(v) != "" {
		summary = v
	}
	if arr, ok := report["risk_warnings"].([]any); ok {
		out := make([]string, 0, len(arr))
		for _, item := range arr {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, s)
			}
		}
		if len(out) > 0 {
			warnings = out
		}
	}
	return payload, summary, warnings
}

func extractConverseText(output bedrocktypes.ConverseOutput) string {
	switch v := output.(type) {
	case *bedrocktypes.ConverseOutputMemberMessage:
		return extractMessageText(v.Value)
	default:
		return ""
	}
}

func extractMessageText(msg bedrocktypes.Message) string {
	parts := make([]string, 0, len(msg.Content))
	for _, block := range msg.Content {
		switch c := block.(type) {
		case *bedrocktypes.ContentBlockMemberText:
			parts = append(parts, c.Value)
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func (a *MockAnalyzer) AnalyzeWallet(_ context.Context, in WalletAnalysisInput) (*WalletAnalysisOutput, error) {
	report := map[string]any{
		"wallet_id":      in.WalletID,
		"wallet_address": in.WalletAddress,
		"as_of":          in.AsOf.UTC().Format(time.RFC3339),
		"layer1_facts": map[string]any{
			"trading_pnl":   in.TradingPnL,
			"maker_rebates": in.MakerRebates,
			"fees_paid":     in.FeesPaid,
			"total_trades":  in.TotalTrades,
			"volume_30d":    in.Volume30D,
		},
		"layer2_strategy": map[string]any{
			"primary_type": in.StrategyType,
			"confidence":   0.68,
			"evidence_points": []string{
				"strategy derived from wallet_features_daily",
				"smart score from latest wallet_score row",
			},
		},
		"layer3_info_edge": map[string]any{
			"level":                in.InfoEdgeLevel,
			"confidence":           0.63,
			"mean_delta_t_minutes": -15,
			"early_entry_rate":     0.21,
		},
		"smart_score":              in.SmartScore,
		"natural_language_summary": fmt.Sprintf("Wallet %s shows %s behavior with score %d. 30d volume %.2f across %d trades.", in.WalletAddress, in.StrategyType, in.SmartScore, in.Volume30D, in.TotalTrades),
	}
	payload, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("marshal report: %w", err)
	}

	summary := fmt.Sprintf("Wallet %s currently maps to %s strategy with smart score %d and %.2f 30d volume.", in.WalletAddress, in.StrategyType, in.SmartScore, in.Volume30D)
	return &WalletAnalysisOutput{
		ModelID:      a.cfg.AnalysisModel,
		ReportJSON:   payload,
		NLSummary:    summary,
		RiskWarnings: []string{"Probabilistic estimate only", "Not investment advice"},
		InputTokens:  400,
		OutputTokens: 240,
		LatencyMS:    120,
	}, nil
}
