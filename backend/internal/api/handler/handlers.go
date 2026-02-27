package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"easy-arbitra/backend/internal/copytrade"
	"easy-arbitra/backend/internal/service"
	"easy-arbitra/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	walletService    *service.WalletService
	marketService    *service.MarketService
	statsService     *service.StatsService
	anomalyService   *service.AnomalyService
	explainService   *service.ExplanationService
	infoEdge         *service.InfoEdgeService
	aiService        *service.AIService
	watchlistService *service.WatchlistService
	portfolioService *service.PortfolioService
	copyTradeService *copytrade.Service
	copyTradeRepo    *copytrade.Repository
	readyCheck       func(*gin.Context) error
}

func New(
	walletService *service.WalletService,
	marketService *service.MarketService,
	statsService *service.StatsService,
	anomalyService *service.AnomalyService,
	explainService *service.ExplanationService,
	infoEdge *service.InfoEdgeService,
	aiService *service.AIService,
	watchlistService *service.WatchlistService,
	portfolioService *service.PortfolioService,
	copyTradeService *copytrade.Service,
	copyTradeRepo *copytrade.Repository,
	readyCheck func(*gin.Context) error,
) *Handlers {
	return &Handlers{
		walletService: walletService, marketService: marketService, statsService: statsService, anomalyService: anomalyService,
		explainService: explainService, infoEdge: infoEdge, aiService: aiService, watchlistService: watchlistService,
		portfolioService: portfolioService, copyTradeService: copyTradeService, copyTradeRepo: copyTradeRepo, readyCheck: readyCheck,
	}
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handlers) Ready(c *gin.Context) {
	if h.readyCheck != nil {
		if err := h.readyCheck(c); err != nil {
			response.Internal(c, "not ready")
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func (h *Handlers) ListWallets(c *gin.Context) {
	page, pageSize := parsePaging(c)
	tracked := parseBoolPtr(c.Query("tracked"))

	rows, err := h.walletService.List(c.Request.Context(), service.WalletListQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   c.DefaultQuery("sort_by", "updated_at"),
		Order:    c.DefaultQuery("order", "desc"),
		Tracked:  tracked,
		Search:   c.Query("q"),
	})
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) GetWallet(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	row, err := h.walletService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handlers) ListPotentialWallets(c *gin.Context) {
	page, pageSize := parsePaging(c)
	minTrades := int64(parsePositiveInt(c.DefaultQuery("min_trades", "100"), 100))
	minRealizedPnL, err := strconv.ParseFloat(strings.TrimSpace(c.DefaultQuery("min_realized_pnl", "0")), 64)
	if err != nil {
		response.BadRequest(c, "invalid min_realized_pnl")
		return
	}
	var hasAIReport *bool
	if raw := strings.TrimSpace(c.Query("has_ai_report")); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			response.BadRequest(c, "invalid has_ai_report")
			return
		}
		hasAIReport = &v
	}

	rows, err := h.walletService.ListPotential(c.Request.Context(), service.PotentialWalletListQuery{
		Page:           page,
		PageSize:       pageSize,
		MinTrades:      minTrades,
		MinRealizedPnL: minRealizedPnL,
		StrategyType:   strings.TrimSpace(c.Query("strategy_type")),
		PoolTier:       strings.TrimSpace(c.Query("pool_tier")),
		HasAIReport:    hasAIReport,
		SortBy:         strings.TrimSpace(c.DefaultQuery("sort_by", "trade_count")),
		Order:          strings.TrimSpace(c.DefaultQuery("order", "desc")),
	})
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) GetWalletProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	row, err := h.walletService.GetProfile(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handlers) GetWalletShareCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	row, err := h.walletService.GetShareCard(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handlers) GetWalletExplanations(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	row, err := h.explainService.GetWalletExplanation(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handlers) GetWalletInfoEdge(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	report, err := h.infoEdge.Evaluate(c.Request.Context(), id)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, report)
}

func (h *Handlers) ListMarkets(c *gin.Context) {
	page, pageSize := parsePaging(c)
	status := parseInt16Ptr(c.Query("status"))

	rows, err := h.marketService.List(c.Request.Context(), service.MarketListQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   c.DefaultQuery("sort_by", "updated_at"),
		Order:    c.DefaultQuery("order", "desc"),
		Category: c.Query("category"),
		Status:   status,
	})
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) GetMarket(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid market id")
		return
	}
	row, err := h.marketService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "market not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handlers) GetOverviewStats(c *gin.Context) {
	stats, err := h.statsService.Overview(c.Request.Context())
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, stats)
}

