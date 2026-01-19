package execution_context

import (
	"log/slog"
	"net/http"
	"time"

	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/config"
	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/logging"
)

// ExecutionContext contains execution context for API operations
type ExecutionContext struct {
	Logger       *slog.Logger
	Config       *config.Config
	EvaluationID string
	ModelURL     string
	ModelName    string
	//BackendSpec    BackendSpec
	//BenchmarkSpec  BenchmarkSpec
	TimeoutMinutes int
	RetryAttempts  int
	StartedAt      *time.Time
	Metadata       map[string]interface{}
	MLflowClient   interface{}
	ExperimentName *string
}

// NewExecutionContext creates a new ExecutionContext with default values
func NewExecutionContext(r *http.Request, logger *slog.Logger, serviceConfig *config.Config) *ExecutionContext {
	// Enhance logger with request-specific fields
	enhancedLogger := logging.LoggerWithRequest(logger, r)

	return &ExecutionContext{
		Logger:         enhancedLogger,
		Config:         serviceConfig,
		TimeoutMinutes: 60,
		RetryAttempts:  3,
		Metadata:       make(map[string]interface{}),
	}
}
