package handler

import (
	"strconv"

	"easy-arbitra/backend/internal/service"
	"easy-arbitra/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetWalletDecisionCard(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	row, err := h.walletService.GetDecisionCard(c.Request.Context(), id)
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

func (h *Handlers) GetWalletShareLanding(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid wallet id")
		return
	}
	row, err := h.walletService.GetShareLanding(c.Request.Context(), id)
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

func (h *Handlers) ListPortfolios(c *gin.Context) {
	if h.portfolioService == nil {
		response.OK(c, []any{})
		return
	}
	rows, err := h.portfolioService.List(c.Request.Context())
	if err != nil {
		response.Internal(c, err.Error())
		return
	}
	response.OK(c, rows)
}
