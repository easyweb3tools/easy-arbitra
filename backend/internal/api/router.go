package api

import (
	"easy-arbitra/backend/internal/api/handler"
	"easy-arbitra/backend/internal/api/middleware"
	"easy-arbitra/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.Handlers, authHandler *auth.Handler, jwtSecret string, frontendURL string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(frontendURL))
	r.Use(middleware.RateLimit(30, 60))
	r.Use(middleware.ErrorHandler())

	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.Ready)

	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		authGroup := v1.Group("/auth")
		authHandler.RegisterRoutes(authGroup)

		// Public routes
		v1.GET("/wallets", h.ListWallets)
		v1.GET("/wallets/potential", h.ListPotentialWallets)
		v1.GET("/wallets/:id", h.GetWallet)
		v1.GET("/wallets/:id/profile", h.GetWalletProfile)
		v1.GET("/wallets/:id/share-card", h.GetWalletShareCard)
		v1.GET("/wallets/:id/share-landing", h.GetWalletShareLanding)
		v1.GET("/wallets/:id/decision-card", h.GetWalletDecisionCard)
		v1.GET("/wallets/:id/explanations", h.GetWalletExplanations)
		v1.GET("/wallets/:id/info-edge", h.GetWalletInfoEdge)
		v1.GET("/wallets/:id/pnl-history", h.GetWalletPnLHistory)
		v1.GET("/wallets/:id/trades", h.ListWalletTrades)
		v1.GET("/wallets/:id/positions", h.ListWalletPositions)
		v1.GET("/portfolios", h.ListPortfolios)

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

		// Public copy-trading endpoint
		v1.GET("/copy-trading/monitor", h.GetCopyTradeMonitor)

		// Protected routes
		protected := v1.Group("")
		protected.Use(auth.AuthRequired(jwtSecret))

		// Protected watchlist routes
		protected.GET("/watchlist", h.ListWatchlist)
		protected.POST("/watchlist", h.AddToWatchlist)
		protected.POST("/watchlist/batch", h.AddToWatchlistBatch)
		protected.DELETE("/watchlist/:wallet_id", h.RemoveFromWatchlist)
		protected.GET("/watchlist/feed", h.GetWatchlistFeed)
		protected.GET("/watchlist/summary", h.GetWatchlistSummary)

		// Protected copy-trading routes
		ct := protected.Group("/copy-trading")
		ct.POST("/enable", h.EnableCopyTrading)
		ct.POST("/disable", h.DisableCopyTrading)
		ct.PUT("/settings", h.UpdateCopyTradeSettings)
		ct.GET("/configs", h.ListCopyTradeConfigs)
		ct.GET("/dashboard", h.GetCopyTradeDashboard)
		ct.GET("/positions", h.ListCopyTradePositions)
		ct.GET("/:wallet_id", h.GetCopyTradeConfig)
		ct.GET("/:wallet_id/decisions", h.ListCopyTradeDecisions)
		ct.GET("/:wallet_id/performance", h.GetCopyTradePerformance)
		ct.POST("/decisions/:id/close", h.CloseCopyTradePosition)
	}

	return r
}
