package copytrade

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/pkg/polyaddr"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("copy trading config already exists")
	ErrBudgetExhausted = errors.New("remaining budget insufficient")
)

// ── View types ──

type ConfigView struct {
	ID              int64   `json:"id"`
	WalletID        int64   `json:"wallet_id"`
	WalletAddress   string  `json:"wallet_address"`
	WalletPseudonym *string `json:"wallet_pseudonym,omitempty"`
	Enabled         bool    `json:"enabled"`
	MaxPositionUSDC float64 `json:"max_position_usdc"`
	RiskPreference  string  `json:"risk_preference"`
	TotalPnL        float64 `json:"total_pnl"`
	TotalCopies     int64   `json:"total_copies"`
	OpenPositions   int64   `json:"open_positions"`
	CreatedAt       string  `json:"created_at"`
}

type DecisionView struct {
	ID            int64    `json:"id"`
	Decision      string   `json:"decision"`
	Confidence    float64  `json:"confidence"`
	MarketTitle   string   `json:"market_title"`
	Outcome       string   `json:"outcome"`
	Action        string   `json:"action"`
	Price         float64  `json:"price"`
	SizeUSDC      float64  `json:"size_usdc"`
	StopLossPrice *float64 `json:"stop_loss_price,omitempty"`
	Reasoning     string   `json:"reasoning"`
	ReasoningEn   string   `json:"reasoning_en"`
	RiskNotes     []string `json:"risk_notes"`
	Status        string   `json:"status"`
	RealizedPnL   *float64 `json:"realized_pnl,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

type DashboardView struct {
	TotalPnL       float64        `json:"total_pnl"`
	WinRate        float64        `json:"win_rate"`
	TotalCopies    int64          `json:"total_copies"`
	TotalSkipped   int64          `json:"total_skipped"`
	OpenPositions  int64          `json:"open_positions"`
	ActiveConfigs  int            `json:"active_configs"`
	Configs        []ConfigView   `json:"configs"`
	RecentDecisions []DecisionView `json:"recent_decisions"`
}

type PerformanceView struct {
	TotalPnL    float64          `json:"total_pnl"`
	WinRate     float64          `json:"win_rate"`
	TotalCopies int64            `json:"total_copies"`
	DailyPoints []DailyPerfPoint `json:"daily_points"`
}

type DailyPerfPoint struct {
	Date        string  `json:"date"`
	PnL         float64 `json:"pnl"`
	Copies      int     `json:"copies"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type DecisionListResult struct {
	Items      []DecisionView `json:"items"`
	Pagination Pagination     `json:"pagination"`
}

// ── Service ──

type Service struct {
	repo        *Repository
	walletRepo  *repository.WalletRepository
	scoreRepo   *repository.ScoreRepository
	featureRepo *repository.FeatureRepository
	tradeRepo   *repository.TradeRepository
	marketRepo  *repository.MarketRepository
	agent       *Agent
}

func NewService(
	repo *Repository,
	walletRepo *repository.WalletRepository,
	scoreRepo *repository.ScoreRepository,
	featureRepo *repository.FeatureRepository,
	tradeRepo *repository.TradeRepository,
	marketRepo *repository.MarketRepository,
	agent *Agent,
) *Service {
	return &Service{
		repo: repo, walletRepo: walletRepo, scoreRepo: scoreRepo,
		featureRepo: featureRepo, tradeRepo: tradeRepo, marketRepo: marketRepo, agent: agent,
	}
}

// ── Config management ──

func (s *Service) EnableCopyTrading(ctx context.Context, userFP string, walletID int64, maxPosition float64, riskPref string) (*ConfigView, error) {
	userFP = strings.TrimSpace(userFP)
	if userFP == "" {
		return nil, errors.New("empty user fingerprint")
	}
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if maxPosition <= 0 {
		maxPosition = 1000
	}
	if maxPosition > 100000 {
		maxPosition = 100000
	}
	riskPref = strings.TrimSpace(riskPref)
	if riskPref != "conservative" && riskPref != "aggressive" {
		riskPref = "moderate"
	}
	cfg := &model.CopyTradingConfig{
		UserFingerprint: userFP,
		WalletID:        walletID,
		Enabled:         true,
		MaxPositionUSDC: maxPosition,
		RiskPreference:  riskPref,
	}
	if err := s.repo.UpsertConfig(ctx, cfg); err != nil {
		return nil, err
	}
	saved, err := s.repo.GetConfig(ctx, userFP, walletID)
	if err != nil {
		return nil, err
	}
	return &ConfigView{
		ID:              saved.ID,
		WalletID:        saved.WalletID,
		WalletAddress:   polyaddr.BytesToHex(wallet.Address),
		WalletPseudonym: wallet.Pseudonym,
		Enabled:         saved.Enabled,
		MaxPositionUSDC: saved.MaxPositionUSDC,
		RiskPreference:  saved.RiskPreference,
		CreatedAt:       saved.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

func (s *Service) DisableCopyTrading(ctx context.Context, userFP string, walletID int64) error {
	userFP = strings.TrimSpace(userFP)
	if userFP == "" {
		return errors.New("empty user fingerprint")
	}
	return s.repo.DisableConfig(ctx, userFP, walletID)
}

func (s *Service) UpdateSettings(ctx context.Context, userFP string, walletID int64, maxPosition float64, riskPref string) (*ConfigView, error) {
	return s.EnableCopyTrading(ctx, userFP, walletID, maxPosition, riskPref)
}

func (s *Service) GetConfig(ctx context.Context, userFP string, walletID int64) (*ConfigView, error) {
	cfg, err := s.repo.GetConfig(ctx, userFP, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	wallet, err := s.walletRepo.GetByID(ctx, cfg.WalletID)
	if err != nil {
		return nil, err
	}
	openCount, _ := s.repo.CountExecutedByConfig(ctx, cfg.ID)
	stats, _ := s.repo.GetDashboardStats(ctx, userFP)
	var totalPnL float64
	if stats != nil {
		totalPnL = stats.TotalPnL
	}
	return &ConfigView{
		ID:              cfg.ID,
		WalletID:        cfg.WalletID,
		WalletAddress:   polyaddr.BytesToHex(wallet.Address),
		WalletPseudonym: wallet.Pseudonym,
		Enabled:         cfg.Enabled,
		MaxPositionUSDC: cfg.MaxPositionUSDC,
		RiskPreference:  cfg.RiskPreference,
		TotalPnL:        totalPnL,
		OpenPositions:   openCount,
		TotalCopies:     stats.TotalCopies,
		CreatedAt:       cfg.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

func (s *Service) ListConfigs(ctx context.Context, userFP string) ([]ConfigView, error) {
	userFP = strings.TrimSpace(userFP)
	if userFP == "" {
		return nil, errors.New("empty user fingerprint")
	}
	configs, err := s.repo.ListConfigsByUser(ctx, userFP)
	if err != nil {
		return nil, err
	}
	views := make([]ConfigView, 0, len(configs))
	for _, cfg := range configs {
		wallet, err := s.walletRepo.GetByID(ctx, cfg.WalletID)
		if err != nil {
			continue
		}
		openCount, _ := s.repo.CountExecutedByConfig(ctx, cfg.ID)
		views = append(views, ConfigView{
			ID:              cfg.ID,
			WalletID:        cfg.WalletID,
			WalletAddress:   polyaddr.BytesToHex(wallet.Address),
			WalletPseudonym: wallet.Pseudonym,
			Enabled:         cfg.Enabled,
			MaxPositionUSDC: cfg.MaxPositionUSDC,
			RiskPreference:  cfg.RiskPreference,
			OpenPositions:   openCount,
			CreatedAt:       cfg.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	return views, nil
}

// ── Decision processing ──

func (s *Service) ProcessNewTrade(ctx context.Context, configID int64, leaderTrade repository.TradeHistoryRow, market *model.Market) (*DecisionView, error) {
	cfg, err := s.repo.GetConfigByID(ctx, configID)
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}

	// Pre-check: frequency limit
	if market != nil {
		recent, _ := s.repo.HasRecentCopyInMarket(ctx, configID, market.ID, time.Hour)
		if recent {
			return nil, nil
		}
	}

	// Gather leader context
	wallet, err := s.walletRepo.GetByID(ctx, cfg.WalletID)
	if err != nil {
		return nil, err
	}
	leaderCtx := LeaderContext{
		WalletID: cfg.WalletID,
		Address:  polyaddr.BytesToHex(wallet.Address),
	}
	score, err := s.scoreRepo.LatestByWalletID(ctx, cfg.WalletID)
	if err == nil && score != nil {
		leaderCtx.StrategyType = score.StrategyType
		leaderCtx.SmartScore = score.SmartScore
		leaderCtx.PoolTier = score.PoolTier
		leaderCtx.InfoEdgeLevel = score.InfoEdgeLevel
		leaderCtx.RiskLevel = score.RiskLevel
		leaderCtx.Momentum = score.Momentum
	}
	feat, err := s.featureRepo.LatestByWalletID(ctx, cfg.WalletID)
	if err == nil && feat != nil {
		leaderCtx.Pnl30D = feat.Pnl30d
	}

	// Exposure check
	currentExposure, _ := s.repo.SumExposureByConfig(ctx, configID)
	remainingBudget := cfg.MaxPositionUSDC - currentExposure
	if remainingBudget < 1 {
		return nil, ErrBudgetExhausted
	}

	outcome := "Yes"
	if leaderTrade.TokenSide == 1 {
		outcome = "No"
	}
	action := "Buy"
	if leaderTrade.TradeSide == 1 {
		action = "Sell"
	}

	tradeCtx := TradeContext{
		MarketTitle: leaderTrade.MarketTitle,
		MarketSlug:  leaderTrade.MarketSlug,
		Outcome:     outcome,
		Action:      action,
		Price:       leaderTrade.Price,
		Size:        leaderTrade.Size,
	}
	if market != nil {
		tradeCtx.MarketCategory = market.Category
		tradeCtx.MarketVolume = market.Volume
		tradeCtx.MarketLiquidity = market.Liquidity
	}

	// Copy history
	dashStats, _ := s.repo.GetDashboardStats(ctx, cfg.UserFingerprint)
	copyHistory := CopyHistoryContext{}
	if dashStats != nil {
		copyHistory.TotalCopies = dashStats.TotalCopies
		copyHistory.ProfitableCopies = dashStats.Profitable
		copyHistory.TotalPnL = dashStats.TotalPnL
	}

	signal := CopyTradeSignal{
		LeaderWallet: leaderCtx,
		NewTrade:     tradeCtx,
		UserSettings: UserSettingsContext{
			MaxPositionUSDC: cfg.MaxPositionUSDC,
			RiskPreference:  cfg.RiskPreference,
			CurrentExposure: currentExposure,
			RemainingBudget: remainingBudget,
		},
		CopyHistory: copyHistory,
	}

	agentDec, err := s.agent.Evaluate(ctx, signal)
	if err != nil {
		return nil, err
	}

	// Apply hard limits
	if agentDec.Decision == "copy" {
		maxSingle := cfg.MaxPositionUSDC * 0.20
		if agentDec.PositionUSDC > maxSingle {
			agentDec.PositionUSDC = maxSingle
		}
		if agentDec.PositionUSDC > remainingBudget {
			agentDec.PositionUSDC = remainingBudget
		}
		if agentDec.PositionUSDC < 1 {
			agentDec.Decision = "skip"
			agentDec.Reasoning = "计算仓位低于最小值 $1，跳过。"
			agentDec.ReasoningEn = "Calculated position below $1 minimum. Skipping."
		}
	}

	riskNotesJSON, _ := json.Marshal(agentDec.RiskNotes)
	var marketID *int64
	if market != nil {
		marketID = &market.ID
	}
	tradeID := leaderTrade.TradeID
	now := time.Now().UTC()

	decision := &model.CopyTradeDecision{
		ConfigID:      configID,
		LeaderTradeID: &tradeID,
		MarketID:      marketID,
		MarketTitle:   leaderTrade.MarketTitle,
		Decision:      agentDec.Decision,
		Confidence:    agentDec.Confidence,
		Outcome:       outcome,
		Action:        action,
		Price:         leaderTrade.Price,
		SizeUSDC:      agentDec.PositionUSDC,
		StopLossPrice: agentDec.StopLossPrice,
		Reasoning:     agentDec.Reasoning,
		ReasoningEn:   agentDec.ReasoningEn,
		RiskNotes:     riskNotesJSON,
		ModelID:       agentDec.ModelID,
		InputTokens:   agentDec.InputTokens,
		OutputTokens:  agentDec.OutputTokens,
		LatencyMS:     agentDec.LatencyMS,
		Status:        "pending",
	}

	if agentDec.Decision == "copy" {
		decision.Status = "executed"
		decision.ExecutedAt = &now
	}

	if err := s.repo.CreateDecision(ctx, decision); err != nil {
		return nil, err
	}

	// Update daily perf
	perfDate := now.Truncate(24 * time.Hour)
	perf := &model.CopyTradeDailyPerf{
		ConfigID: configID,
		PerfDate: perfDate,
	}
	if agentDec.Decision == "copy" {
		perf.TotalCopies = 1
		perf.TotalExposure = agentDec.PositionUSDC
	} else {
		perf.Skipped = 1
	}
	_ = s.repo.UpsertDailyPerf(ctx, perf)

	return decisionToView(decision), nil
}

// ── Query methods ──

func (s *Service) GetDashboard(ctx context.Context, userFP string) (*DashboardView, error) {
	userFP = strings.TrimSpace(userFP)
	if userFP == "" {
		return nil, errors.New("empty user fingerprint")
	}
	stats, err := s.repo.GetDashboardStats(ctx, userFP)
	if err != nil {
		return nil, err
	}
	configs, err := s.ListConfigs(ctx, userFP)
	if err != nil {
		return nil, err
	}
	recentDecs, err := s.repo.ListRecentDecisionsByUser(ctx, userFP, 10)
	if err != nil {
		return nil, err
	}
	decViews := make([]DecisionView, 0, len(recentDecs))
	for _, d := range recentDecs {
		decViews = append(decViews, *decisionToView(&d))
	}

	activeCount := 0
	for _, c := range configs {
		if c.Enabled {
			activeCount++
		}
	}

	var winRate float64
	if stats.TotalCopies > 0 {
		winRate = float64(stats.Profitable) / float64(stats.TotalCopies)
	}

	return &DashboardView{
		TotalPnL:        stats.TotalPnL,
		WinRate:         winRate,
		TotalCopies:     stats.TotalCopies,
		TotalSkipped:    stats.TotalSkipped,
		OpenPositions:   stats.OpenCount,
		ActiveConfigs:   activeCount,
		Configs:         configs,
		RecentDecisions: decViews,
	}, nil
}

func (s *Service) ListDecisions(ctx context.Context, userFP string, walletID int64, page, pageSize int) (*DecisionListResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	cfg, err := s.repo.GetConfig(ctx, userFP, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	rows, total, err := s.repo.ListDecisionsByConfig(ctx, cfg.ID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	items := make([]DecisionView, 0, len(rows))
	for _, r := range rows {
		items = append(items, *decisionToView(&r))
	}
	return &DecisionListResult{
		Items: items,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *Service) ListOpenPositions(ctx context.Context, userFP string) ([]DecisionView, error) {
	userFP = strings.TrimSpace(userFP)
	if userFP == "" {
		return nil, errors.New("empty user fingerprint")
	}
	rows, err := s.repo.ListOpenPositionsByUser(ctx, userFP)
	if err != nil {
		return nil, err
	}
	views := make([]DecisionView, 0, len(rows))
	for _, r := range rows {
		views = append(views, *decisionToView(&r))
	}
	return views, nil
}

func (s *Service) GetPerformance(ctx context.Context, userFP string, walletID int64) (*PerformanceView, error) {
	cfg, err := s.repo.GetConfig(ctx, userFP, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	perfRows, err := s.repo.ListDailyPerfByConfig(ctx, cfg.ID, 90)
	if err != nil {
		return nil, err
	}
	dailyPoints := make([]DailyPerfPoint, 0, len(perfRows))
	var totalPnL float64
	var totalCopies int64
	var profitable int64
	for _, p := range perfRows {
		dailyPoints = append(dailyPoints, DailyPerfPoint{
			Date:   p.PerfDate.Format("2006-01-02"),
			PnL:    p.TotalPnL,
			Copies: p.TotalCopies,
		})
		totalPnL += p.TotalPnL
		totalCopies += int64(p.TotalCopies)
		profitable += int64(p.Profitable)
	}
	var winRate float64
	if totalCopies > 0 {
		winRate = float64(profitable) / float64(totalCopies)
	}
	return &PerformanceView{
		TotalPnL:    totalPnL,
		WinRate:     winRate,
		TotalCopies: totalCopies,
		DailyPoints: dailyPoints,
	}, nil
}

func (s *Service) ClosePosition(ctx context.Context, userFP string, decisionID int64) (*DecisionView, error) {
	dec, err := s.repo.GetDecision(ctx, decisionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	cfg, err := s.repo.GetConfigByID(ctx, dec.ConfigID)
	if err != nil {
		return nil, err
	}
	if cfg.UserFingerprint != userFP {
		return nil, ErrNotFound
	}

	closePrice := dec.Price
	realizedPnL := 0.0
	if err := s.repo.CloseDecision(ctx, decisionID, closePrice, realizedPnL); err != nil {
		return nil, err
	}

	updated, _ := s.repo.GetDecision(ctx, decisionID)
	return decisionToView(updated), nil
}

// ── helpers ──

func decisionToView(d *model.CopyTradeDecision) *DecisionView {
	var riskNotes []string
	_ = json.Unmarshal(d.RiskNotes, &riskNotes)

	return &DecisionView{
		ID:            d.ID,
		Decision:      d.Decision,
		Confidence:    d.Confidence,
		MarketTitle:   d.MarketTitle,
		Outcome:       d.Outcome,
		Action:        d.Action,
		Price:         d.Price,
		SizeUSDC:      d.SizeUSDC,
		StopLossPrice: d.StopLossPrice,
		Reasoning:     d.Reasoning,
		ReasoningEn:   d.ReasoningEn,
		RiskNotes:     riskNotes,
		Status:        d.Status,
		RealizedPnL:   d.RealizedPnL,
		CreatedAt:     d.CreatedAt.UTC().Format(time.RFC3339),
	}
}
