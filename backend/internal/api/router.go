package api

import (
	"easy-arbitra/backend/internal/api/handler"
	"easy-arbitra/backend/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.Handlers, frontendURL string) *gin.Engine {
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
		// Wallets
		v1.GET("/wallets", h.ListWallets)
		v1.GET("/wallets/potential", h.ListPotentialWallets)
		v1.GET("/wallets/:id", h.GetWallet)
		v1.GET("/wallets/:id/profile", h.GetWalletProfile)
		v1.GET("/wallets/:id/trades", h.ListWalletTrades)
		v1.GET("/wallets/:id/positions", h.ListWalletPositions)

		// Markets
		v1.GET("/markets", h.ListMarkets)
		v1.GET("/markets/:id", h.GetMarket)

		// Leaderboard & Stats
		v1.GET("/leaderboard", h.GetLeaderboard)
		v1.GET("/stats/overview", h.GetOverviewStats)

		// Daily Pick
		v1.GET("/daily-pick", h.GetDailyPick)
		v1.GET("/daily-pick/history", h.ListDailyPickHistory)

		// Nova Sessions (thinking timeline)
		v1.GET("/nova/sessions", h.ListNovaSessions)

		// Nova Insight (AI brain visualization)
		v1.GET("/nova/status", h.GetNovaStatus)
		v1.GET("/nova/timeline/:date", h.GetNovaTimeline)
		v1.GET("/nova/candidates/:date", h.GetNovaCandidates)
		v1.GET("/nova/decision-explain/:pick_id", h.GetNovaDecisionExplanation)
		v1.GET("/nova/memory", h.GetNovaMemory)
	}

	return r
}
