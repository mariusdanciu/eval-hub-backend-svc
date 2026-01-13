package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"eval-hub-backend-svc/internal/constants"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// BackendSpec represents the backend specification
type BackendSpec struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// BenchmarkSpec represents the benchmark specification
type BenchmarkSpec struct {
	BenchmarkID string                 `json:"benchmark_id"`
	ProviderID  string                 `json:"provider_id"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// ExecutionContext contains execution context for evaluation operations
type ExecutionContext struct {
	EvaluationID   string                 `json:"evaluation_id"`
	ModelURL       string                 `json:"model_url"`
	ModelName      string                 `json:"model_name"`
	BackendSpec    BackendSpec            `json:"backend_spec"`
	BenchmarkSpec  BenchmarkSpec          `json:"benchmark_spec"`
	TimeoutMinutes int                    `json:"timeout_minutes"`
	RetryAttempts  int                    `json:"retry_attempts"`
	StartedAt      *time.Time             `json:"started_at,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	MLflowClient   interface{}            `json:"-"` // Not serialized
	ExperimentName *string                `json:"experiment_name,omitempty"`
	Logger         *zap.Logger            `json:"-"` // Not serialized
}

// enhanceLoggerWithRequest enhances a logger with request-specific fields
func enhanceLoggerWithRequest(logger *zap.Logger, r *http.Request) *zap.Logger {
	// Extract RequestID from X-Global-Transaction-Id header, or generate a UUID if not present
	requestID := r.Header.Get("X-Global-Transaction-Id")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// Add request_id to logger using With
	enhancedLogger := logger.With(zap.String(constants.LOG_REQUEST_ID, requestID))

	// Extract and add HTTP request fields to logger if they exist
	userAgent := r.Header.Get("User-Agent")
	if userAgent != "" {
		enhancedLogger = enhancedLogger.With(zap.String(constants.LOG_USER_AGENT, userAgent))
	}

	remoteAddr := r.RemoteAddr
	if remoteAddr != "" {
		enhancedLogger = enhancedLogger.With(zap.String(constants.LOG_REMOTE_ADR, remoteAddr))
	}

	// Extract remote_user from URL user info or header
	remoteUser := ""
	if r.URL != nil && r.URL.User != nil {
		remoteUser = r.URL.User.Username()
	}
	if remoteUser == "" {
		remoteUser = r.Header.Get("Remote-User")
	}
	if remoteUser != "" {
		enhancedLogger = enhancedLogger.With(zap.String(constants.LOG_USER, remoteUser))
	}

	referer := r.Header.Get("Referer")
	if referer != "" {
		enhancedLogger = enhancedLogger.With(zap.String(constants.LOG_REFERER, referer))
	}

	return enhancedLogger
}

// NewExecutionContext creates a new ExecutionContext with default values
func NewExecutionContext(r *http.Request, logger *zap.Logger) ExecutionContext {
	// Enhance logger with request-specific fields
	enhancedLogger := enhanceLoggerWithRequest(logger, r)

	return ExecutionContext{
		TimeoutMinutes: 60,
		RetryAttempts:  3,
		BackendSpec:    BackendSpec{},
		BenchmarkSpec:  BenchmarkSpec{},
		Metadata:       make(map[string]interface{}),
		Logger:         enhancedLogger,
	}
}

// HandleCreateEvaluation handles POST /api/v1/evaluations/jobs
func (h *Handlers) HandleCreateEvaluation(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Evaluation creation not yet implemented",
	})
}

// HandleListEvaluations handles GET /api/v1/evaluations/jobs
func (h *Handlers) HandleListEvaluations(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":       []interface{}{},
		"total_count": 0,
		"limit":       50,
		"first":       map[string]string{"href": ""},
		"next":        nil,
	})
}

// HandleGetEvaluation handles GET /api/v1/evaluations/jobs/{id}
func (h *Handlers) HandleGetEvaluation(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	pathParts := strings.Split(r.URL.Path, "/")
	id := pathParts[len(pathParts)-1]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Evaluation retrieval not yet implemented",
		"id":      id,
	})
}

// HandleCancelEvaluation handles DELETE /api/v1/evaluations/jobs/{id}
func (h *Handlers) HandleCancelEvaluation(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Evaluation cancellation not yet implemented",
	})
}

// HandleGetEvaluationSummary handles GET /api/v1/evaluations/jobs/{id}/summary
func (h *Handlers) HandleGetEvaluationSummary(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Evaluation summary not yet implemented",
	})
}

// HandleListBenchmarks handles GET /api/v1/evaluations/benchmarks
func (h *Handlers) HandleListBenchmarks(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"benchmarks":         []interface{}{},
		"total_count":        0,
		"providers_included": []string{},
	})
}

// HandleListCollections handles GET /api/v1/evaluations/collections
func (h *Handlers) HandleListCollections(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"collections":       []interface{}{},
		"total_collections": 0,
	})
}

// HandleCreateCollection handles POST /api/v1/evaluations/collections
func (h *Handlers) HandleCreateCollection(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection creation not yet implemented",
	})
}

// HandleGetCollection handles GET /api/v1/evaluations/collections/{collection_id}
func (h *Handlers) HandleGetCollection(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract collection_id from path
	pathParts := strings.Split(r.URL.Path, "/")
	collectionID := pathParts[len(pathParts)-1]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Collection retrieval not yet implemented",
		"collection_id": collectionID,
	})
}

// HandleUpdateCollection handles PUT /api/v1/evaluations/collections/{collection_id}
func (h *Handlers) HandleUpdateCollection(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection update not yet implemented",
	})
}

// HandlePatchCollection handles PATCH /api/v1/evaluations/collections/{collection_id}
func (h *Handlers) HandlePatchCollection(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection patch not yet implemented",
	})
}

// HandleDeleteCollection handles DELETE /api/v1/evaluations/collections/{collection_id}
func (h *Handlers) HandleDeleteCollection(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection deletion not yet implemented",
	})
}

// HandleListProviders handles GET /api/v1/evaluations/providers
func (h *Handlers) HandleListProviders(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"providers":        []interface{}{},
		"total_providers":  0,
		"total_benchmarks": 0,
	})
}

// HandleGetProvider handles GET /api/v1/evaluations/providers/{provider_id}
func (h *Handlers) HandleGetProvider(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract provider_id from path
	pathParts := strings.Split(r.URL.Path, "/")
	providerID := pathParts[len(pathParts)-1]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Provider retrieval not yet implemented",
		"provider_id": providerID,
	})
}

// HandleGetSystemMetrics handles GET /api/v1/metrics/system
func (h *Handlers) HandleGetSystemMetrics(ctx ExecutionContext, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "System metrics not yet implemented",
	})
}
