package modules

import (
	clubApp "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/app"
	clubHandler "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/infra/handler"
	clubPostgres "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/infra/postgres"

	matchApp "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/app"
	matchHandler "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/infra/handler"
	matchPostgres "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/infra/postgres"

	reportingApp "github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/app"
	reportingHandler "github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/infra/handler"
	reportingPostgres "github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/infra/postgres"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterModules registers all application modules
func RegisterModules(db *pgxpool.Pool, router *gin.Engine) {
	api := router.Group("/api/v1")

	registerClubModule(db, api)
	registerMatchModule(db, api)
	registerReportingModule(db, api)
}

func registerClubModule(db *pgxpool.Pool, rg *gin.RouterGroup) {
	// Repositories
	teamRepo := clubPostgres.NewTeamRepository(db)
	playerRepo := clubPostgres.NewPlayerRepository(db)

	// Services
	teamService := clubApp.NewTeamService(teamRepo)
	playerService := clubApp.NewPlayerService(playerRepo, teamRepo)

	// Handlers
	teamH := clubHandler.NewTeamHandler(teamService)
	playerH := clubHandler.NewPlayerHandler(playerService)

	// Register routes
	clubHandler.RegisterRoutes(rg, teamH, playerH)
}

func registerMatchModule(db *pgxpool.Pool, rg *gin.RouterGroup) {
	// Repositories
	matchRepo := matchPostgres.NewMatchRepository(db)
	resultRepo := matchPostgres.NewMatchResultRepository(db)
	reportRepo := matchPostgres.NewReportRepository(db)

	// Service
	matchService := matchApp.NewMatchService(matchRepo, resultRepo, reportRepo)

	// Handler
	matchH := matchHandler.NewMatchHandler(matchService)

	// Register routes
	matchHandler.RegisterRoutes(rg, matchH)
}

func registerReportingModule(db *pgxpool.Pool, rg *gin.RouterGroup) {
	// Repositories
	repo := reportingPostgres.NewReportingRepository(db)

	// Service
	service := reportingApp.NewReportingService(repo)

	// Handler
	h := reportingHandler.NewReportingHandler(service)

	// Register routes
	reportingHandler.RegisterRoutes(rg, h)
}
