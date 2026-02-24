package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/pkg/polyaddr"
	"gorm.io/gorm"
)

func (s *WalletService) GetDecisionCard(ctx context.Context, walletID int64) (*WalletDecisionCardView, error) {
	if _, err := s.walletRepo.GetByID(ctx, walletID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var score *model.WalletScore
	if row, err := s.scoreRepo.LatestByWalletID(ctx, walletID); err == nil {
		score = row
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var feat *model.WalletFeaturesDaily
	if s.featureRepo != nil {
		if row, err := s.featureRepo.LatestByWalletID(ctx, walletID); err == nil {
			feat = row
		}
	}

	strategy := "unknown"
	scoreValue := 0
	poolTier := "observation"
	updatedAt := time.Now().UTC()
	if score != nil {
		strategy = score.StrategyType
		scoreValue = score.SmartScore
		poolTier = normalizePoolTier(score.PoolTier)
		updatedAt = score.ScoredAt.UTC()
	}

	suitableFor := deriveSuitableFor(valueOr(score, func(v *model.WalletScore) string { return v.SuitableFor }), strategy)
	riskLevel := deriveRiskLevel(valueOr(score, func(v *model.WalletScore) string { return v.RiskLevel }), scoreValue, feat)
	suggestedPosition := deriveSuggestedPosition(valueOr(score, func(v *model.WalletScore) string { return v.SuggestedPosition }), riskLevel)
	momentum := deriveMomentum(valueOr(score, func(v *model.WalletScore) string { return v.Momentum }), feat)

	recommendation := fmt.Sprintf("This wallet is suitable for %s users, with %s risk and %s momentum. Suggested allocation is %s.", suitableFor, riskLevel, momentum, suggestedPosition)
	recommendationZh := fmt.Sprintf("该钱包更适合%s用户，当前风险%s、状态%s，建议仓位区间为%s。", suitableForZH(suitableFor), riskLevelZH(riskLevel), momentumZH(momentum), suggestedPosition)

	return &WalletDecisionCardView{
		WalletID:          walletID,
		PoolTier:          poolTier,
		SuitableFor:       suitableFor,
		RiskLevel:         riskLevel,
		SuggestedPosition: suggestedPosition,
		Momentum:          momentum,
		Status7D:          momentum,
		Recommendation:    recommendation,
		RecommendationZh:  recommendationZh,
		Disclaimer:        "This is a probabilistic judgment for research only, not investment advice.",
		DisclaimerZh:      "以上为概率判断，仅供研究参考，不构成投资建议。",
		LastUpdated:       updatedAt.Format(time.RFC3339),
	}, nil
}

func (s *WalletService) GetShareLanding(ctx context.Context, walletID int64) (*WalletShareLandingView, error) {
	shareCard, err := s.GetShareCard(ctx, walletID)
	if err != nil {
		return nil, err
	}
	decision, err := s.GetDecisionCard(ctx, walletID)
	if err != nil {
		return nil, err
	}

	pnl7d := 0.0
	pnl30d := 0.0
	if s.featureRepo != nil {
		if f, err := s.featureRepo.LatestByWalletID(ctx, walletID); err == nil {
			pnl7d = f.Pnl7d
			pnl30d = f.Pnl30d
		}
	}
	maxDrawdown7D := 0.0
	if pnl7d < 0 {
		maxDrawdown7D = -pnl7d
	}

	summary := shareCard.NLSummary
	if strings.TrimSpace(summary) == "" {
		summary = fallbackSummary(shareCard.StrategyType, shareCard.SmartScore, "", shareCard.TotalTrades, shareCard.RealizedPnL)
	}

	return &WalletShareLandingView{
		Wallet:         shareCard.Wallet,
		PoolTier:       shareCard.PoolTier,
		StrategyType:   shareCard.StrategyType,
		SmartScore:     shareCard.SmartScore,
		Pnl7D:          pnl7d,
		Pnl30D:         pnl30d,
		MaxDrawdown7D:  maxDrawdown7D,
		StabilityScore: shareCard.SmartScore,
		NLSummary:      summary,
		FollowerCount:  shareCard.FollowerCount,
		NewFollowers7D: shareCard.NewFollowers7D,
		DecisionCard:   *decision,
		UpdatedAt:      shareCard.UpdatedAt,
	}, nil
}

func (s *PortfolioService) List(ctx context.Context) ([]PortfolioItem, error) {
	rows, err := s.portfolioRepo.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return s.defaultPortfolios(ctx)
	}
	items := make([]PortfolioItem, 0, len(rows))
	for _, row := range rows {
		walletIDs := parseWalletIDs(row.WalletIDs)
		wallets := loadWalletViews(ctx, s.walletRepo, walletIDs)
		nameZh := row.Name
		if row.NameZh != nil && strings.TrimSpace(*row.NameZh) != "" {
			nameZh = strings.TrimSpace(*row.NameZh)
		}
		desc := ""
		if row.Description != nil {
			desc = strings.TrimSpace(*row.Description)
		}
		expectedReturn, maxDrawdown := portfolioMeta(row.RiskLevel)
		items = append(items, PortfolioItem{
			ID:             row.ID,
			Name:           row.Name,
			NameZh:         nameZh,
			Description:    desc,
			RiskLevel:      normalizeRiskLevel(row.RiskLevel),
			ExpectedReturn: expectedReturn,
			MaxDrawdown:    maxDrawdown,
			WalletIDs:      walletIDs,
			Wallets:        wallets,
		})
	}
	return items, nil
}

func (s *PortfolioService) defaultPortfolios(ctx context.Context) ([]PortfolioItem, error) {
	stableRows, err := s.walletRepo.ListPotentialWallets(ctx, repository.PotentialWalletFilter{
		MinTrades:      100,
		MinRealizedPnL: 0,
		SortBy:         "smart_score",
		Order:          "desc",
		Limit:          3,
	})
	if err != nil {
		return nil, err
	}
	aggressiveRows, err := s.walletRepo.ListPotentialWallets(ctx, repository.PotentialWalletFilter{
		MinTrades:      100,
		MinRealizedPnL: 0,
		SortBy:         "realized_pnl",
		Order:          "desc",
		Limit:          3,
	})
	if err != nil {
		return nil, err
	}

	makePortfolio := func(id int64, name string, nameZh string, desc string, risk string, rows []repository.PotentialWalletRow, expected string, drawdown string) PortfolioItem {
		walletIDs := make([]int64, 0, len(rows))
		wallets := make([]WalletView, 0, len(rows))
		for _, row := range rows {
			walletIDs = append(walletIDs, row.WalletID)
			wallets = append(wallets, WalletView{ID: row.WalletID, Address: polyaddr.BytesToHex(row.Address), Pseudonym: row.Pseudonym, Tracked: row.IsTracked})
		}
		return PortfolioItem{
			ID:             id,
			Name:           name,
			NameZh:         nameZh,
			Description:    desc,
			RiskLevel:      risk,
			ExpectedReturn: expected,
			MaxDrawdown:    drawdown,
			WalletIDs:      walletIDs,
			Wallets:        wallets,
		}
	}

	items := make([]PortfolioItem, 0, 2)
	if len(stableRows) > 0 {
		items = append(items, makePortfolio(
			1,
			"Stable Pack",
			"稳健组合",
			"Low-volatility wallets for conservative follow strategy.",
			"low",
			stableRows,
			"5%-8%",
			"<10%",
		))
	}
	if len(aggressiveRows) > 0 {
		items = append(items, makePortfolio(
			2,
			"Aggressive Pack",
			"进攻组合",
			"Higher-return wallets with higher volatility.",
			"high",
			aggressiveRows,
			"15%-30%",
			"<30%",
		))
	}
	return items, nil
}

func (s *WatchlistService) Summary(ctx context.Context, userFingerprint string) (*WatchlistSummary, error) {
	cleanUser := strings.TrimSpace(userFingerprint)
	if cleanUser == "" {
		return nil, errors.New("empty user fingerprint")
	}

	total, err := s.watchlistRepo.CountByUser(ctx, cleanUser)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &WatchlistSummary{FollowedWallets: 0, StyleDistribution: map[string]int{}, ActionRequired: 0, HealthyWallets: 0}, nil
	}

	rows, err := s.watchlistRepo.ListByUser(ctx, cleanUser, 200, 0)
	if err != nil {
		return nil, err
	}
	styles := map[string]int{}
	for _, row := range rows {
		key := strings.TrimSpace(row.StrategyType)
		if key == "" {
			key = "unknown"
		}
		styles[key]++
	}

	actionRequired, err := s.watchlistRepo.CountActionRequiredByUser(ctx, cleanUser)
	if err != nil {
		return nil, err
	}
	healthy := total - actionRequired
	if healthy < 0 {
		healthy = 0
	}

	return &WatchlistSummary{
		FollowedWallets:   total,
		StyleDistribution: styles,
		ActionRequired:    actionRequired,
		HealthyWallets:    healthy,
	}, nil
}

