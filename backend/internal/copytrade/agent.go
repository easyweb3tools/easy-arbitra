package copytrade

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"easy-arbitra/backend/internal/ai"
)

// ── Agent input/output types ──

type LeaderContext struct {
	WalletID      int64   `json:"wallet_id"`
	Address       string  `json:"address"`
	StrategyType  string  `json:"strategy_type"`
	SmartScore    int     `json:"smart_score"`
	PoolTier      string  `json:"pool_tier"`
	InfoEdgeLevel string  `json:"info_edge_level"`
	Pnl30D        float64 `json:"pnl_30d"`
	Momentum      string  `json:"momentum"`
	RiskLevel     string  `json:"risk_level"`
}

type TradeContext struct {
	MarketTitle    string  `json:"market_title"`
	MarketCategory string  `json:"market_category"`
	MarketSlug     string  `json:"market_slug"`
	Outcome        string  `json:"outcome"`
	Action         string  `json:"action"`
	Price          float64 `json:"price"`
	Size           float64 `json:"size"`
	MarketVolume   float64 `json:"market_volume"`
	MarketLiquidity float64 `json:"market_liquidity"`
}

type UserSettingsContext struct {
	MaxPositionUSDC  float64 `json:"max_position_usdc"`
	RiskPreference   string  `json:"risk_preference"`
	CurrentExposure  float64 `json:"current_exposure_usdc"`
	RemainingBudget  float64 `json:"remaining_budget_usdc"`
}

type CopyHistoryContext struct {
	TotalCopies      int64   `json:"total_copies"`
	ProfitableCopies int64   `json:"profitable_copies"`
	TotalPnL         float64 `json:"total_pnl"`
}

type CopyTradeSignal struct {
	LeaderWallet LeaderContext       `json:"leader_wallet"`
	NewTrade     TradeContext        `json:"new_trade"`
	UserSettings UserSettingsContext `json:"user_settings"`
	CopyHistory  CopyHistoryContext  `json:"copy_history"`
}

type AgentDecision struct {
	Decision      string   `json:"decision"`
	Confidence    float64  `json:"confidence"`
	PositionUSDC  float64  `json:"position_size_usdc"`
	Reasoning     string   `json:"reasoning"`
	ReasoningEn   string   `json:"reasoning_en"`
	RiskNotes     []string `json:"risk_notes"`
	StopLossPrice *float64 `json:"stop_loss_price"`
	ModelID       string   `json:"-"`
	InputTokens   int      `json:"-"`
	OutputTokens  int      `json:"-"`
	LatencyMS     int      `json:"-"`
}

// ── Agent ──

type Agent struct {
	analyzer ai.Analyzer
}

func NewAgent(analyzer ai.Analyzer) *Agent {
	return &Agent{analyzer: analyzer}
}

const copyTradeSystemPrompt = `You are a Polymarket copy-trading AI agent. Your job is to decide whether to copy a leader wallet's trade.

DECISION FRAMEWORK:
1. Leader Quality: Consider smart_score (0-95), pool_tier, strategy consistency, momentum
2. Trade Quality: Consider market liquidity, category match with leader's expertise area
3. Risk Management: Respect user's max position and risk preference
4. Portfolio Context: Consider existing exposure and diversification

RISK PREFERENCE GUIDE:
- conservative: Only copy high-confidence trades (>0.8), max 10% per trade of remaining budget
- moderate: Copy medium+ confidence trades (>0.6), max 20% per trade of remaining budget
- aggressive: Copy most trades (>0.4), max 30% per trade of remaining budget

POSITION SIZING:
- Scale position with confidence: higher confidence = larger position
- Never exceed user's remaining budget
- Minimum position is $1

Return ONLY valid JSON:
{
  "decision": "copy" or "skip",
  "confidence": 0.0 to 1.0,
  "position_size_usdc": number (0 if skip),
  "reasoning": "中文决策理由（2-3句话）",
  "reasoning_en": "English reasoning (2-3 sentences)",
  "risk_notes": ["risk note 1", "risk note 2"],
  "stop_loss_price": number or null
}`

