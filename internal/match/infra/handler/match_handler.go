package handler

import (
	"net/http"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/app"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/infra/handler/request"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/infra/handler/response"
	common "github.com/ZyoGo/ayo-indonesia-footbal/pkg/http"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	service app.MatchServicePort
}

func NewMatchHandler(service app.MatchServicePort) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var req request.CreateMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		return
	}

	id, err := h.service.CreateMatch(c.Request.Context(), req.ToDomain())
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewCreatedSuccessResponse(id))
}

func (h *MatchHandler) GetAllMatches(c *gin.Context) {
	matches, err := h.service.GetAllMatches(c.Request.Context())
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(response.FromMatches(matches)))
}

func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	id := c.Param("id")

	match, err := h.service.GetMatchByID(c.Request.Context(), id)
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(response.FromMatch(match)))
}

func (h *MatchHandler) ReportResult(c *gin.Context) {
	matchID := c.Param("id")

	var req request.ReportResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		return
	}

	id, err := h.service.ReportResult(c.Request.Context(), matchID, req.ToDomain())
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewCreatedSuccessResponse(id))
}

func (h *MatchHandler) GetMatchReport(c *gin.Context) {
	matchID := c.Param("id")

	report, err := h.service.GetMatchReport(c.Request.Context(), matchID)
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(response.FromMatchReport(report)))
}

func (h *MatchHandler) GetAllMatchReports(c *gin.Context) {
	reports, err := h.service.GetAllMatchReports(c.Request.Context())
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(response.FromMatchReports(reports)))
}
