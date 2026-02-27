package handler

import (
	"errors"
	"strconv"

	"easy-arbitra/backend/internal/copytrade"
	"easy-arbitra/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// ── Copy-Trading Handlers ──

func (h *Handlers) EnableCopyTrading(c *gin.Context) {
	var req struct {
		WalletID        int64   `json:"wallet_id"`
		MaxPositionUSDC float64 `json:"max_position_usdc"`
		RiskPreference  string  `json:"risk_preference"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if req.WalletID <= 0 {
		response.BadRequest(c, "wallet_id is required")
		return
	}
	cfg, err := h.copyTradeService.EnableCopyTrading(c.Request.Context(), userIdentifier(c), req.WalletID, req.MaxPositionUSDC, req.RiskPreference)
	if err != nil {
		if errors.Is(err, copytrade.ErrNotFound) {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.Created(c, cfg)
}

func (h *Handlers) DisableCopyTrading(c *gin.Context) {
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
	if err := h.copyTradeService.DisableCopyTrading(c.Request.Context(), userIdentifier(c), req.WalletID); err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, gin.H{"disabled": true})
}

func (h *Handlers) UpdateCopyTradeSettings(c *gin.Context) {
	var req struct {
		WalletID        int64   `json:"wallet_id"`
		MaxPositionUSDC float64 `json:"max_position_usdc"`
		RiskPreference  string  `json:"risk_preference"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid body")
		return
	}
	if req.WalletID <= 0 {
		response.BadRequest(c, "wallet_id is required")
		return
	}
	cfg, err := h.copyTradeService.UpdateSettings(c.Request.Context(), userIdentifier(c), req.WalletID, req.MaxPositionUSDC, req.RiskPreference)
	if err != nil {
		if errors.Is(err, copytrade.ErrNotFound) {
			response.NotFound(c, "wallet not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, cfg)
}

func (h *Handlers) ListCopyTradeConfigs(c *gin.Context) {
	configs, err := h.copyTradeService.ListConfigs(c.Request.Context(), userIdentifier(c))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, configs)
}

func (h *Handlers) GetCopyTradeConfig(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet_id")
		return
	}
	cfg, err := h.copyTradeService.GetConfig(c.Request.Context(), userIdentifier(c), walletID)
	if err != nil {
		if errors.Is(err, copytrade.ErrNotFound) {
			response.NotFound(c, "config not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, cfg)
}

func (h *Handlers) GetCopyTradeDashboard(c *gin.Context) {
	dashboard, err := h.copyTradeService.GetDashboard(c.Request.Context(), userIdentifier(c))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, dashboard)
}

func (h *Handlers) ListCopyTradeDecisions(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet_id")
		return
	}
	page, pageSize := parsePaging(c)
	result, err := h.copyTradeService.ListDecisions(c.Request.Context(), userIdentifier(c), walletID, page, pageSize)
	if err != nil {
		if errors.Is(err, copytrade.ErrNotFound) {
			response.NotFound(c, "config not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, result)
}

func (h *Handlers) GetCopyTradePerformance(c *gin.Context) {
	walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet_id")
		return
	}
	perf, err := h.copyTradeService.GetPerformance(c.Request.Context(), userIdentifier(c), walletID)
	if err != nil {
		if errors.Is(err, copytrade.ErrNotFound) {
			response.NotFound(c, "config not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, perf)
}

func (h *Handlers) ListCopyTradePositions(c *gin.Context) {
	positions, err := h.copyTradeService.ListOpenPositions(c.Request.Context(), userIdentifier(c))
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, positions)
}

func (h *Handlers) CloseCopyTradePosition(c *gin.Context) {
	decisionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid decision id")
		return
	}
	dec, err := h.copyTradeService.ClosePosition(c.Request.Context(), userIdentifier(c), decisionID)
	if err != nil {
		if errors.Is(err, copytrade.ErrNotFound) {
			response.NotFound(c, "position not found")
			return
		}
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, dec)
}

func (h *Handlers) GetCopyTradeMonitor(c *gin.Context) {
	ctx := c.Request.Context()

	runs, err := h.copyTradeRepo.GetSyncerRunHistory(ctx, 20)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}

	hourlyStats, err := h.copyTradeRepo.GetHourlySyncStats(ctx, 24)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}

	copyableWallets, err := h.copyTradeRepo.GetCopyableWallets(ctx, 20)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}

	enabledConfigs, err := h.copyTradeRepo.ListEnabledConfigs(ctx)
	if err != nil {
		response.Internal(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"enabled_configs":  len(enabledConfigs),
		"recent_runs":      runs,
		"hourly_stats":     hourlyStats,
		"copyable_wallets": copyableWallets,
	})
}
