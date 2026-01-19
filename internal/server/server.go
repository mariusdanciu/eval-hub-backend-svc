package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/config"
	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/execution_context"
	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/handlers"
	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer    *http.Server
	port          int
	logger        *slog.Logger
	serviceConfig *config.Config
}

func NewServer(logger *slog.Logger, serviceConfig *config.Config) (*Server, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required for the server")
	}
	if (serviceConfig == nil) || (serviceConfig.Service == nil) {
		return nil, fmt.Errorf("service config is required for the server")
	}

	return &Server{
		port:          serviceConfig.Service.Port,
		logger:        logger,
		serviceConfig: serviceConfig,
	}, nil
}

func (s *Server) setupRoutes() (http.Handler, error) {
	router := http.NewServeMux()
	h := handlers.New()

	// Health and status endpoints
	router.HandleFunc("/api/v1/health", h.HandleHealth)
	router.HandleFunc("/api/v1/status", h.HandleStatus)

	// Evaluation jobs endpoints
	router.HandleFunc("/api/v1/evaluations/jobs", func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.NewExecutionContext(r, s.logger, s.serviceConfig)
		switch r.Method {
		case http.MethodPost:
			h.HandleCreateEvaluation(ctx, w, r)
		case http.MethodGet:
			h.HandleListEvaluations(ctx, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// Handle summary endpoint first (more specific)
	router.HandleFunc("/api/v1/evaluations/jobs/", func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.NewExecutionContext(r, s.logger, s.serviceConfig)
		path := r.URL.Path
		if strings.HasSuffix(path, "/summary") && r.Method == http.MethodGet {
			h.HandleGetEvaluationSummary(ctx, w, r)
			return
		}
		// Handle individual job endpoints
		switch r.Method {
		case http.MethodGet:
			h.HandleGetEvaluation(ctx, w, r)
		case http.MethodDelete:
			h.HandleCancelEvaluation(ctx, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Benchmarks endpoint
	router.HandleFunc("/api/v1/evaluations/benchmarks", func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.NewExecutionContext(r, s.logger, s.serviceConfig)
		h.HandleListBenchmarks(ctx, w, r)
	})

	// Collections endpoints
	router.HandleFunc("/api/v1/evaluations/collections", func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.NewExecutionContext(r, s.logger, s.serviceConfig)
		switch r.Method {
		case http.MethodPost:
			h.HandleCreateCollection(ctx, w, r)
		case http.MethodGet:
			h.HandleListCollections(ctx, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc("/api/v1/evaluations/collections/", func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.NewExecutionContext(r, s.logger, s.serviceConfig)
		switch r.Method {
		case http.MethodGet:
			h.HandleGetCollection(ctx, w, r)
		case http.MethodPut:
			h.HandleUpdateCollection(ctx, w, r)
		case http.MethodPatch:
			h.HandlePatchCollection(ctx, w, r)
		case http.MethodDelete:
			h.HandleDeleteCollection(ctx, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Providers endpoints
	router.HandleFunc("/api/v1/evaluations/providers", func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.NewExecutionContext(r, s.logger, s.serviceConfig)
		h.HandleListProviders(ctx, w, r)
	})
	router.HandleFunc("/api/v1/evaluations/providers/", func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.NewExecutionContext(r, s.logger, s.serviceConfig)
		h.HandleGetProvider(ctx, w, r)
	})

	// System metrics endpoint
	router.HandleFunc("/api/v1/metrics/system", func(w http.ResponseWriter, r *http.Request) {
		ctx := execution_context.NewExecutionContext(r, s.logger, s.serviceConfig)
		h.HandleGetSystemMetrics(ctx, w, r)
	})

	// OpenAPI documentation endpoints
	router.HandleFunc("/openapi.yaml", h.HandleOpenAPI)
	router.HandleFunc("/docs", h.HandleDocs)

	// Prometheus metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// Wrap router with metrics middleware
	return metrics.Middleware(router), nil
}

// SetupRoutes exposes the route setup for testing
func (s *Server) SetupRoutes() (http.Handler, error) {
	return s.setupRoutes()
}

func (s *Server) Start() error {
	handler, err := s.setupRoutes()
	if err != nil {
		return err
	}
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Server starting on port %d\n", s.port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	fmt.Println("Shutting down server gracefully...")
	return s.httpServer.Shutdown(ctx)
}
