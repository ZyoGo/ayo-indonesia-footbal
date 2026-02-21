package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/config"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabaseConnection() *pgxpool.Pool {
	const (
		defaultMaxConns          int32         = 50
		defaultMinConns          int32         = 2
		defaultMaxConnLifetime   time.Duration = 30 * time.Minute
		defaultMaxConnIdleTime   time.Duration = 5 * time.Minute
		defaultHealthCheckPeriod time.Duration = 30 * time.Second
		defaultConnTimeout       time.Duration = 5 * time.Second
	)

	// Load database configuration from environment variables
	// example: postgres://user:password@localhost:5432/dbname?sslmode=disable
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.GetConfig().Database.Username,
		config.GetConfig().Database.Password,
		config.GetConfig().Database.Address,
		config.GetConfig().Database.Port,
		config.GetConfig().Database.Name,
	)

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Get().With().ErrorContext(context.Background(), "Failed to create database config", slog.Any("error", err))
		os.Exit(1)
	}

	// Apply pooling settings
	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnTimeout

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Get().With().ErrorContext(context.Background(), "Failed to establish database pool", slog.Any("error", err))
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		logger.Get().With().ErrorContext(context.Background(), "Failed to ping database connection", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Get().With().Info("Connected to the database successfully!")
	return pool
}
