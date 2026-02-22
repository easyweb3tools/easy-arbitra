package api

import (
	"easy-arbitra/backend/internal/api/handler"
	"easy-arbitra/backend/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(30, 60))
	r.Use(middleware.ErrorHandler())

	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.Ready)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/wallets", h.ListWallets)
		v1.GET("/wallets/potential", h.ListPotentialWallets)
		v1.GET("/wallets/:id", h.GetWallet)
		v1.GET("/wallets/:id/profile", h.GetWalletProfile)
		v1.GET("/wallets/:id/share-card", h.GetWalletShareCard)
		v1.GET("/wallets/:id/explanations", h.GetWalletExplanations)
		v1.GET("/wallets/:id/info-edge", h.GetWalletInfoEdge)

		v1.GET("/markets", h.ListMarkets)
		v1.GET("/markets/:id", h.GetMarket)

		v1.GET("/leaderboard", h.GetLeaderboard)
		v1.GET("/stats/overview", h.GetOverviewStats)
		v1.GET("/ops/highlights", h.GetOpsHighlights)
		v1.GET("/anomalies", h.ListAnomalies)
		v1.GET("/anomalies/:id", h.GetAnomaly)
		v1.PATCH("/anomalies/:id/acknowledge", h.AcknowledgeAnomaly)

		ai := v1.Group("/ai")
		ai.POST("/analyze/:wallet_id", h.TriggerAIAnalysis)
		ai.GET("/report/:wallet_id", h.GetAIReport)
		ai.GET("/report/:wallet_id/history", h.ListAIReports)
	}

	return r
}
