package modules

import (
	authApp "github.com/ZyoGo/ayo-indonesia-footbal/internal/auth/app"
	authHandler "github.com/ZyoGo/ayo-indonesia-footbal/internal/auth/infra/handler"
	authPostgres "github.com/ZyoGo/ayo-indonesia-footbal/internal/auth/infra/postgres"

	clubApp "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/app"
	clubHandler "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/infra/handler"
	clubPostgres "github.com/ZyoGo/ayo-indonesia-footbal/internal/club/infra/postgres"

	matchApp "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/app"
	matchHandler "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/infra/handler"
	matchPostgres "github.com/ZyoGo/ayo-indonesia-footbal/internal/match/infra/postgres"

	reportingApp "github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/app"
	reportingHandler "github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/infra/handler"
	reportingPostgres "github.com/ZyoGo/ayo-indonesia-footbal/internal/reporting/infra/postgres"

	uploadHandler "github.com/ZyoGo/ayo-indonesia-footbal/internal/upload/handler"

	"github.com/ZyoGo/ayo-indonesia-footbal/config"
	authguard "github.com/ZyoGo/ayo-indonesia-footbal/pkg/http/middleware/authguard"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/jwt"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/upload"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterModules(db *pgxpool.Pool, router *gin.Engine) {
	cfg := config.GetConfig()
	jwtKeys := config.GetJWTKeys()

	jwtService := jwt.NewService(jwtKeys.PrivateKey, jwtKeys.PublicKey, cfg.JWT.Issuer, cfg.JWT.Subject)

	guard := authguard.NewAuthGuard(jwtService)
	authMW := guard.Guard()

	uploader := upload.NewUploader(
		cfg.Upload.BasePath,
		cfg.Upload.URLPrefix,
		cfg.Upload.MaxSize,
		cfg.Upload.AllowedTypes,
	)

	router.Static("/uploads", cfg.Upload.BasePath)

	api := router.Group("/api/v1")

	registerAuthModule(db, api, jwtService)
	registerUploadModule(api, uploader, authMW)
	registerClubModule(db, api, authMW)
	registerMatchModule(db, api, authMW)
	registerReportingModule(db, api)
}

func registerAuthModule(db *pgxpool.Pool, rg *gin.RouterGroup, jwtService *jwt.Service) {
	userRepo := authPostgres.NewUserRepository(db)
	service := authApp.NewAuthService(userRepo, jwtService)
	h := authHandler.NewAuthHandler(service)
	authHandler.RegisterRoutes(rg, h)
}

func registerUploadModule(rg *gin.RouterGroup, uploader *upload.Uploader, authMW gin.HandlerFunc) {
	uploadH := uploadHandler.NewUploadHandler(uploader)
	rg.POST("/uploads", authMW, uploadH.Upload)
}

func registerClubModule(db *pgxpool.Pool, rg *gin.RouterGroup, authMW gin.HandlerFunc) {
	teamRepo := clubPostgres.NewTeamRepository(db)
	playerRepo := clubPostgres.NewPlayerRepository(db)

	teamService := clubApp.NewTeamService(teamRepo)
	playerService := clubApp.NewPlayerService(playerRepo, teamRepo)

	teamH := clubHandler.NewTeamHandler(teamService)
	playerH := clubHandler.NewPlayerHandler(playerService)

	clubHandler.RegisterRoutes(rg, teamH, playerH, authMW)
}

func registerMatchModule(db *pgxpool.Pool, rg *gin.RouterGroup, authMW gin.HandlerFunc) {
	matchRepo := matchPostgres.NewMatchRepository(db)
	resultRepo := matchPostgres.NewMatchResultRepository(db)
	reportRepo := matchPostgres.NewReportRepository(db)

	matchService := matchApp.NewMatchService(matchRepo, resultRepo, reportRepo)

	matchH := matchHandler.NewMatchHandler(matchService)

	matchHandler.RegisterRoutes(rg, matchH, authMW)
}

func registerReportingModule(db *pgxpool.Pool, rg *gin.RouterGroup) {
	repo := reportingPostgres.NewReportingRepository(db)
	service := reportingApp.NewReportingService(repo)
	h := reportingHandler.NewReportingHandler(service)
	reportingHandler.RegisterRoutes(rg, h)
}
