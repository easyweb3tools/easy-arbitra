package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/pkg/polyaddr"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")
var ErrInsufficientTrades = errors.New("wallet has fewer than 100 trades")
var ErrNonPositivePnL = errors.New("wallet pnl is not positive")

type WalletService struct {
	walletRepo   *repository.WalletRepository
	scoreRepo    *repository.ScoreRepository
	tradeRepo    *repository.TradeRepository
	featureRepo  *repository.FeatureRepository
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
	StrategyType   string
	PoolTier       string
	HasAIReport    *bool
	SortBy         string
	Order          string
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
	PoolTier       string     `json:"pool_tier"`
	HasAIReport    bool       `json:"has_ai_report"`
	NLSummary      string     `json:"nl_summary"`
	Summary        string     `json:"summary"`
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
	PoolTier      string  `json:"pool_tier"`
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

type WalletView struct {
	ID        int64   `json:"id"`
	Address   string  `json:"address"`
	Pseudonym *string `json:"pseudonym,omitempty"`
	Tracked   bool    `json:"tracked"`
}

func walletToView(w model.Wallet) WalletView {
	return WalletView{ID: w.ID, Address: polyaddr.BytesToHex(w.Address), Pseudonym: w.Pseudonym, Tracked: w.IsTracked}
}

type PnLHistoryPoint struct {
	Date          string  `json:"date"`
	Pnl7D         float64 `json:"pnl_7d"`
	Pnl30D        float64 `json:"pnl_30d"`
	Pnl90D        float64 `json:"pnl_90d"`
	TradeCount30D int     `json:"trade_count_30d"`
	ActiveDays30D int     `json:"active_days_30d"`
	AvgEdge       float64 `json:"avg_edge"`
}

type TradeHistoryView struct {
	ID          int64   `json:"id"`
	BlockTime   string  `json:"block_time"`
	MarketTitle string  `json:"market_title"`
	MarketSlug  string  `json:"market_slug"`
	Outcome     string  `json:"outcome"`
	Action      string  `json:"action"`
	Price       float64 `json:"price"`
	Size        float64 `json:"size"`
	FeePaid     float64 `json:"fee_paid"`
	IsMaker     bool    `json:"is_maker"`
}

type TradeHistoryResult struct {
	Items      []TradeHistoryView `json:"items"`
	Pagination Pagination         `json:"pagination"`
}

type WalletPositionView struct {
	MarketID    int64   `json:"market_id"`
	MarketTitle string  `json:"market_title"`
	MarketSlug  string  `json:"market_slug"`
	Category    string  `json:"category"`
	NetSize     float64 `json:"net_size"`
	AvgPrice    float64 `json:"avg_price"`
	TotalVolume float64 `json:"total_volume"`
	TradeCount  int64   `json:"trade_count"`
	LastTradeAt string  `json:"last_trade_at"`
}

// ── Constructors ──

func NewWalletService(
	walletRepo *repository.WalletRepository,
	scoreRepo *repository.ScoreRepository,
	tradeRepo *repository.TradeRepository,
	featureRepo *repository.FeatureRepository,
	aiReportRepo *repository.AIReportRepository,
	infoEdge *InfoEdgeService,
) *WalletService {
	return &WalletService{
		walletRepo: walletRepo, scoreRepo: scoreRepo, tradeRepo: tradeRepo,
		featureRepo: featureRepo, aiReportRepo: aiReportRepo, infoEdge: infoEdge,
	}
}

func NewMarketService(marketRepo *repository.MarketRepository) *MarketService {
	return &MarketService{marketRepo: marketRepo}
}

func NewStatsService(walletRepo *repository.WalletRepository, marketRepo *repository.MarketRepository, scoreRepo *repository.ScoreRepository) *StatsService {
	return &StatsService{walletRepo: walletRepo, marketRepo: marketRepo, scoreRepo: scoreRepo}
}

// ── WalletService ──

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
	total, err := s.walletRepo.CountPotentialWallets(ctx, repository.PotentialWalletFilter{
		MinTrades:      q.MinTrades,
		MinRealizedPnL: q.MinRealizedPnL,
		StrategyType:   strings.TrimSpace(q.StrategyType),
		PoolTier:       strings.TrimSpace(q.PoolTier),
		HasAIReport:    q.HasAIReport,
		SortBy:         strings.TrimSpace(q.SortBy),
		Order:          strings.TrimSpace(q.Order),
	})
	if err != nil {
		return nil, err
	}
	rows, err := s.walletRepo.ListPotentialWallets(ctx, repository.PotentialWalletFilter{
		MinTrades:      q.MinTrades,
		MinRealizedPnL: q.MinRealizedPnL,
		StrategyType:   strings.TrimSpace(q.StrategyType),
		PoolTier:       strings.TrimSpace(q.PoolTier),
		HasAIReport:    q.HasAIReport,
		SortBy:         strings.TrimSpace(q.SortBy),
		Order:          strings.TrimSpace(q.Order),
		Limit:          q.PageSize,
		Offset:         (q.Page - 1) * q.PageSize,
	})
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
			PoolTier:       row.PoolTier,
			HasAIReport:    row.LastAnalyzedAt != nil,
			NLSummary:      row.NLSummary,
			Summary:        buildFallbackSummary(row.StrategyType, row.SmartScore, row.NLSummary, row.TradeCount, row.RealizedPnL),
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
		poolTier := score.PoolTier
		if strings.TrimSpace(poolTier) == "" {
			poolTier = "observation"
		}
		profile.Strategy = &StrategySummary{
			StrategyType:  score.StrategyType,
			SmartScore:    score.SmartScore,
			InfoEdgeLevel: score.InfoEdgeLevel,
			PoolTier:      poolTier,
			Confidence:    score.StrategyConfidence,
			ScoredAt:      score.ScoredAt.UTC().Format(time.RFC3339),
		}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return profile, nil
}

