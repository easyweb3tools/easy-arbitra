package handler

import (
	"net/http"
	"strconv"
	"strings"

	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/internal/service"
	"easy-arbitra/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	walletService *service.WalletService
	marketService *service.MarketService
	statsService  *service.StatsService
	dailyPickRepo *repository.DailyPickRepository
	walletRepo    *repository.WalletRepository
	readyCheck    func(*gin.Context) error
}

func New(
	walletService *service.WalletService,
	marketService *service.MarketService,
	statsService *service.StatsService,
	dailyPickRepo *repository.DailyPickRepository,
	walletRepo *repository.WalletRepository,
	readyCheck func(*gin.Context) error,
) *Handlers {
	return &Handlers{
		walletService: walletService, marketService: marketService, statsService: statsService,
		dailyPickRepo: dailyPickRepo, walletRepo: walletRepo, readyCheck: readyCheck,
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

func (h *Handlers) GetDailyPick(c *gin.Context) {
	pick, err := h.dailyPickRepo.GetLatest(c.Request.Context())
	if err != nil {
		response.NotFound(c, "no daily pick available")
		return
	}
	wallet, _ := h.walletRepo.GetByID(c.Request.Context(), pick.WalletID)
	response.OK(c, gin.H{
		"pick":   pick,
		"wallet": wallet,
	})
}

func (h *Handlers) ListDailyPickHistory(c *gin.Context) {
	limit := parsePositiveInt(c.DefaultQuery("limit", "14"), 14)
	if limit > 90 {
		limit = 90
	}
	rows, err := h.dailyPickRepo.ListRecent(c.Request.Context(), limit)
	if err != nil {
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
