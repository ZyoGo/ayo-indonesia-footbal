package handler

import (
	"net/http"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/app"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/infra/handler/request"
	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/infra/handler/response"
	common "github.com/ZyoGo/ayo-indonesia-footbal/pkg/http"
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	service app.PlayerServicePort
}

func NewPlayerHandler(service app.PlayerServicePort) *PlayerHandler {
	return &PlayerHandler{service: service}
}

func (h *PlayerHandler) Create(c *gin.Context) {
	var req request.CreatePlayerRequest
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

func (h *PlayerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	player, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(response.FromPlayer(player)))
}

func (h *PlayerHandler) GetByTeamID(c *gin.Context) {
	teamID := c.Param("id")

	players, err := h.service.GetByTeamID(c.Request.Context(), teamID)
	if err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(response.FromPlayers(players)))
}

func (h *PlayerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdatePlayerRequest
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

func (h *PlayerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		resp := common.RenderErrorResponse(err)
		c.JSON(resp.Code, resp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponse())
}
