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

type WalletService struct {
	walletRepo *repository.WalletRepository
	scoreRepo  *repository.ScoreRepository
	tradeRepo  *repository.TradeRepository
	infoEdge   *InfoEdgeService
}

type MarketService struct {
	marketRepo *repository.MarketRepository
}

type StatsService struct {
	walletRepo *repository.WalletRepository
	marketRepo *repository.MarketRepository
	scoreRepo  *repository.ScoreRepository
}

type AIService struct {
	walletRepo   *repository.WalletRepository
	scoreRepo    *repository.ScoreRepository
	tradeRepo    *repository.TradeRepository
	aiReportRepo *repository.AIReportRepository
	analyzer     ai.Analyzer
	cacheTTL     time.Duration
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

func NewWalletService(
	walletRepo *repository.WalletRepository,
	scoreRepo *repository.ScoreRepository,
	tradeRepo *repository.TradeRepository,
	infoEdge *InfoEdgeService,
) *WalletService {
	return &WalletService{walletRepo: walletRepo, scoreRepo: scoreRepo, tradeRepo: tradeRepo, infoEdge: infoEdge}
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
	analyzer ai.Analyzer,
	cacheTTL time.Duration,
) *AIService {
	return &AIService{walletRepo: walletRepo, scoreRepo: scoreRepo, tradeRepo: tradeRepo, aiReportRepo: aiReportRepo, analyzer: analyzer, cacheTTL: cacheTTL}
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

func (s *AIService) AnalyzeByWalletID(ctx context.Context, walletID int64) (*model.AIAnalysisReport, error) {
	if s.cacheTTL > 0 {
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
	if s.tradeRepo != nil {
		pnl, err := s.tradeRepo.AggregateByWalletID(ctx, walletID)
		if err == nil {
			in.TradingPnL = pnl.TradingPnL
			in.MakerRebates = pnl.MakerRebates
			in.FeesPaid = pnl.FeesPaid
			in.TotalTrades = pnl.TotalTrades
			in.Volume30D = pnl.Volume30D
		}
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