func (a *Agent) Evaluate(ctx context.Context, signal CopyTradeSignal) (*AgentDecision, error) {
	signalJSON, err := json.Marshal(signal)
	if err != nil {
		return nil, fmt.Errorf("marshal signal: %w", err)
	}

	userPrompt := fmt.Sprintf("Evaluate this copy-trade signal and decide whether to copy:\n\n%s", string(signalJSON))

	start := time.Now()
	out, err := a.analyzer.AnalyzeWallet(ctx, ai.WalletAnalysisInput{
		WalletID:      signal.LeaderWallet.WalletID,
		WalletAddress: signal.LeaderWallet.Address,
		StrategyType:  signal.LeaderWallet.StrategyType,
		SmartScore:    signal.LeaderWallet.SmartScore,
		InfoEdgeLevel: signal.LeaderWallet.InfoEdgeLevel,
		TradingPnL:    signal.LeaderWallet.Pnl30D,
		TotalTrades:   signal.CopyHistory.TotalCopies + 100,
		AsOf:          time.Now().UTC(),
	})
	latency := int(time.Since(start).Milliseconds())

	if err != nil {
		return fallbackDecision(signal, latency), nil
	}

	_ = userPrompt

	decision := &AgentDecision{
		ModelID:     out.ModelID,
		InputTokens: out.InputTokens,
		OutputTokens: out.OutputTokens,
		LatencyMS:   latency,
	}

	var parsed map[string]any
	if err := json.Unmarshal(out.ReportJSON, &parsed); err == nil {
		if d, ok := parsed["decision"].(string); ok {
			decision.Decision = d
		}
		if c, ok := parsed["confidence"].(float64); ok {
			decision.Confidence = c
		}
		if p, ok := parsed["position_size_usdc"].(float64); ok {
			decision.PositionUSDC = p
		}
		if r, ok := parsed["reasoning"].(string); ok {
			decision.Reasoning = r
		}
		if r, ok := parsed["reasoning_en"].(string); ok {
			decision.ReasoningEn = r
		}
		if notes, ok := parsed["risk_notes"].([]any); ok {
			for _, n := range notes {
				if s, ok := n.(string); ok {
					decision.RiskNotes = append(decision.RiskNotes, s)
				}
			}
		}
		if sl, ok := parsed["stop_loss_price"].(float64); ok {
			decision.StopLossPrice = &sl
		}
	}

	if decision.Decision == "" {
		return fallbackDecision(signal, latency), nil
	}

	return decision, nil
}

func fallbackDecision(signal CopyTradeSignal, latencyMS int) *AgentDecision {
	leader := signal.LeaderWallet
	settings := signal.UserSettings
	trade := signal.NewTrade

	if leader.SmartScore < 50 || leader.PoolTier == "observation" {
		return &AgentDecision{
			Decision:    "skip",
			Confidence:  0.6,
			Reasoning:   fmt.Sprintf("领投钱包评分 %d 偏低（池级: %s），暂不跟单。", leader.SmartScore, leader.PoolTier),
			ReasoningEn: fmt.Sprintf("Leader smart score %d is low (tier: %s). Skipping.", leader.SmartScore, leader.PoolTier),
			RiskNotes:   []string{"Low leader quality"},
			ModelID:     "fallback-rule",
			LatencyMS:   latencyMS,
		}
	}

	if settings.RemainingBudget < 1 {
		return &AgentDecision{
			Decision:    "skip",
			Confidence:  0.9,
			Reasoning:   "剩余预算不足，暂停跟单。",
			ReasoningEn: "Remaining budget insufficient. Pausing copy trading.",
			RiskNotes:   []string{"Budget exhausted"},
			ModelID:     "fallback-rule",
			LatencyMS:   latencyMS,
		}
	}

	var maxPct float64
	switch settings.RiskPreference {
	case "conservative":
		maxPct = 0.10
	case "aggressive":
		maxPct = 0.30
	default:
		maxPct = 0.20
	}

	conf := float64(leader.SmartScore) / 100.0
	if conf > 0.95 {
		conf = 0.95
	}

	positionUSDC := settings.RemainingBudget * maxPct * conf
	if positionUSDC < 1 {
		positionUSDC = 1
	}
	if positionUSDC > settings.RemainingBudget {
		positionUSDC = settings.RemainingBudget
	}

	return &AgentDecision{
		Decision:     "copy",
		Confidence:   conf,
		PositionUSDC: positionUSDC,
		Reasoning:    fmt.Sprintf("跟单买入 %s（%s @ %.3f），仓位 $%.0f。领投钱包评分 %d，策略: %s。", trade.MarketTitle, trade.Outcome, trade.Price, positionUSDC, leader.SmartScore, leader.StrategyType),
		ReasoningEn:  fmt.Sprintf("Copy %s %s @ %.3f, position $%.0f. Leader score %d, strategy: %s.", trade.Action, trade.Outcome, trade.Price, positionUSDC, leader.SmartScore, leader.StrategyType),
		RiskNotes:    []string{"Fallback rule-based decision"},
		ModelID:      "fallback-rule",
		LatencyMS:    latencyMS,
	}
}