func (s *WalletService) GetPnLHistory(ctx context.Context, walletID int64, limit int) ([]PnLHistoryPoint, error) {
	if _, err := s.walletRepo.GetByID(ctx, walletID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	rows, err := s.featureRepo.ListByWalletID(ctx, walletID, limit)
	if err != nil {
		return nil, err
	}
	points := make([]PnLHistoryPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, PnLHistoryPoint{
			Date:          row.FeatureDate.Format("2006-01-02"),
			Pnl7D:         row.Pnl7d,
			Pnl30D:        row.Pnl30d,
			Pnl90D:        row.Pnl90d,
			TradeCount30D: row.TradeCount30d,
			ActiveDays30D: row.ActiveDays30d,
			AvgEdge:       row.AvgEdge,
		})
	}
	return points, nil
}

func (s *WalletService) ListTrades(ctx context.Context, walletID int64, page int, pageSize int) (*TradeHistoryResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	if _, err := s.walletRepo.GetByID(ctx, walletID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	rows, total, err := s.tradeRepo.ListByWalletID(ctx, walletID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	items := make([]TradeHistoryView, 0, len(rows))
	for _, row := range rows {
		outcome := "Yes"
		if row.TokenSide == 1 {
			outcome = "No"
		}
		action := "Buy"
		if row.TradeSide == 1 {
			action = "Sell"
		}
		items = append(items, TradeHistoryView{
			ID:          row.TradeID,
			BlockTime:   row.BlockTime.UTC().Format(time.RFC3339),
			MarketTitle: row.MarketTitle,
			MarketSlug:  row.MarketSlug,
			Outcome:     outcome,
			Action:      action,
			Price:       row.Price,
			Size:        row.Size,
			FeePaid:     row.FeePaid,
			IsMaker:     row.IsMaker,
		})
	}
	return &TradeHistoryResult{
		Items: items,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}, nil
}

func (s *WalletService) ListPositions(ctx context.Context, walletID int64) ([]WalletPositionView, error) {
	if _, err := s.walletRepo.GetByID(ctx, walletID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	rows, err := s.tradeRepo.AggregatePositionsByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	views := make([]WalletPositionView, 0, len(rows))
	for _, row := range rows {
		views = append(views, WalletPositionView{
			MarketID:    row.MarketID,
			MarketTitle: row.MarketTitle,
			MarketSlug:  row.MarketSlug,
			Category:    row.Category,
			NetSize:     row.NetSize,
			AvgPrice:    row.AvgPrice,
			TotalVolume: row.TotalVolume,
			TradeCount:  row.TradeCount,
			LastTradeAt: row.LastTradeAt.UTC().Format(time.RFC3339),
		})
	}
	return views, nil
}

// ── MarketService ──

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

// ── StatsService ──

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

// ── Helpers ──

func buildFallbackSummary(strategyType string, smartScore int, nlSummary string, tradeCount int64, realizedPnL float64) string {
	if strings.TrimSpace(nlSummary) != "" {
		return nlSummary
	}
	return fmt.Sprintf(
		"%s trader (score %d) with %d trades and %.2f USDC realized PnL",
		strategyType, smartScore, tradeCount, realizedPnL,
	)
}
