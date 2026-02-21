package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all Auth routes (public, no middleware).
func RegisterRoutes(rg *gin.RouterGroup, authHandler *AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}
}
