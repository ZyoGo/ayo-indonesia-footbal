package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all Club Management routes.
// Routes are separated from handler logic for clarity.
func RegisterRoutes(rg *gin.RouterGroup, teamHandler *TeamHandler, playerHandler *PlayerHandler) {
	// Team routes
	teams := rg.Group("/teams")
	{
		teams.POST("", teamHandler.Create)
		teams.GET("", teamHandler.GetAll)
		teams.GET("/:id", teamHandler.GetByID)
		teams.PUT("/:id", teamHandler.Update)
		teams.DELETE("/:id", teamHandler.Delete)

		// Nested: players under a team
		teams.GET("/:id/players", playerHandler.GetByTeamID)
	}

	// Player routes
	players := rg.Group("/players")
	{
		players.POST("", playerHandler.Create)
		players.GET("/:id", playerHandler.GetByID)
		players.PUT("/:id", playerHandler.Update)
		players.DELETE("/:id", playerHandler.Delete)
	}
}
