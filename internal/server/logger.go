package server

import (
	"go.uber.org/zap"
)

// NewLogger creates a new zap logger.
// It creates a production JSON logger that outputs to stdout.
// If production logger creation fails, it falls back to development logger.
func NewLogger() *zap.Logger {
	// Create a production JSON logger (outputs JSON to stdout)
	logger, err := zap.NewProduction()
	if err != nil {
		// Fallback to development logger if production logger fails
		logger, _ = zap.NewDevelopment()
	}
	return logger
}
