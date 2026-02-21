package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/cmd/http/modules"
	"github.com/ZyoGo/ayo-indonesia-footbal/config"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/db"
	middleware "github.com/ZyoGo/ayo-indonesia-footbal/pkg/http/middleware/logger"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()
	dbConn := db.NewDatabaseConnection()
	defer dbConn.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.GinLogger())

	// Register modules
	modules.RegisterModules(dbConn, router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: router,
	}

	go func() {
		logger.Get().With().Info(fmt.Sprintf("Starting server on port %d", cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Get().With().ErrorContext(context.Background(), "Server failed", slog.Any("error", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Get().With().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Get().With().ErrorContext(ctx, "Server forced to shutdown", slog.Any("error", err))
	}

	logger.Get().With().Info("Server stopped gracefully")
}
