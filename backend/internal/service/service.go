package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"easy-arbitra/backend/internal/ai"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/pkg/polyaddr"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")
var ErrInsufficientTrades = errors.New("wallet has fewer than 100 trades")
var ErrNonPositivePnL = errors.New("wallet pnl is not positive")

type WalletService struct {
	walletRepo   *repository.WalletRepository
	scoreRepo    *repository.ScoreRepository
	tradeRepo    *repository.TradeRepository
	aiReportRepo *repository.AIReportRepository
	infoEdge     *InfoEdgeService
}

type MarketService struct {
	marketRepo *repository.MarketRepository
}

type StatsService struct {
	walletRepo *repository.WalletRepository
	marketRepo *repository.MarketRepository
	scoreRepo  *repository.ScoreRepository
}

type OpsHighlights struct {
	AsOf                  string                         `json:"as_of"`
	NewPotentialWallets24 int64                          `json:"new_potential_wallets_24h"`
	TopRealizedPnL24h     []OpsTopRealizedWalletView     `json:"top_realized_pnl_24h"`
	TopAIConfidence       []OpsTopAIConfidenceWalletView `json:"top_ai_confidence"`
}

type OpsTopRealizedWalletView struct {
	Wallet         WalletView `json:"wallet"`
	TradeCount     int64      `json:"trade_count"`
	RealizedPnL    float64    `json:"realized_pnl"`
	RealizedPnL24h float64    `json:"realized_pnl_24h"`
	HasAIReport    bool       `json:"has_ai_report"`
	NLSummary      string     `json:"nl_summary"`
	ModelID        string     `json:"model_id"`
	LastAnalyzedAt *string    `json:"last_analyzed_at,omitempty"`
}

type OpsTopAIConfidenceWalletView struct {
	Wallet         WalletView `json:"wallet"`
	TradeCount     int64      `json:"trade_count"`
	RealizedPnL    float64    `json:"realized_pnl"`
	SmartScore     int        `json:"smart_score"`
	InfoEdgeLevel  string     `json:"info_edge_level"`
	StrategyType   string     `json:"strategy_type"`
	NLSummary      string     `json:"nl_summary"`
	LastAnalyzedAt *string    `json:"last_analyzed_at,omitempty"`
}

type AIService struct {
	walletRepo    *repository.WalletRepository
	scoreRepo     *repository.ScoreRepository
	tradeRepo     *repository.TradeRepository
	aiReportRepo  *repository.AIReportRepository
	watchlistRepo *repository.WatchlistRepository
	analyzer      ai.Analyzer
	cacheTTL      time.Duration
}

type WatchlistService struct {
	walletRepo    *repository.WalletRepository
	watchlistRepo *repository.WatchlistRepository
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type WalletListQuery struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
	Tracked  *bool
	Search   string
}

type PotentialWalletListQuery struct {
	Page           int
	PageSize       int
	MinTrades      int64
	MinRealizedPnL float64
}

type MarketListQuery struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
	Category string
	Status   *int16
}

type WalletListResult struct {
	Items      []WalletView `json:"items"`
	Pagination Pagination   `json:"pagination"`
}

type PotentialWalletView struct {
	Wallet         WalletView `json:"wallet"`
	TotalTrades    int64      `json:"total_trades"`
	TradingPnL     float64    `json:"trading_pnl"`
	MakerRebates   float64    `json:"maker_rebates"`
	RealizedPnL    float64    `json:"realized_pnl"`
	SmartScore     int        `json:"smart_score"`
	InfoEdgeLevel  string     `json:"info_edge_level"`
	StrategyType   string     `json:"strategy_type"`
	HasAIReport    bool       `json:"has_ai_report"`
	LastAnalyzedAt *string    `json:"last_analyzed_at,omitempty"`
}

type PotentialWalletListResult struct {
	Items      []PotentialWalletView `json:"items"`
	Pagination Pagination            `json:"pagination"`
}

type MarketListResult struct {
	Items      []model.Market `json:"items"`
	Pagination Pagination     `json:"pagination"`
}

