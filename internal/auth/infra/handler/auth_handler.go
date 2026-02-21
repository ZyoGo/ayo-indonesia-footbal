package handler

import (
	"net/http"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/auth/app"
	common "github.com/ZyoGo/ayo-indonesia-footbal/pkg/http"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	service app.AuthServicePort
}

func NewAuthHandler(service app.AuthServicePort) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		errResp := common.RenderErrorResponse(err)
		c.JSON(errResp.Code, errResp)
		return
	}

	c.JSON(http.StatusOK, common.NewSuccessResponseWithData(loginResponse{Token: token}))
}
