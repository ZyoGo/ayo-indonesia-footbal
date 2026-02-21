package handler

import (
	"net/http"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/app"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/infra/handler/response"
	"github.com/gin-gonic/gin"
)

type ReportingHandler struct {
	service app.ReportingServicePort
}

func NewReportingHandler(service app.ReportingServicePort) *ReportingHandler {
	return &ReportingHandler{service: service}
}

func (h *ReportingHandler) GetStandings(c *gin.Context) {
	standings, err := h.service.GetStandings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.StandingResponse, 0, len(standings))
	for _, s := range standings {
		resp = append(resp, response.FromStandingDomain(s))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ReportingHandler) GetTopScorers(c *gin.Context) {
	scorers, err := h.service.GetTopScorers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]response.TopScorerResponse, 0, len(scorers))
	for _, s := range scorers {
		resp = append(resp, response.FromTopScorerDomain(s))
	}

	c.JSON(http.StatusOK, resp)
}

func RegisterRoutes(rg *gin.RouterGroup, h *ReportingHandler) {
	reporting := rg.Group("/reporting")
	{
		reporting.GET("/standings", h.GetStandings)
		reporting.GET("/top-scorers", h.GetTopScorers)
	}
}
