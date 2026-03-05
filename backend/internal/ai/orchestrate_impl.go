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

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

// ── Orchestrate implementations ──

func (a *NovaDevAPIAnalyzer) Orchestrate(ctx context.Context, in OrchestrateInput) (*OrchestrateOutput, error) {
	start := time.Now()
	systemPrompt, userPrompt := BuildOrchestratePrompts(in)

	maxTokens := a.cfg.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
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
		return a.fallback.Orchestrate(ctx, in)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return a.fallback.Orchestrate(ctx, in)
	}
	req.Header.Set("Authorization", "Bearer "+a.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return a.fallback.Orchestrate(ctx, in)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return a.fallback.Orchestrate(ctx, in)
	}
	if resp.StatusCode >= 400 {
		return a.fallback.Orchestrate(ctx, in)
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
		return a.fallback.Orchestrate(ctx, in)
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

	latency := int(time.Since(start).Milliseconds())
	return ParseOrchestrateResponse(content, a.cfg.AnalysisModel, out.Usage.PromptTokens, out.Usage.CompletionTokens, latency)
}

func (a *BedrockAnalyzer) Orchestrate(ctx context.Context, in OrchestrateInput) (*OrchestrateOutput, error) {
	start := time.Now()
	systemPrompt, userPrompt := BuildOrchestratePrompts(in)

	maxTokens := int32(a.cfg.MaxTokens)
	if maxTokens <= 0 {
		maxTokens = 4096
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
		return a.fallback.Orchestrate(ctx, in)
	}

	text := extractConverseText(resp.Output)
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

	latency := int(time.Since(start).Milliseconds())
	return ParseOrchestrateResponse(text, a.cfg.AnalysisModel, inputTokens, outputTokens, latency)
}

func (m *MockAnalyzer) Orchestrate(_ context.Context, in OrchestrateInput) (*OrchestrateOutput, error) {
	out := &OrchestrateOutput{
		Phase:       "analyzing",
		ModelID:     "mock",
		NLSummary:   fmt.Sprintf("Mock analysis round %d/%d with %d candidates", in.Round, in.TotalRounds, len(in.Candidates)),
		NLSummaryZh: fmt.Sprintf("模拟分析第 %d/%d 轮，共 %d 个候选", in.Round, in.TotalRounds, len(in.Candidates)),
	}

	for i, c := range in.Candidates {
		out.Rankings = append(out.Rankings, CandidateRank{
			WalletID: c.WalletID,
			Rank:     i + 1,
			Score:    float64(c.SmartScore),
			Reason:   fmt.Sprintf("Mock: %s strategy, score %d", c.StrategyType, c.SmartScore),
		})
	}

	out.Observations = fmt.Sprintf("Round %d: evaluated %d candidates.", in.Round, len(in.Candidates))

	// Force final on last round
	if in.IsLastRound && len(in.Candidates) > 0 {
		best := in.Candidates[0]
		out.Phase = "final"
		out.FinalPick = &FinalPickDecision{
			WalletID:    best.WalletID,
			Confidence:  0.75,
			Rationale:   fmt.Sprintf("Selected wallet %d (%s, score %d) after %d rounds.", best.WalletID, best.StrategyType, best.SmartScore, in.Round),
			RationaleZh: fmt.Sprintf("经 %d 轮分析选择钱包 %d（%s，评分 %d）。", in.Round, best.WalletID, best.StrategyType, best.SmartScore),
		}
		out.NLSummary = fmt.Sprintf("Final: wallet %d (score %d)", best.WalletID, best.SmartScore)
		out.NLSummaryZh = fmt.Sprintf("最终推荐：钱包 %d（评分 %d）", best.WalletID, best.SmartScore)
	}

	return out, nil
}