func (h *Handlers) GetOpsHighlights(c *gin.Context) {
	limit := parsePositiveInt(c.DefaultQuery("limit", "5"), 5)
	if limit > 20 {
		limit = 20
	}
	rows, err := h.statsService.OpsHighlights(c.Request.Context(), limit)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) GetLeaderboard(c *gin.Context) {
	page, pageSize := parsePaging(c)
	rows, err := h.statsService.Leaderboard(c.Request.Context(), service.LeaderboardQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   c.DefaultQuery("sort_by", "smart_score"),
		Order:    c.DefaultQuery("order", "desc"),
	})
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) ListAnomalies(c *gin.Context) {
	page, pageSize := parsePaging(c)
	severity, err := service.ParseOptionalInt16(strings.TrimSpace(c.Query("severity")))
	if err != nil {
		response.BadRequest(c, "invalid severity")
		return
	}
	ack, err := service.ParseOptionalBool(strings.TrimSpace(c.Query("acknowledged")))
	if err != nil {
		response.BadRequest(c, "invalid acknowledged")
		return
	}

	rows, err := h.anomalyService.List(c.Request.Context(), service.AnomalyListQuery{
		Page:         page,
		PageSize:     pageSize,
		Severity:     severity,
		AlertType:    c.Query("type"),
		Acknowledged: ack,
	})
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) AcknowledgeAnomaly(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid anomaly id")
		return
	}
	if err := h.anomalyService.Acknowledge(c.Request.Context(), id); err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, gin.H{"acknowledged": true})
}

func (h *Handlers) GetAnomaly(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid anomaly id")
		return
	}
	row, err := h.anomalyService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "anomaly not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, row)
}

func (h *Handlers) TriggerAIAnalysis(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	force := false
	if raw := strings.TrimSpace(c.Query("force")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			response.BadRequest(c, "invalid force query")
			return
		}
		force = parsed
	}
	report, err := h.aiService.AnalyzeByWalletID(c.Request.Context(), walletID, force)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		if err == service.ErrInsufficientTrades || err == service.ErrNonPositivePnL {
			response.BadRequest(c, err.Error())
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.Created(c, report)
}

func (h *Handlers) GetAIReport(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	report, err := h.aiService.LatestByWalletID(c.Request.Context(), walletID)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "report not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, report)
}

func (h *Handlers) ListAIReports(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	limit := parsePositiveInt(c.DefaultQuery("limit", "10"), 10)
	rows, err := h.aiService.ListByWalletID(c.Request.Context(), walletID, limit)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) AddToWatchlist(c *gin.Context) {
	var req struct {
		WalletID int64 `json:"wallet_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if req.WalletID <= 0 {
		response.BadRequest(c, "wallet_id is required")
		return
	}
	if err := h.watchlistService.AddByWalletID(c.Request.Context(), req.WalletID, userIdentifier(c)); err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.Created(c, gin.H{"wallet_id": req.WalletID, "watching": true})
}

func (h *Handlers) AddToWatchlistBatch(c *gin.Context) {
	var req struct {
		WalletIDs []int64 `json:"wallet_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if len(req.WalletIDs) == 0 {
		response.BadRequest(c, "wallet_ids is required")
		return
	}
	created := make([]int64, 0, len(req.WalletIDs))
	for _, walletID := range req.WalletIDs {
		if walletID <= 0 {
			continue
		}
		if err := h.watchlistService.AddByWalletID(c.Request.Context(), walletID, userIdentifier(c)); err != nil && err != service.ErrNotFound {
			response.Internal(c, err.Error())
			return
		}
		created = append(created, walletID)
	}
	response.Created(c, gin.H{"wallet_ids": created, "watching": true})
}

func (h *Handlers) RemoveFromWatchlist(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)
	if err != nil || walletID <= 0 {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	if err := h.watchlistService.RemoveByWalletID(c.Request.Context(), walletID, userIdentifier(c)); err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, gin.H{"wallet_id": walletID, "watching": false})
}

func (h *Handlers) ListWatchlist(c *gin.Context) {
	page, pageSize := parsePaging(c)
	rows, err := h.watchlistService.List(c.Request.Context(), service.WatchlistListQuery{
		Page:            page,
		PageSize:        pageSize,
		UserFingerprint: userIdentifier(c),
	})
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) GetWatchlistFeed(c *gin.Context) {
	page, pageSize := parsePaging(c)
	rows, err := h.watchlistService.Feed(c.Request.Context(), service.WatchlistFeedQuery{
		Page:            page,
		PageSize:        pageSize,
		UserFingerprint: userIdentifier(c),
		EventType:       strings.TrimSpace(c.Query("type")),
	})
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) GetWatchlistSummary(c *gin.Context) {
	rows, err := h.watchlistService.Summary(c.Request.Context(), userIdentifier(c))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) GetWalletPnLHistory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	limit := parsePositiveInt(c.DefaultQuery("limit", "90"), 90)
	if limit > 365 {
		limit = 365
	}
	rows, err := h.walletService.GetPnLHistory(c.Request.Context(), id, limit)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) ListWalletTrades(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	page, pageSize := parsePaging(c)
	rows, err := h.walletService.ListTrades(c.Request.Context(), id, page, pageSize)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func (h *Handlers) ListWalletPositions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	rows, err := h.walletService.ListPositions(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}

func parsePaging(c *gin.Context) (int, int) {
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := parsePositiveInt(c.DefaultQuery("page_size", "20"), 20)
	if pageSize > 200 {
		pageSize = 200
	}
	return page, pageSize
}

func parsePositiveInt(input string, fallback int) int {
	v, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func parseBoolPtr(input string) *bool {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	v, err := strconv.ParseBool(input)
	if err != nil {
		return nil
	}
	return &v
}

func parseInt16Ptr(input string) *int16 {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	v, err := strconv.ParseInt(input, 10, 16)
	if err != nil {
		return nil
	}
	val := int16(v)
	return &val
}

func userIdentifier(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("uid:%d", userID.(int64))
	}
	return ""
}
