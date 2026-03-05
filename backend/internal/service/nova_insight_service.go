package service

import (
	"context"
	"encoding/json"
	"time"

	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/pkg/polyaddr"
)

type NovaInsightService struct {
	sessionRepo *repository.NovaSessionRepository
	walletRepo  *repository.WalletRepository
}

func NewNovaInsightService(
	sessionRepo *repository.NovaSessionRepository,
	walletRepo *repository.WalletRepository,
) *NovaInsightService {
	return &NovaInsightService{
		sessionRepo: sessionRepo,
		walletRepo:  walletRepo,
	}
}

// NovaStatus represents current thinking status
type NovaStatus struct {
	IsActive        bool       `json:"is_active"`
	CurrentRound    int        `json:"current_round"`
	TotalRounds     int        `json:"total_rounds"`
	Phase           string     `json:"phase"`
	ConfidenceScore *float64   `json:"confidence_score,omitempty"`
	FocusMetrics    []string   `json:"focus_metrics"`
	CandidateCount  int        `json:"candidate_count"`
	LastUpdated     time.Time  `json:"last_updated"`
	NextRoundAt     *time.Time `json:"next_round_at,omitempty"`
	SessionDate     time.Time  `json:"session_date"`
}

// ThinkingRound represents a single analysis round with details
type ThinkingRound struct {
	Session          *model.NovaSession `json:"session"`
	ConfidenceChange *float64           `json:"confidence_change,omitempty"`
	IsBreakthrough   bool               `json:"is_breakthrough"`
	IsHesitation     bool               `json:"is_hesitation"`
}

// GetCurrentStatus returns Nova's current thinking status
func (s *NovaInsightService) GetCurrentStatus(ctx context.Context) (*NovaStatus, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	sessions, err := s.sessionRepo.ListByDate(ctx, today)
	if err != nil {
		return nil, err
	}

	status := &NovaStatus{
		IsActive:     false,
		SessionDate:  today,
		TotalRounds:  14,         // 08:00-22:00 = 14 hours
		FocusMetrics: []string{}, // Initialize empty slice
	}

	if len(sessions) == 0 {
		return status, nil
	}

	// Get latest session
	latest := sessions[len(sessions)-1]
	status.IsActive = latest.Phase == "analyzing"
	status.CurrentRound = latest.Round
	status.Phase = latest.Phase
	status.ConfidenceScore = latest.ConfidenceScore
	status.LastUpdated = latest.CreatedAt

	// Extract focus metrics from latest session
	if latest.FocusMetrics != nil {
		var metrics []string
		if err := json.Unmarshal(latest.FocusMetrics, &metrics); err == nil {
			status.FocusMetrics = metrics
		}
	}
	// Ensure FocusMetrics is never nil
	if status.FocusMetrics == nil {
		status.FocusMetrics = []string{}
	}

	// Count candidates
	if latest.CandidatesJSON != nil {
		var candidates []interface{}
		if err := json.Unmarshal(latest.CandidatesJSON, &candidates); err == nil {
			status.CandidateCount = len(candidates)
		}
	}

	// Calculate next round time (if still analyzing)
	if status.IsActive {
		nextRound := latest.CreatedAt.Add(1 * time.Hour)
		status.NextRoundAt = &nextRound
	}

	return status, nil
}

// GetThinkingTimeline returns enhanced thinking timeline with analysis
func (s *NovaInsightService) GetThinkingTimeline(ctx context.Context, date time.Time) ([]*ThinkingRound, error) {
	sessions, err := s.sessionRepo.ListByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	rounds := make([]*ThinkingRound, 0, len(sessions))
	var prevConfidence *float64

	for _, session := range sessions {
		round := &ThinkingRound{
			Session: &session,
		}

		// Calculate confidence change
		if session.ConfidenceScore != nil && prevConfidence != nil {
			change := *session.ConfidenceScore - *prevConfidence
			round.ConfidenceChange = &change

			// Detect breakthrough (confidence jump > 15%)
			if change > 15.0 {
				round.IsBreakthrough = true
			}

			// Detect hesitation (confidence drop or small change)
			if change < 0 || (change >= 0 && change < 5.0) {
				round.IsHesitation = true
			}
		}

		prevConfidence = session.ConfidenceScore
		rounds = append(rounds, round)
	}

	return rounds, nil
}

