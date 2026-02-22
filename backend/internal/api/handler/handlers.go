package handler

import (
	"net/http"
	"strconv"
	"strings"

	"easy-arbitra/backend/internal/service"
	"easy-arbitra/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	walletService  *service.WalletService
	marketService  *service.MarketService
	statsService   *service.StatsService
	anomalyService *service.AnomalyService
	explainService *service.ExplanationService
	infoEdge       *service.InfoEdgeService
	aiService      *service.AIService
	readyCheck     func(*gin.Context) error
}

func New(
	walletService *service.WalletService,
	marketService *service.MarketService,
	statsService *service.StatsService,
	anomalyService *service.AnomalyService,
	explainService *service.ExplanationService,
	infoEdge *service.InfoEdgeService,
	aiService *service.AIService,
	readyCheck func(*gin.Context) error,
) *Handlers {
	return &Handlers{
		walletService: walletService, marketService: marketService, statsService: statsService, anomalyService: anomalyService,
		explainService: explainService, infoEdge: infoEdge, aiService: aiService, readyCheck: readyCheck,
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

	rows, err := h.walletService.ListPotential(c.Request.Context(), service.PotentialWalletListQuery{
		Page:           page,
		PageSize:       pageSize,
		MinTrades:      minTrades,
		MinRealizedPnL: minRealizedPnL,
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