type WalletProfile struct {
	Wallet   WalletView          `json:"wallet"`
	Layer1   Layer1Facts         `json:"layer1_facts"`
	Strategy *StrategySummary    `json:"strategy,omitempty"`
	Layer3   Layer3InfoEdge      `json:"layer3_info_edge"`
	Meta     map[string][]string `json:"meta"`
}

type WalletShareCardView struct {
	Wallet        WalletView `json:"wallet"`
	TotalTrades   int64      `json:"total_trades"`
	RealizedPnL   float64    `json:"realized_pnl"`
	SmartScore    int        `json:"smart_score"`
	InfoEdgeLevel string     `json:"info_edge_level"`
	StrategyType  string     `json:"strategy_type"`
	HasAIReport   bool       `json:"has_ai_report"`
	NLSummary     string     `json:"nl_summary"`
	UpdatedAt     string     `json:"updated_at"`
}

type Layer1Facts struct {
	RealizedPnL   float64 `json:"realized_pnl"`
	UnrealizedPnL float64 `json:"unrealized_pnl"`
	TradingPnL    float64 `json:"trading_pnl"`
	MakerRebates  float64 `json:"maker_rebates"`
	FeesPaid      float64 `json:"fees_paid"`
	TotalTrades   int64   `json:"total_trades"`
	Volume30D     float64 `json:"volume_30d"`
}

type StrategySummary struct {
	StrategyType  string  `json:"strategy_type"`
	SmartScore    int     `json:"smart_score"`
	InfoEdgeLevel string  `json:"info_edge_level"`
	Confidence    float64 `json:"confidence"`
	ScoredAt      string  `json:"scored_at"`
}

type Layer3InfoEdge struct {
	MeanDeltaMinutes float64 `json:"mean_delta_minutes"`
	StdDevMinutes    float64 `json:"stddev_minutes"`
	Samples          int64   `json:"samples"`
	PValue           float64 `json:"p_value"`
	Label            string  `json:"label"`
}

type OverviewStats struct {
	TrackedWallets int64 `json:"tracked_wallets"`
	IndexedMarkets int64 `json:"indexed_markets"`
}

type LeaderboardQuery struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

type LeaderboardItem struct {
	WalletID      int64   `json:"wallet_id"`
	Address       string  `json:"address"`
	Pseudonym     *string `json:"pseudonym,omitempty"`
	StrategyType  string  `json:"strategy_type"`
	SmartScore    int     `json:"smart_score"`
	InfoEdgeLevel string  `json:"info_edge_level"`
	ScoredAt      string  `json:"scored_at"`
}

type LeaderboardResult struct {
	Items      []LeaderboardItem `json:"items"`
	Pagination Pagination        `json:"pagination"`
}

type WatchlistListQuery struct {
	Page            int
	PageSize        int
	UserFingerprint string
}

type WatchlistFeedQuery struct {
	Page            int
	PageSize        int
	UserFingerprint string
	EventType       string
}

type WatchlistItem struct {
	WatchlistID    int64      `json:"watchlist_id"`
	WatchlistedAt  string     `json:"watchlisted_at"`
	Wallet         WalletView `json:"wallet"`
	TotalTrades    int64      `json:"total_trades"`
	TradingPnL     float64    `json:"trading_pnl"`
	MakerRebates   float64    `json:"maker_rebates"`
	RealizedPnL    float64    `json:"realized_pnl"`
	SmartScore     int        `json:"smart_score"`
	InfoEdgeLevel  string     `json:"info_edge_level"`
	StrategyType   string     `json:"strategy_type"`
	HasAIReport    bool       `json:"has_ai_report"`
	LastAnalyzedAt *string    `json:"last_analyzed_at,omitempty"`
}

type WatchlistListResult struct {
	Items      []WatchlistItem `json:"items"`
	Pagination Pagination      `json:"pagination"`
}

type WatchlistFeedItem struct {
	EventID      int64          `json:"event_id"`
	Wallet       WalletView     `json:"wallet"`
	EventType    string         `json:"event_type"`
	EventPayload map[string]any `json:"event_payload"`
	EventTime    string         `json:"event_time"`
}

type WatchlistFeedResult struct {
	Items      []WatchlistFeedItem `json:"items"`
	Pagination Pagination          `json:"pagination"`
}

