package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"easy-arbitra/backend/internal/ai"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/pkg/polyaddr"

	"gorm.io/datatypes"
)

// NovaOrchestrator is the hourly worker that wakes Nova to analyze candidates.
// Nova drives the entire flow: it evaluates candidates, maintains memory across
// rounds, self-determines when to make the final pick, and writes the result.
type NovaOrchestrator struct {
	sessionRepo   *repository.NovaSessionRepository
	dailyPickRepo *repository.DailyPickRepository
	scoreRepo     *repository.ScoreRepository
	tradeRepo     *repository.TradeRepository
	walletRepo    *repository.WalletRepository
	analyzer      ai.Analyzer
	startHour     int // UTC hour to start daily analysis (e.g. 8)
	endHour       int // UTC hour for last round  (e.g. 22)
}

func NewNovaOrchestrator(
	sessionRepo *repository.NovaSessionRepository,
	dailyPickRepo *repository.DailyPickRepository,
	scoreRepo *repository.ScoreRepository,
	tradeRepo *repository.TradeRepository,
	walletRepo *repository.WalletRepository,
	analyzer ai.Analyzer,
	startHour, endHour int,
) *NovaOrchestrator {
	if startHour < 0 {
		startHour = 8
	}
	if endHour <= startHour {
		endHour = 22
	}
	return &NovaOrchestrator{
		sessionRepo:   sessionRepo,
		dailyPickRepo: dailyPickRepo,
		scoreRepo:     scoreRepo,
		tradeRepo:     tradeRepo,
		walletRepo:    walletRepo,
		analyzer:      analyzer,
		startHour:     startHour,
		endHour:       endHour,
	}
}

func (o *NovaOrchestrator) Name() string { return "nova_orchestrator" }

func (o *NovaOrchestrator) RunOnce(ctx context.Context) error {
	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour)
	currentHour := now.Hour()

	// Step 0: Backfill yesterday's daily pick follow PnL (non-fatal)
	_ = o.backfillYesterdayPick(ctx, today)

	// Check if within analysis window
	if currentHour < o.startHour || currentHour > o.endHour {
		return nil // outside analysis window
	}

	// Step 1: Already have a final pick today? Skip.
	hasFinal, _ := o.sessionRepo.HasFinalByDate(ctx, today)
	if hasFinal {
		return nil
	}

	// Step 2: Determine round number
	round := currentHour - o.startHour + 1
	totalRounds := o.endHour - o.startHour + 1
	isLastRound := currentHour >= o.endHour

	// Step 3: Collect top 20 candidates
	candidates, err := o.collectCandidates(ctx, 20)
	if err != nil {
		return fmt.Errorf("collect candidates: %w", err)
	}
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates available")
	}

	// Step 4: Load Nova's memory (previous rounds today)
	memory, err := o.loadMemory(ctx, today)
	if err != nil {
		return fmt.Errorf("load memory: %w", err)
	}

	// Step 5: Load yesterday's result for feedback
	yesterdayResult := o.loadYesterdayResult(ctx, today)

	// Step 6: Call Nova
	input := ai.OrchestrateInput{
		CurrentTime:     now,
		Round:           round,
		TotalRounds:     totalRounds,
		IsLastRound:     isLastRound,
		Candidates:      candidates,
		Memory:          memory,
		YesterdayResult: yesterdayResult,
	}

	output, err := o.analyzer.Orchestrate(ctx, input)
	if err != nil {
		return fmt.Errorf("nova orchestrate: %w", err)
	}

	// Step 7: Save session (Nova's memory for this round)
	candidatesJSON, _ := json.Marshal(output.Rankings)
	observationsJSON, _ := json.Marshal(map[string]string{"notes": output.Observations})
	decisionJSON, _ := json.Marshal(output)

	session := &model.NovaSession{
		SessionDate:      today,
		Round:            round,
		Phase:            output.Phase,
		CandidatesJSON:   datatypes.JSON(candidatesJSON),
		ObservationsJSON: datatypes.JSON(observationsJSON),
		DecisionJSON:     datatypes.JSON(decisionJSON),
		NLSummary:        output.NLSummary,
		NLSummaryZh:      output.NLSummaryZh,
		ModelID:          output.ModelID,
		InputTokens:      output.InputTokens,
		OutputTokens:     output.OutputTokens,
		LatencyMS:        output.LatencyMS,
	}

	// Step 8: If Nova decided "final", write the daily pick
	if output.Phase == "final" && output.FinalPick != nil {
		walletID := output.FinalPick.WalletID
		session.PickedWalletID = &walletID

		// Get PnL for the chosen wallet
		pnl, pnlErr := o.tradeRepo.AggregateByWalletID(ctx, walletID)
		realizedPnL := 0.0
		totalTrades := int64(0)
		if pnlErr == nil {
			realizedPnL = pnl.TradingPnL + pnl.MakerRebates
			totalTrades = pnl.TotalTrades
		}

		// Find the wallet's score from candidates
		smartScore := 0
		for _, c := range candidates {
			if c.WalletID == walletID {
				smartScore = c.SmartScore
				break
			}
		}

		pick := &model.DailyPick{
			PickDate:        today,
			WalletID:        walletID,
			SmartScore:      smartScore,
			RealizedPnL:     realizedPnL,
			TotalTrades:     totalTrades,
			WinRate:         output.FinalPick.Confidence,
			ReasonJSON:      datatypes.JSON(decisionJSON),
			ReasonSummary:   output.FinalPick.Rationale,
			ReasonSummaryZh: output.FinalPick.RationaleZh,
			ModelID:         output.ModelID,
		}
		if err := o.dailyPickRepo.Create(ctx, pick); err != nil {
			return fmt.Errorf("create daily pick: %w", err)
		}
	}

	return o.sessionRepo.Create(ctx, session)
}