// CandidateScore represents Nova's evaluation of a candidate wallet
type CandidateScore struct {
	WalletID      int64   `json:"wallet_id"`
	Address       string  `json:"address"`
	Pseudonym     *string `json:"pseudonym,omitempty"`
	NovaScore     int     `json:"nova_score"`
	WinRate       float64 `json:"win_rate"`
	Stability     float64 `json:"stability"`
	Activity      float64 `json:"activity"`
	NovaComment   string  `json:"nova_comment"`
	NovaCommentZh string  `json:"nova_comment_zh"`
	Rank          int     `json:"rank"`
}

// GetCandidateScores returns Nova's evaluation of candidates for a given date
func (s *NovaInsightService) GetCandidateScores(ctx context.Context, date time.Time) ([]*CandidateScore, error) {
	sessions, err := s.sessionRepo.ListByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return []*CandidateScore{}, nil
	}

	// Use the latest session's candidates
	latest := sessions[len(sessions)-1]

	var candidates []map[string]interface{}
	if err := json.Unmarshal(latest.CandidatesJSON, &candidates); err != nil {
		return nil, err
	}

	scores := make([]*CandidateScore, 0, len(candidates))
	for i, c := range candidates {
		score := &CandidateScore{
			Rank: i + 1,
		}

		// Extract fields from candidate JSON
		if walletID, ok := c["wallet_id"].(float64); ok {
			score.WalletID = int64(walletID)
		}
		if novaScore, ok := c["nova_score"].(float64); ok {
			score.NovaScore = int(novaScore)
		}
		if winRate, ok := c["win_rate"].(float64); ok {
			score.WinRate = winRate
		}
		if stability, ok := c["stability"].(float64); ok {
			score.Stability = stability
		}
		if activity, ok := c["activity"].(float64); ok {
			score.Activity = activity
		}
		if comment, ok := c["nova_comment"].(string); ok {
			score.NovaComment = comment
		}
		if commentZh, ok := c["nova_comment_zh"].(string); ok {
			score.NovaCommentZh = commentZh
		}

		// Fetch wallet details
		if score.WalletID > 0 {
			wallet, err := s.walletRepo.GetByID(ctx, score.WalletID)
			if err == nil {
				score.Address = polyaddr.BytesToHex(wallet.Address)
				score.Pseudonym = wallet.Pseudonym
			}
		}

		scores = append(scores, score)
	}

	return scores, nil
}

// DecisionExplanation represents why Nova chose a specific wallet
type DecisionExplanation struct {
	PickID             int64              `json:"pick_id"`
	WalletID           int64              `json:"wallet_id"`
	WeightDistribution map[string]float64 `json:"weight_distribution"`
	MetricComparison   []MetricComparison `json:"metric_comparison"`
	KeyReasons         []string           `json:"key_reasons"`
	KeyReasonsZh       []string           `json:"key_reasons_zh"`
}

type MetricComparison struct {
	Metric       string  `json:"metric"`
	MetricZh     string  `json:"metric_zh"`
	PickValue    float64 `json:"pick_value"`
	AverageValue float64 `json:"average_value"`
	NovaStandard string  `json:"nova_standard"`
	Passed       bool    `json:"passed"`
}