func NewWalletService(
	walletRepo *repository.WalletRepository,
	scoreRepo *repository.ScoreRepository,
	tradeRepo *repository.TradeRepository,
	aiReportRepo *repository.AIReportRepository,
	infoEdge *InfoEdgeService,
) *WalletService {
	return &WalletService{walletRepo: walletRepo, scoreRepo: scoreRepo, tradeRepo: tradeRepo, aiReportRepo: aiReportRepo, infoEdge: infoEdge}
}

func NewMarketService(marketRepo *repository.MarketRepository) *MarketService {
	return &MarketService{marketRepo: marketRepo}
}

func NewStatsService(walletRepo *repository.WalletRepository, marketRepo *repository.MarketRepository, scoreRepo *repository.ScoreRepository) *StatsService {
	return &StatsService{walletRepo: walletRepo, marketRepo: marketRepo, scoreRepo: scoreRepo}
}

func NewAIService(
	walletRepo *repository.WalletRepository,
	scoreRepo *repository.ScoreRepository,
	tradeRepo *repository.TradeRepository,
	aiReportRepo *repository.AIReportRepository,
	watchlistRepo *repository.WatchlistRepository,
	analyzer ai.Analyzer,
	cacheTTL time.Duration,
) *AIService {
	return &AIService{walletRepo: walletRepo, scoreRepo: scoreRepo, tradeRepo: tradeRepo, aiReportRepo: aiReportRepo, watchlistRepo: watchlistRepo, analyzer: analyzer, cacheTTL: cacheTTL}
}

func NewWatchlistService(walletRepo *repository.WalletRepository, watchlistRepo *repository.WatchlistRepository) *WatchlistService {
	return &WatchlistService{walletRepo: walletRepo, watchlistRepo: watchlistRepo}
}

func (s *WalletService) List(ctx context.Context, q WalletListQuery) (*WalletListResult, error) {
	rows, total, err := s.walletRepo.List(ctx, repository.WalletListFilter{
		Tracked: q.Tracked,
		Search:  strings.TrimSpace(q.Search),
		SortBy:  q.SortBy,
		Order:   q.Order,
		Limit:   q.PageSize,
		Offset:  (q.Page - 1) * q.PageSize,
	})
	if err != nil {
		return nil, err
	}
	items := make([]WalletView, 0, len(rows))
	for _, row := range rows {
		items = append(items, walletToView(row))
	}
	return &WalletListResult{
		Items: items,
		Pagination: Pagination{
			Page:     q.Page,
			PageSize: q.PageSize,
			Total:    total,
		},
	}, nil
}

func (s *WalletService) GetByID(ctx context.Context, id int64) (*WalletView, error) {
	row, err := s.walletRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	view := walletToView(*row)
	return &view, nil
}

func (s *WalletService) ListPotential(ctx context.Context, q PotentialWalletListQuery) (*PotentialWalletListResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 200 {
		q.PageSize = 200
	}
	if q.MinTrades <= 0 {
		q.MinTrades = 100
	}
	total, err := s.walletRepo.CountPotentialWallets(ctx, q.MinTrades, q.MinRealizedPnL)
	if err != nil {
		return nil, err
	}
	rows, err := s.walletRepo.ListPotentialWallets(ctx, q.MinTrades, q.MinRealizedPnL, q.PageSize, (q.Page-1)*q.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]PotentialWalletView, 0, len(rows))
	for _, row := range rows {
		var analyzedAt *string
		if row.LastAnalyzedAt != nil {
			v := row.LastAnalyzedAt.UTC().Format(time.RFC3339)
			analyzedAt = &v
		}
		items = append(items, PotentialWalletView{
			Wallet: WalletView{
				ID:        row.WalletID,
				Address:   polyaddr.BytesToHex(row.Address),
				Pseudonym: row.Pseudonym,
				Tracked:   row.IsTracked,
			},
			TotalTrades:    row.TradeCount,
			TradingPnL:     row.TradingPnL,
			MakerRebates:   row.MakerRebates,
			RealizedPnL:    row.RealizedPnL,
			SmartScore:     row.SmartScore,
			InfoEdgeLevel:  row.InfoEdgeLevel,
			StrategyType:   row.StrategyType,
			HasAIReport:    row.LastAnalyzedAt != nil,
			LastAnalyzedAt: analyzedAt,
		})
	}
	return &PotentialWalletListResult{
		Items: items,
		Pagination: Pagination{
			Page:     q.Page,
			PageSize: q.PageSize,
			Total:    total,
		},
	}, nil
}

