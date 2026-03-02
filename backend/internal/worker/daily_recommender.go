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

// DailyRecommender runs once per day (via interval config).
// 1. Backfills yesterday's pick with follow-up PnL
// 2. Picks today's best trader using score leaderboard
// 3. Calls Nova to generate recommendation reasoning
// 4. Writes the DailyPick record
type DailyRecommender struct {
	dailyPickRepo *repository.DailyPickRepository
	scoreRepo     *repository.ScoreRepository
	tradeRepo     *repository.TradeRepository
	walletRepo    *repository.WalletRepository
	analyzer      ai.Analyzer
}

func NewDailyRecommender(
	dailyPickRepo *repository.DailyPickRepository,
	scoreRepo *repository.ScoreRepository,
	tradeRepo *repository.TradeRepository,
	walletRepo *repository.WalletRepository,
	analyzer ai.Analyzer,
) *DailyRecommender {
	return &DailyRecommender{
		dailyPickRepo: dailyPickRepo,
		scoreRepo:     scoreRepo,
		tradeRepo:     tradeRepo,
		walletRepo:    walletRepo,
		analyzer:      analyzer,
	}
}

func (d *DailyRecommender) Name() string { return "daily_recommender" }

func (d *DailyRecommender) RunOnce(ctx context.Context) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Step 1: Backfill yesterday's pick with follow PnL
	_ = d.backfillYesterdayPick(ctx, today) // non-fatal

	// Step 2: Check if today's pick already exists
	if existing, err := d.dailyPickRepo.GetByDate(ctx, today); err == nil && existing != nil {
		return nil // already picked today
	}

	// Step 3: Pick best trader from leaderboard
	rows, _, err := d.scoreRepo.Leaderboard(ctx, 10, 0, "smart_score", "desc")
	if err != nil {
		return fmt.Errorf("leaderboard query: %w", err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("no traders available for daily pick")
	}

	// Find the best candidate that wasn't picked in the last 7 days
	var chosen *repository.LeaderboardRow
	recentPicks, _ := d.dailyPickRepo.ListRecent(ctx, 7)
	recentSet := make(map[int64]bool)
	for _, p := range recentPicks {
		recentSet[p.WalletID] = true
	}
	for i := range rows {
		if !recentSet[rows[i].WalletID] {
			chosen = &rows[i]
			break
		}
	}
	if chosen == nil {
		chosen = &rows[0] // fallback to top if all were recent
	}

	// Get PnL data for the chosen wallet
	pnl, err := d.tradeRepo.AggregateByWalletID(ctx, chosen.WalletID)
	if err != nil {
		return fmt.Errorf("aggregate pnl for wallet %d: %w", chosen.WalletID, err)
	}
	realizedPnL := pnl.TradingPnL + pnl.MakerRebates

	// Step 4: Call Nova for recommendation reasoning
	wallet, err := d.walletRepo.GetByID(ctx, chosen.WalletID)
	if err != nil {
		return fmt.Errorf("get wallet %d: %w", chosen.WalletID, err)
	}

	analysisInput := ai.WalletAnalysisInput{
		WalletID:      chosen.WalletID,
		WalletAddress: polyaddr.BytesToHex(wallet.Address),
		StrategyType:  chosen.StrategyType,
		SmartScore:    chosen.SmartScore,
		InfoEdgeLevel: chosen.InfoEdgeLevel,
		TradingPnL:    pnl.TradingPnL,
		MakerRebates:  pnl.MakerRebates,
		FeesPaid:      pnl.FeesPaid,
		TotalTrades:   pnl.TotalTrades,
		Volume30D:     pnl.Volume30D,
		AsOf:          today,
	}

	var reasonJSON datatypes.JSON
	var reasonSummary string
	var modelID string

	output, err := d.analyzer.AnalyzeWallet(ctx, analysisInput)
	if err != nil {
		// Nova call failed, use fallback summary
		reasonSummary = fmt.Sprintf(
			"Today's recommended trader (score %d, %s strategy) with %.2f USDC realized PnL across %d trades.",
			chosen.SmartScore, chosen.StrategyType, realizedPnL, pnl.TotalTrades,
		)
		fallback := map[string]any{"fallback": true, "reason": reasonSummary}
		reasonJSON, _ = json.Marshal(fallback)
		modelID = "fallback"
	} else {
		reasonJSON = datatypes.JSON(output.ReportJSON)
		reasonSummary = output.NLSummary
		modelID = output.ModelID
	}

	// Step 5: Write daily pick
	pick := &model.DailyPick{
		PickDate:      today,
		WalletID:      chosen.WalletID,
		SmartScore:    chosen.SmartScore,
		RealizedPnL:   realizedPnL,
		TotalTrades:   pnl.TotalTrades,
		WinRate:       0, // win rate requires per-trade analysis; left at 0 for now
		ReasonJSON:    reasonJSON,
		ReasonSummary: reasonSummary,
		ModelID:       modelID,
	}

	return d.dailyPickRepo.Create(ctx, pick)
}

func (d *DailyRecommender) backfillYesterdayPick(ctx context.Context, today time.Time) error {
	yesterday := today.Add(-24 * time.Hour)
	pick, err := d.dailyPickRepo.GetByDate(ctx, yesterday)
	if err != nil {
		return nil // no pick yesterday, nothing to backfill
	}
	if pick.FollowPnL != nil {
		return nil // already backfilled
	}

	// Calculate PnL for the picked wallet now vs at pick time
	pnlNow, err := d.tradeRepo.AggregateByWalletID(ctx, pick.WalletID)
	if err != nil {
		return err
	}
	realizedNow := pnlNow.TradingPnL + pnlNow.MakerRebates

	followPnL := realizedNow - pick.RealizedPnL
	tradesFollowed := int(pnlNow.TotalTrades - pick.TotalTrades)
	if tradesFollowed < 0 {
		tradesFollowed = 0
	}

	return d.dailyPickRepo.UpdateFollowResult(ctx, pick.ID, tradesFollowed, followPnL)
}