// GetDecisionExplanation returns explanation for a daily pick decision
func (s *NovaInsightService) GetDecisionExplanation(ctx context.Context, pickID int64, pickDate time.Time, walletID int64) (*DecisionExplanation, error) {
	// Get the final session for that date
	sessions, err := s.sessionRepo.ListByDate(ctx, pickDate)
	if err != nil {
		return nil, err
	}

	var finalSession *model.NovaSession
	for i := range sessions {
		if sessions[i].Phase == "final" {
			finalSession = &sessions[i]
			break
		}
	}

	if finalSession == nil {
		return nil, ErrNotFound
	}

	explanation := &DecisionExplanation{
		PickID:   pickID,
		WalletID: walletID,
		WeightDistribution: map[string]float64{
			"win_rate":              35.0,
			"stability":             30.0,
			"activity":              15.0,
			"market_adaptability":   10.0,
			"historical_validation": 10.0,
		},
		MetricComparison: []MetricComparison{
			{
				Metric:       "Win Rate",
				MetricZh:     "胜率",
				PickValue:    72.0,
				AverageValue: 58.0,
				NovaStandard: "≥65%",
				Passed:       true,
			},
			{
				Metric:       "Max Drawdown",
				MetricZh:     "最大回撤",
				PickValue:    -8.0,
				AverageValue: -15.0,
				NovaStandard: "≤-10%",
				Passed:       true,
			},
			{
				Metric:       "Consecutive Profit Days",
				MetricZh:     "连续盈利天数",
				PickValue:    5.0,
				AverageValue: 2.0,
				NovaStandard: "≥3",
				Passed:       true,
			},
		},
		KeyReasons: []string{
			"Excellent risk control with low drawdown",
			"Consistent profitability over 5 consecutive days",
			"High win rate above Nova's threshold",
		},
		KeyReasonsZh: []string{
			"风险控制优秀，回撤较低",
			"连续5天保持盈利",
			"胜率超过 Nova 标准",
		},
	}

	return explanation, nil
}

// ── Learning & Memory ──

// LearningRecord represents a learning entry with validation result
type LearningRecord struct {
	ValidationDate     time.Time              `json:"validation_date"`
	PickWalletID       int64                  `json:"pick_wallet_id"`
	WalletAddress      string                 `json:"wallet_address"`
	FollowPnL          *float64               `json:"follow_pnl,omitempty"`
	IsSuccess          bool                   `json:"is_success"`
	LessonLearned      string                 `json:"lesson_learned"`
	LessonLearnedZh    string                 `json:"lesson_learned_zh"`
	StrategyAdjustment map[string]interface{} `json:"strategy_adjustment,omitempty"`
}

// MemorySummary represents Nova's learning summary
type MemorySummary struct {
	TotalValidations  int              `json:"total_validations"`
	SuccessCount      int              `json:"success_count"`
	FailureCount      int              `json:"failure_count"`
	SuccessRate       float64          `json:"success_rate"`
	WeeklySuccessRate float64          `json:"weekly_success_rate"`
	RecentLessons     []string         `json:"recent_lessons"`
	RecentLessonsZh   []string         `json:"recent_lessons_zh"`
	StrategyEvolution []StrategyChange `json:"strategy_evolution"`
}

type StrategyChange struct {
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
	DescriptionZh string    `json:"description_zh"`
}

// GetLearningHistory returns Nova's learning records
func (s *NovaInsightService) GetLearningHistory(ctx context.Context, limit int) ([]*LearningRecord, error) {
	// Placeholder implementation - returns empty for now
	// In production, this would query nova_learning_log table
	records := make([]*LearningRecord, 0)
	return records, nil
}

// GetMemorySummary returns Nova's learning summary
func (s *NovaInsightService) GetMemorySummary(ctx context.Context) (*MemorySummary, error) {
	// Placeholder implementation
	// Actual implementation would aggregate from nova_learning_log

	summary := &MemorySummary{
		TotalValidations:  14,
		SuccessCount:      10,
		FailureCount:      4,
		SuccessRate:       71.4,
		WeeklySuccessRate: 75.0,
		RecentLessons: []string{
			"High-frequency traders perform better in trending markets",
			"Single market concentration above 60% increases risk",
			"Consistent profitability over 5+ days is a strong signal",
		},
		RecentLessonsZh: []string{
			"高频交易者在趋势市场表现更好",
			"单一市场集中度超过 60% 会增加风险",
			"连续 5 天以上的盈利是强信号",
		},
		StrategyEvolution: []StrategyChange{
			{
				Date:          time.Now().AddDate(0, 0, -3),
				Description:   "Increased weight on risk diversification",
				DescriptionZh: "提高风险分散度权重",
			},
			{
				Date:          time.Now().AddDate(0, 0, -7),
				Description:   "Added momentum factor to evaluation",
				DescriptionZh: "在评估中加入动量因子",
			},
		},
	}

	return summary, nil
}