func (s *WalletService) GetProfile(ctx context.Context, id int64) (*WalletProfile, error) {
	wallet, err := s.walletRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	pnl, err := s.tradeRepo.AggregateByWalletID(ctx, id)
	if err != nil {
		return nil, err
	}

	profile := &WalletProfile{
		Wallet: walletToView(*wallet),
		Layer1: Layer1Facts{
			TradingPnL:    pnl.TradingPnL,
			MakerRebates:  pnl.MakerRebates,
			FeesPaid:      pnl.FeesPaid,
			TotalTrades:   pnl.TotalTrades,
			Volume30D:     pnl.Volume30D,
			RealizedPnL:   pnl.TradingPnL + pnl.MakerRebates,
			UnrealizedPnL: 0,
		},
		Meta: map[string][]string{
			"disclosures": {
				"Scores are probabilistic estimates based on publicly available data.",
				"Classification does not constitute evidence of wrongdoing.",
			},
		},
	}
	info, err := s.infoEdge.Evaluate(ctx, id)
	if err == nil {
		profile.Layer3 = Layer3InfoEdge{
			MeanDeltaMinutes: info.MeanDeltaMinutes,
			StdDevMinutes:    info.StdDevMinutes,
			Samples:          info.Samples,
			PValue:           info.PValue,
			Label:            info.Classification,
		}
	}

	score, err := s.scoreRepo.LatestByWalletID(ctx, id)
	if err == nil {
		profile.Strategy = &StrategySummary{
			StrategyType:  score.StrategyType,
			SmartScore:    score.SmartScore,
			InfoEdgeLevel: score.InfoEdgeLevel,
			Confidence:    score.StrategyConfidence,
			ScoredAt:      score.ScoredAt.UTC().Format(time.RFC3339),
		}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return profile, nil
}

func (s *WalletService) GetShareCard(ctx context.Context, id int64) (*WalletShareCardView, error) {
	wallet, err := s.walletRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	pnl, err := s.tradeRepo.AggregateByWalletID(ctx, id)
	if err != nil {
		return nil, err
	}
	realized := pnl.TradingPnL + pnl.MakerRebates

	strategyType := "unknown"
	smartScore := 0
	infoEdgeLevel := "unknown"
	updatedAt := wallet.UpdatedAt.UTC()
	score, err := s.scoreRepo.LatestByWalletID(ctx, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if score != nil {
		strategyType = score.StrategyType
		smartScore = score.SmartScore
		infoEdgeLevel = score.InfoEdgeLevel
		if score.ScoredAt.After(updatedAt) {
			updatedAt = score.ScoredAt.UTC()
		}
	}

	hasAI := false
	summary := ""
	report, err := s.aiReportRepo.LatestByWalletID(ctx, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if report != nil {
		hasAI = true
		summary = report.NLSummary
		if report.CreatedAt.After(updatedAt) {
			updatedAt = report.CreatedAt.UTC()
		}
	}

	return &WalletShareCardView{
		Wallet: WalletView{
			ID:        wallet.ID,
			Address:   polyaddr.BytesToHex(wallet.Address),
			Pseudonym: wallet.Pseudonym,
			Tracked:   wallet.IsTracked,
		},
		TotalTrades:   pnl.TotalTrades,
		RealizedPnL:   realized,
		SmartScore:    smartScore,
		InfoEdgeLevel: infoEdgeLevel,
		StrategyType:  strategyType,
		HasAIReport:   hasAI,
		NLSummary:     summary,
		UpdatedAt:     updatedAt.Format(time.RFC3339),
	}, nil
}

func (s *MarketService) List(ctx context.Context, q MarketListQuery) (*MarketListResult, error) {
	items, total, err := s.marketRepo.List(ctx, repository.MarketListFilter{
		Category: strings.TrimSpace(q.Category),
		Status:   q.Status,
		SortBy:   q.SortBy,
		Order:    q.Order,
		Limit:    q.PageSize,
		Offset:   (q.Page - 1) * q.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &MarketListResult{
		Items: items,
		Pagination: Pagination{
			Page:     q.Page,
			PageSize: q.PageSize,
			Total:    total,
		},
	}, nil
}

func (s *MarketService) GetByID(ctx context.Context, id int64) (*model.Market, error) {
	m, err := s.marketRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return m, nil
}

func (s *StatsService) Overview(ctx context.Context) (*OverviewStats, error) {
	tracked, err := s.walletRepo.CountTracked(ctx)
	if err != nil {
		return nil, err
	}
	markets, err := s.marketRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	return &OverviewStats{TrackedWallets: tracked, IndexedMarkets: markets}, nil
}

func (s *StatsService) Leaderboard(ctx context.Context, q LeaderboardQuery) (*LeaderboardResult, error) {
	rows, total, err := s.scoreRepo.Leaderboard(ctx, q.PageSize, (q.Page-1)*q.PageSize, q.SortBy, q.Order)
	if err != nil {
		return nil, err
	}
	items := make([]LeaderboardItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, LeaderboardItem{
			WalletID:      row.WalletID,
			Address:       polyaddr.BytesToHex(row.Address),
			Pseudonym:     row.Pseudonym,
			StrategyType:  row.StrategyType,
			SmartScore:    row.SmartScore,
			InfoEdgeLevel: row.InfoEdgeLevel,
			ScoredAt:      row.ScoredAt,
		})
	}
	return &LeaderboardResult{
		Items: items,
		Pagination: Pagination{
			Page:     q.Page,
			PageSize: q.PageSize,
			Total:    total,
		},
	}, nil
}

func (s *StatsService) OpsHighlights(ctx context.Context, limit int) (*OpsHighlights, error) {
	if limit <= 0 {
		limit = 5
	}
	const minTrades int64 = 100
	const minRealizedPnL float64 = 0

	newPotential, err := s.walletRepo.CountNewPotentialWallets24h(ctx, minTrades, minRealizedPnL)
	if err != nil {
		return nil, err
	}
	topRealized, err := s.walletRepo.ListOpsTopRealizedWallets(ctx, limit, minTrades, minRealizedPnL)
	if err != nil {
		return nil, err
	}
	topAI, err := s.walletRepo.ListOpsTopAIConfidenceWallets(ctx, limit, minTrades, minRealizedPnL)
	if err != nil {
		return nil, err
	}

	realizedViews := make([]OpsTopRealizedWalletView, 0, len(topRealized))
	for _, row := range topRealized {
		var analyzedAt *string
		if row.LastAnalyzedRaw != nil {
			v := row.LastAnalyzedRaw.UTC().Format(time.RFC3339)
			analyzedAt = &v
		}
		realizedViews = append(realizedViews, OpsTopRealizedWalletView{
			Wallet: WalletView{
				ID:      row.WalletID,
				Address: polyaddr.BytesToHex(row.Address),
			},
			TradeCount:     row.TradeCount,
			RealizedPnL:    row.RealizedPnL,
			RealizedPnL24h: row.RealizedPnL24h,
			HasAIReport:    row.HasAIReport,
			NLSummary:      row.LatestSummary,
			ModelID:        row.LatestModelID,
			LastAnalyzedAt: analyzedAt,
		})
	}

	aiViews := make([]OpsTopAIConfidenceWalletView, 0, len(topAI))
	for _, row := range topAI {
		var analyzedAt *string
		if row.LastAnalyzedAt != nil {
			v := row.LastAnalyzedAt.UTC().Format(time.RFC3339)
			analyzedAt = &v
		}
		aiViews = append(aiViews, OpsTopAIConfidenceWalletView{
			Wallet: WalletView{
				ID:      row.WalletID,
				Address: polyaddr.BytesToHex(row.Address),
			},
			TradeCount:     row.TradeCount,
			RealizedPnL:    row.RealizedPnL,
			SmartScore:     row.SmartScore,
			InfoEdgeLevel:  row.InfoEdgeLevel,
			StrategyType:   row.StrategyType,
			NLSummary:      row.NLSummary,
			LastAnalyzedAt: analyzedAt,
		})
	}

	return &OpsHighlights{
		AsOf:                  time.Now().UTC().Format(time.RFC3339),
		NewPotentialWallets24: newPotential,
		TopRealizedPnL24h:     realizedViews,
		TopAIConfidence:       aiViews,
	}, nil
}

func (s *AIService) AnalyzeByWalletID(ctx context.Context, walletID int64, force bool) (*model.AIAnalysisReport, error) {
	if !force && s.cacheTTL > 0 {
		latest, err := s.aiReportRepo.LatestByWalletID(ctx, walletID)
		if err == nil && time.Since(latest.CreatedAt) <= s.cacheTTL {
			return latest, nil
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	score, err := s.scoreRepo.LatestByWalletID(ctx, walletID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	in := ai.WalletAnalysisInput{
		WalletID:      walletID,
		WalletAddress: polyaddr.BytesToHex(wallet.Address),
		AsOf:          time.Now().UTC(),
		StrategyType:  "unknown",
		SmartScore:    0,
		InfoEdgeLevel: "luck",
	}
	var pnl *repository.WalletPnLSummary
	if s.tradeRepo != nil {
		pnl, err = s.tradeRepo.AggregateByWalletID(ctx, walletID)
		if err != nil {
			return nil, err
		}
		in.TradingPnL = pnl.TradingPnL
		in.MakerRebates = pnl.MakerRebates
		in.FeesPaid = pnl.FeesPaid
		in.TotalTrades = pnl.TotalTrades
		in.Volume30D = pnl.Volume30D
	}
	if pnl == nil || pnl.TotalTrades < 100 {
		return nil, ErrInsufficientTrades
	}
	realizedPnL := pnl.TradingPnL + pnl.MakerRebates
	if realizedPnL <= 0 {
		return nil, ErrNonPositivePnL
	}
	if score != nil {
		in.StrategyType = score.StrategyType
		in.SmartScore = score.SmartScore
		in.InfoEdgeLevel = score.InfoEdgeLevel
	}

	out, err := s.analyzer.AnalyzeWallet(ctx, in)
	if err != nil {
		return nil, err
	}

	warnings, _ := json.Marshal(out.RiskWarnings)
	report := &model.AIAnalysisReport{
		WalletID:     walletID,
		ModelID:      out.ModelID,
		Report:       datatypes.JSON(out.ReportJSON),
		NLSummary:    out.NLSummary,
		RiskWarnings: datatypes.JSON(warnings),
		InputTokens:  out.InputTokens,
		OutputTokens: out.OutputTokens,
		LatencyMS:    out.LatencyMS,
	}
	if err := s.aiReportRepo.Create(ctx, report); err != nil {
		return nil, err
	}
	if s.watchlistRepo != nil {
		evt, _ := json.Marshal(map[string]any{
			"model_id":      report.ModelID,
			"summary":       report.NLSummary,
			"input_tokens":  report.InputTokens,
			"output_tokens": report.OutputTokens,
		})
		_ = s.watchlistRepo.CreateUpdateEvent(ctx, walletID, "ai_report", evt)
	}
	return report, nil
}

func (s *AIService) LatestByWalletID(ctx context.Context, walletID int64) (*model.AIAnalysisReport, error) {
	report, err := s.aiReportRepo.LatestByWalletID(ctx, walletID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return report, nil
}

func (s *AIService) ListByWalletID(ctx context.Context, walletID int64, limit int) ([]model.AIAnalysisReport, error) {
	rows, err := s.aiReportRepo.ListByWalletID(ctx, walletID, limit)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

type WalletView struct {
	ID        int64   `json:"id"`
	Address   string  `json:"address"`
	Pseudonym *string `json:"pseudonym,omitempty"`
	Tracked   bool    `json:"tracked"`
}

func walletToView(w model.Wallet) WalletView {
	return WalletView{ID: w.ID, Address: polyaddr.BytesToHex(w.Address), Pseudonym: w.Pseudonym, Tracked: w.IsTracked}
}

func (s *WatchlistService) AddByWalletID(ctx context.Context, walletID int64, userFingerprint string) error {
	userFingerprint = strings.TrimSpace(userFingerprint)
	if userFingerprint == "" {
		return errors.New("empty user fingerprint")
	}
	if _, err := s.walletRepo.GetByID(ctx, walletID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.watchlistRepo.Add(ctx, walletID, userFingerprint)
}

func (s *WatchlistService) RemoveByWalletID(ctx context.Context, walletID int64, userFingerprint string) error {
	userFingerprint = strings.TrimSpace(userFingerprint)
	if userFingerprint == "" {
		return errors.New("empty user fingerprint")
	}
	return s.watchlistRepo.Remove(ctx, walletID, userFingerprint)
}

func (s *WatchlistService) List(ctx context.Context, q WatchlistListQuery) (*WatchlistListResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 200 {
		q.PageSize = 200
	}
	if strings.TrimSpace(q.UserFingerprint) == "" {
		return nil, errors.New("empty user fingerprint")
	}
	total, err := s.watchlistRepo.CountByUser(ctx, q.UserFingerprint)
	if err != nil {
		return nil, err
	}
	rows, err := s.watchlistRepo.ListByUser(ctx, q.UserFingerprint, q.PageSize, (q.Page-1)*q.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]WatchlistItem, 0, len(rows))
	for _, row := range rows {
		var analyzedAt *string
		if row.LastAnalyzedAt != nil {
			v := row.LastAnalyzedAt.UTC().Format(time.RFC3339)
			analyzedAt = &v
		}
		items = append(items, WatchlistItem{
			WatchlistID:   row.WatchlistID,
			WatchlistedAt: row.WatchlistedAt.UTC().Format(time.RFC3339),
			Wallet: WalletView{
				ID:        row.WalletID,
				Address:   polyaddr.BytesToHex(row.Address),
				Pseudonym: row.Pseudonym,
				Tracked:   row.IsTracked,
			},
			TotalTrades:    row.TradeCount,
			TradingPnL:     row.TradingPnL,
			MakerRebates:   row.MakerRebates,
			RealizedPnL:    row.RealizedPnL,
			SmartScore:     row.SmartScore,
			InfoEdgeLevel:  row.InfoEdgeLevel,
			StrategyType:   row.StrategyType,
			HasAIReport:    row.LastAnalyzedAt != nil,
			LastAnalyzedAt: analyzedAt,
		})
	}
	return &WatchlistListResult{
		Items: items,
		Pagination: Pagination{
			Page:     q.Page,
			PageSize: q.PageSize,
			Total:    total,
		},
	}, nil
}

func (s *WatchlistService) Feed(ctx context.Context, q WatchlistFeedQuery) (*WatchlistFeedResult, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 200 {
		q.PageSize = 200
	}
	if strings.TrimSpace(q.UserFingerprint) == "" {
		return nil, errors.New("empty user fingerprint")
	}
	total, err := s.watchlistRepo.CountFeedByUser(ctx, q.UserFingerprint, q.EventType)
	if err != nil {
		return nil, err
	}
	rows, err := s.watchlistRepo.ListFeedByUser(ctx, q.UserFingerprint, q.EventType, q.PageSize, (q.Page-1)*q.PageSize)
	if err != nil {
		return nil, err
	}
	items := make([]WatchlistFeedItem, 0, len(rows))
	for _, row := range rows {
		payload := map[string]any{}
		_ = json.Unmarshal(row.EventPayload, &payload)
		items = append(items, WatchlistFeedItem{
			EventID: row.EventID,
			Wallet: WalletView{
				ID:        row.WalletID,
				Address:   polyaddr.BytesToHex(row.Address),
				Pseudonym: row.Pseudonym,
			},
			EventType:    row.EventType,
			EventPayload: payload,
			EventTime:    row.EventTime.UTC().Format(time.RFC3339),
		})
	}
	return &WatchlistFeedResult{
		Items: items,
		Pagination: Pagination{
			Page:     q.Page,
			PageSize: q.PageSize,
			Total:    total,
		},
	}, nil
}
