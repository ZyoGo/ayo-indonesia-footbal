package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all Match Context routes.
func RegisterRoutes(rg *gin.RouterGroup, matchHandler *MatchHandler) {
	matches := rg.Group("/matches")
	{
		matches.POST("", matchHandler.CreateMatch)
		matches.GET("", matchHandler.GetAllMatches)
		matches.GET("/:id", matchHandler.GetMatchByID)
		matches.POST("/:id/result", matchHandler.ReportResult)
		matches.GET("/:id/report", matchHandler.GetMatchReport)
	}

	// Reports
	rg.GET("/reports/matches", matchHandler.GetAllMatchReports)
}
