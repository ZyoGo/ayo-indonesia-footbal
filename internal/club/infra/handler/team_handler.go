package handler

import (
	"net/http"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/app"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/infra/handler/request"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/infra/handler/response"
	common "github.com/ZyoGo/ayo-indonesia-footbal/pkg/http"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service app.TeamServicePort
}

func NewTeamHandler(service app.TeamServicePort) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req request.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		return
	}

	id, err := h.service.Create(c.Request.Context(), req.ToDomain())
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewCreatedSuccessResponse(id))
}

func (h *TeamHandler) GetAll(c *gin.Context) {
	teams, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(response.FromTeams(teams)))
}

func (h *TeamHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	team, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(response.FromTeam(team)))
}

func (h *TeamHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req.ToDomain()); err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse())
}

func (h *TeamHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse())
}