func fallbackSummary(strategyType string, smartScore int, nlSummary string, tradeCount int64, realizedPnL float64) string {
	if strings.TrimSpace(nlSummary) != "" {
		return strings.TrimSpace(nlSummary)
	}
	risk := "showing elevated volatility"
	switch {
	case smartScore >= 80:
		risk = "showing stable behavior"
	case smartScore >= 60:
		risk = "showing acceptable stability"
	}
	if strings.TrimSpace(strategyType) == "" {
		strategyType = "unknown"
	}
	return fmt.Sprintf("This wallet follows %s style with %d trades and %.2f realized PnL, %s.", strategyType, tradeCount, realizedPnL, risk)
}

func normalizePoolTier(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "star":
		return "star"
	case "strategy":
		return "strategy"
	default:
		return "observation"
	}
}

func normalizeRiskLevel(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "low":
		return "low"
	case "high":
		return "high"
	default:
		return "medium"
	}
}

func deriveSuitableFor(raw string, strategy string) string {
	if v := strings.ToLower(strings.TrimSpace(raw)); v != "" {
		return v
	}
	switch strings.TrimSpace(strategy) {
	case "market_maker", "quant":
		return "conservative"
	case "event_trader", "arbitrage":
		return "event_driven"
	default:
		return "aggressive"
	}
}

