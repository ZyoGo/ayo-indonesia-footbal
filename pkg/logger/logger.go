package logger

import (
	"log/slog"
	"os"
	"sync"
)

// Global logger instance
var (
	logger *slog.Logger
	once   sync.Once
)

// InitializeLogger sets up structured logging with slog
func InitializeLogger() *slog.Logger {
	once.Do(func() {
		// Write to stdout only
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	})

	return logger
}

// Get returns the global logger instance
func Get() *slog.Logger {
	return InitializeLogger()
}
