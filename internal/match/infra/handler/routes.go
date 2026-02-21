package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all Match Context routes.
// Write routes (POST) are protected by the auth middleware.
// Read routes (GET) are public.
func RegisterRoutes(rg *gin.RouterGroup, matchHandler *MatchHandler, authMiddleware ...gin.HandlerFunc) {
	matches := rg.Group("/matches")
	{
		// Public (read-only)
		matches.GET("", matchHandler.GetAllMatches)
		matches.GET("/:id", matchHandler.GetMatchByID)
		matches.GET("/:id/report", matchHandler.GetMatchReport)

		// Protected (write) — middleware applied per-route
		matches.POST("", append(authMiddleware, matchHandler.CreateMatch)...)
		matches.POST("/:id/result", append(authMiddleware, matchHandler.ReportResult)...)
	}

	// Reports (public, read-only)
	rg.GET("/reports/matches", matchHandler.GetAllMatchReports)
}