func deriveRiskLevel(raw string, smartScore int, feat *model.WalletFeaturesDaily) string {
	if v := strings.ToLower(strings.TrimSpace(raw)); v != "" {
		return normalizeRiskLevel(v)
	}
	if feat != nil {
		if feat.Pnl7d < 0 || feat.Pnl30d < 0 {
			return "high"
		}
		if feat.Pnl30d > 0 && feat.Pnl7d > 0 {
			return "low"
		}
	}
	if smartScore >= 80 {
		return "low"
	}
	if smartScore >= 60 {
		return "medium"
	}
	return "high"
}

func deriveSuggestedPosition(raw string, risk string) string {
	if v := strings.TrimSpace(raw); v != "" {
		return v
	}
	switch risk {
	case "low":
		return "5-10%"
	case "medium":
		return "3-5%"
	default:
		return "1-3%"
	}
}

func deriveMomentum(raw string, feat *model.WalletFeaturesDaily) string {
	if v := strings.ToLower(strings.TrimSpace(raw)); v != "" {
		switch v {
		case "heating", "stable", "cooling":
			return v
		}
	}
	if feat == nil {
		return "stable"
	}
	if feat.Pnl7d > 0 {
		return "heating"
	}
	if feat.Pnl7d < 0 {
		return "cooling"
	}
	return "stable"
}

func suitableForZH(v string) string {
	switch v {
	case "conservative":
		return "稳健型"
	case "event_driven":
		return "事件驱动"
	default:
		return "进攻型"
	}
}

func riskLevelZH(v string) string {
	switch v {
	case "low":
		return "低"
	case "high":
		return "高"
	default:
		return "中"
	}
}

func momentumZH(v string) string {
	switch v {
	case "heating":
		return "升温"
	case "cooling":
		return "降温"
	default:
		return "平稳"
	}
}

func parseWalletIDs(raw []byte) []int64 {
	ids := []int64{}
	if len(raw) == 0 {
		return ids
	}
	var ints []int64
	if err := json.Unmarshal(raw, &ints); err == nil {
		return ints
	}
	var mixed []any
	if err := json.Unmarshal(raw, &mixed); err != nil {
		return ids
	}
	for _, item := range mixed {
		switch v := item.(type) {
		case float64:
			ids = append(ids, int64(v))
		}
	}
	return ids
}

func loadWalletViews(ctx context.Context, walletRepo *repository.WalletRepository, walletIDs []int64) []WalletView {
	rows := make([]WalletView, 0, len(walletIDs))
	for _, walletID := range walletIDs {
		wallet, err := walletRepo.GetByID(ctx, walletID)
		if err != nil {
			continue
		}
		rows = append(rows, WalletView{
			ID:        wallet.ID,
			Address:   polyaddr.BytesToHex(wallet.Address),
			Pseudonym: wallet.Pseudonym,
			Tracked:   wallet.IsTracked,
		})
	}
	return rows
}

func portfolioMeta(risk string) (expectedReturn string, maxDrawdown string) {
	switch normalizeRiskLevel(risk) {
	case "low":
		return "5%-8%", "<10%"
	case "high":
		return "15%-30%", "<30%"
	default:
		return "8%-15%", "<20%"
	}
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func valueOr[T any](row *model.WalletScore, fn func(*model.WalletScore) T) T {
	var zero T
	if row == nil {
		return zero
	}
	return fn(row)
}
