package ai

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ── Orchestrate types: Nova-as-Brain ──

// CandidateData is a single candidate wallet for Nova to evaluate.
type CandidateData struct {
	WalletID      int64   `json:"wallet_id"`
	Address       string  `json:"address"`
	SmartScore    int     `json:"smart_score"`
	StrategyType  string  `json:"strategy_type"`
	InfoEdgeLevel string  `json:"info_edge_level"`
	TradingPnL    float64 `json:"trading_pnl"`
	MakerRebates  float64 `json:"maker_rebates"`
	FeesPaid      float64 `json:"fees_paid"`
	TotalTrades   int64   `json:"total_trades"`
	Volume30D     float64 `json:"volume_30d"`
}

// SessionMemory is a summary of a previous round for Nova to recall.
type SessionMemory struct {
	Round        int    `json:"round"`
	Phase        string `json:"phase"`
	Observations string `json:"observations"`
	TopPick      string `json:"top_pick,omitempty"`
}

// YesterdayResult is validation feedback from yesterday's pick.
type YesterdayResult struct {
	WalletID       int64   `json:"wallet_id"`
	Address        string  `json:"address"`
	FollowPnL      float64 `json:"follow_pnl"`
	TradesFollowed int     `json:"trades_followed"`
	SmartScore     int     `json:"smart_score"`
}

// OrchestrateInput is everything Nova needs to make decisions.
type OrchestrateInput struct {
	CurrentTime     time.Time        `json:"current_time"`
	Round           int              `json:"round"`
	TotalRounds     int              `json:"total_rounds"`
	IsLastRound     bool             `json:"is_last_round"`
	Candidates      []CandidateData  `json:"candidates"`
	Memory          []SessionMemory  `json:"memory"`
	YesterdayResult *YesterdayResult `json:"yesterday_result,omitempty"`
}

// CandidateRank is Nova's evaluation of a single candidate.
type CandidateRank struct {
	WalletID int64   `json:"wallet_id"`
	Rank     int     `json:"rank"`
	Score    float64 `json:"score"`
	Reason   string  `json:"reason"`
}

// FinalPickDecision is Nova's final recommendation.
type FinalPickDecision struct {
	WalletID    int64   `json:"wallet_id"`
	Confidence  float64 `json:"confidence"`
	Rationale   string  `json:"rationale"`
	RationaleZh string  `json:"rationale_zh"`
}

// OrchestrateOutput is Nova's response.
type OrchestrateOutput struct {
	Phase        string             `json:"phase"` // "analyzing" | "final"
	Rankings     []CandidateRank    `json:"rankings"`
	Observations string             `json:"observations"`
	FinalPick    *FinalPickDecision `json:"final_pick,omitempty"`
	NLSummary    string             `json:"nl_summary"`
	NLSummaryZh  string             `json:"nl_summary_zh"`
	ModelID      string             `json:"model_id"`
	InputTokens  int                `json:"input_tokens"`
	OutputTokens int                `json:"output_tokens"`
	LatencyMS    int                `json:"latency_ms"`
}

// buildOrchestratePrompts creates the system + user prompts for Nova orchestration.
func BuildOrchestratePrompts(in OrchestrateInput) (string, string) {
	systemPrompt := `You are the analytical brain of a Polymarket trading intelligence system.
You are called every hour to analyze candidate wallets and decide if you are ready to make your "Daily Pick" — the single best trader to recommend today.

Your responsibilities:
1. Evaluate all candidate wallets based on their trading performance, strategy, and risk
2. Maintain observations across rounds (your memory from previous rounds is provided)
3. Decide whether to continue analyzing ("analyzing") or make your final pick ("final")
4. When making a final pick, provide detailed rationale in both English and Chinese

Rules:
- You MUST return strict JSON matching the schema below
- You CAN make a "final" decision early if you are confident enough
- If this is the last round (is_last_round=true), you MUST make a "final" decision
- Learn from yesterday's result if provided — avoid repeating poor picks
- Consider: PnL consistency, trade volume, strategy type, risk-adjusted returns

Return JSON with these fields:
{
  "phase": "analyzing" | "final",
  "rankings": [{"wallet_id": N, "rank": N, "score": 0-100, "reason": "..."}],
  "observations": "your analysis notes for memory (will be recalled next round)",
  "final_pick": {"wallet_id": N, "confidence": 0-1, "rationale": "...", "rationale_zh": "..."} // only when phase="final"
  "nl_summary": "one-line English summary of this round",
  "nl_summary_zh": "本轮分析的中文摘要"
}`

	candidatesJSON, _ := json.MarshalIndent(in.Candidates, "  ", "  ")
	memoryJSON, _ := json.MarshalIndent(in.Memory, "  ", "  ")

	var yesterdayBlock string
	if in.YesterdayResult != nil {
		yJSON, _ := json.MarshalIndent(in.YesterdayResult, "  ", "  ")
		yesterdayBlock = fmt.Sprintf(`
Yesterday's Pick Result (learn from this):
%s
`, string(yJSON))
	} else {
		yesterdayBlock = "Yesterday's Pick Result: None (first day or no result yet)"
	}

	userPrompt := fmt.Sprintf(`Current Time: %s
Round: %d / %d (is_last_round: %v)

%s

Candidate Wallets (%d total):
%s

Your Memory from Previous Rounds (%d rounds):
%s

Based on the above, analyze the candidates and respond with your decision.`,
		in.CurrentTime.UTC().Format(time.RFC3339),
		in.Round, in.TotalRounds, in.IsLastRound,
		yesterdayBlock,
		len(in.Candidates), string(candidatesJSON),
		len(in.Memory), string(memoryJSON),
	)

	return systemPrompt, userPrompt
}

// ParseOrchestrateResponse parses Nova's JSON response into OrchestrateOutput.
func ParseOrchestrateResponse(raw string, modelID string, inputTokens, outputTokens, latencyMS int) (*OrchestrateOutput, error) {
	// Extract JSON from response (Nova may wrap in markdown code blocks)
	cleaned := raw
	if idx := strings.Index(cleaned, "{"); idx >= 0 {
		cleaned = cleaned[idx:]
	}
	if idx := strings.LastIndex(cleaned, "}"); idx >= 0 {
		cleaned = cleaned[:idx+1]
	}

	var out OrchestrateOutput
	if err := json.Unmarshal([]byte(cleaned), &out); err != nil {
		return nil, fmt.Errorf("parse nova orchestrate response: %w", err)
	}

	out.ModelID = modelID
	out.InputTokens = inputTokens
	out.OutputTokens = outputTokens
	out.LatencyMS = latencyMS

	// Validate phase
	if out.Phase != "analyzing" && out.Phase != "final" {
		out.Phase = "analyzing"
	}

	// If final but no pick, force analyzing
	if out.Phase == "final" && out.FinalPick == nil {
		out.Phase = "analyzing"
	}

	return &out, nil
}