func (o *NovaOrchestrator) collectCandidates(ctx context.Context, limit int) ([]ai.CandidateData, error) {
	rows, _, err := o.scoreRepo.Leaderboard(ctx, limit, 0, "smart_score", "desc")
	if err != nil {
		return nil, err
	}

	candidates := make([]ai.CandidateData, 0, len(rows))
	for _, r := range rows {
		wallet, wErr := o.walletRepo.GetByID(ctx, r.WalletID)
		if wErr != nil {
			continue
		}
		pnl, pErr := o.tradeRepo.AggregateByWalletID(ctx, r.WalletID)
		if pErr != nil {
			continue
		}
		candidates = append(candidates, ai.CandidateData{
			WalletID:      r.WalletID,
			Address:       polyaddr.BytesToHex(wallet.Address),
			SmartScore:    r.SmartScore,
			StrategyType:  r.StrategyType,
			InfoEdgeLevel: r.InfoEdgeLevel,
			TradingPnL:    pnl.TradingPnL,
			MakerRebates:  pnl.MakerRebates,
			FeesPaid:      pnl.FeesPaid,
			TotalTrades:   pnl.TotalTrades,
			Volume30D:     pnl.Volume30D,
		})
	}
	return candidates, nil
}

func (o *NovaOrchestrator) loadMemory(ctx context.Context, today time.Time) ([]ai.SessionMemory, error) {
	sessions, err := o.sessionRepo.ListByDate(ctx, today)
	if err != nil {
		return nil, err
	}
	memory := make([]ai.SessionMemory, 0, len(sessions))
	for _, s := range sessions {
		var obs struct {
			Notes string `json:"notes"`
		}
		_ = json.Unmarshal(s.ObservationsJSON, &obs)

		topPick := ""
		if s.PickedWalletID != nil {
			topPick = fmt.Sprintf("wallet_%d", *s.PickedWalletID)
		}
		memory = append(memory, ai.SessionMemory{
			Round:        s.Round,
			Phase:        s.Phase,
			Observations: obs.Notes,
			TopPick:      topPick,
		})
	}
	return memory, nil
}

func (o *NovaOrchestrator) loadYesterdayResult(ctx context.Context, today time.Time) *ai.YesterdayResult {
	yesterday := today.Add(-24 * time.Hour)
	pick, err := o.dailyPickRepo.GetByDate(ctx, yesterday)
	if err != nil || pick == nil || pick.FollowPnL == nil {
		return nil
	}

	wallet, err := o.walletRepo.GetByID(ctx, pick.WalletID)
	if err != nil {
		return nil
	}

	return &ai.YesterdayResult{
		WalletID:       pick.WalletID,
		Address:        polyaddr.BytesToHex(wallet.Address),
		FollowPnL:      *pick.FollowPnL,
		TradesFollowed: pick.TradesFollowed,
		SmartScore:     pick.SmartScore,
	}
}

func (o *NovaOrchestrator) backfillYesterdayPick(ctx context.Context, today time.Time) error {
	yesterday := today.Add(-24 * time.Hour)
	pick, err := o.dailyPickRepo.GetByDate(ctx, yesterday)
	if err != nil || pick == nil {
		return nil
	}
	if pick.FollowPnL != nil {
		return nil // already backfilled
	}

	pnlNow, err := o.tradeRepo.AggregateByWalletID(ctx, pick.WalletID)
	if err != nil {
		return err
	}
	realizedNow := pnlNow.TradingPnL + pnlNow.MakerRebates
	followPnL := realizedNow - pick.RealizedPnL
	tradesFollowed := int(pnlNow.TotalTrades - pick.TotalTrades)
	if tradesFollowed < 0 {
		tradesFollowed = 0
	}

	return o.dailyPickRepo.UpdateFollowResult(ctx, pick.ID, tradesFollowed, followPnL)
}
