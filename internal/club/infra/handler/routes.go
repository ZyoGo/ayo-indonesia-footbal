package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all Club Management routes.
// Write routes (POST, PUT, DELETE) are protected by the auth middleware.
// Read routes (GET) are public.
func RegisterRoutes(rg *gin.RouterGroup, teamHandler *TeamHandler, playerHandler *PlayerHandler, authMiddleware ...gin.HandlerFunc) {
	// Team routes
	teams := rg.Group("/teams")
	{
		// Public (read-only)
		teams.GET("", teamHandler.GetAll)
		teams.GET("/:id", teamHandler.GetByID)
		teams.GET("/:id/players", playerHandler.GetByTeamID)

		// Protected (write) — middleware applied per-route
		teams.POST("", append(authMiddleware, teamHandler.Create)...)
		teams.PUT("/:id", append(authMiddleware, teamHandler.Update)...)
		teams.DELETE("/:id", append(authMiddleware, teamHandler.Delete)...)
	}

	// Player routes
	players := rg.Group("/players")
	{
		// Public (read-only)
		players.GET("/:id", playerHandler.GetByID)

		// Protected (write) — middleware applied per-route
		players.POST("", append(authMiddleware, playerHandler.Create)...)
		players.PUT("/:id", append(authMiddleware, playerHandler.Update)...)
		players.DELETE("/:id", append(authMiddleware, playerHandler.Delete)...)
	}
}
