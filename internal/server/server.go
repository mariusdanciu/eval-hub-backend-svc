package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"eval-hub-backend-svc/internal/handlers"
	"eval-hub-backend-svc/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
	port       string
}

func NewServer() *Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{
		port: port,
	}
}

func (s *Server) setupRoutes() http.Handler {
	router := http.NewServeMux()
	h := handlers.New()

	// Health and status endpoints
	router.HandleFunc("/api/v1/health", h.HandleHealth)
	router.HandleFunc("/api/v1/status", h.HandleStatus)

	// Evaluation jobs endpoints
	router.HandleFunc("/api/v1/evaluations/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.HandleCreateEvaluation(w, r)
		case http.MethodGet:
			h.HandleListEvaluations(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	// Handle summary endpoint first (more specific)
	router.HandleFunc("/api/v1/evaluations/jobs/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/summary") && r.Method == http.MethodGet {
			h.HandleGetEvaluationSummary(w, r)
			return
		}
		// Handle individual job endpoints
		switch r.Method {
		case http.MethodGet:
			h.HandleGetEvaluation(w, r)
		case http.MethodDelete:
			h.HandleCancelEvaluation(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Benchmarks endpoint
	router.HandleFunc("/api/v1/evaluations/benchmarks", h.HandleListBenchmarks)

	// Collections endpoints
	router.HandleFunc("/api/v1/evaluations/collections", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.HandleCreateCollection(w, r)
		case http.MethodGet:
			h.HandleListCollections(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc("/api/v1/evaluations/collections/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.HandleGetCollection(w, r)
		case http.MethodPut:
			h.HandleUpdateCollection(w, r)
		case http.MethodPatch:
			h.HandlePatchCollection(w, r)
		case http.MethodDelete:
			h.HandleDeleteCollection(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Providers endpoints
	router.HandleFunc("/api/v1/evaluations/providers", h.HandleListProviders)
	router.HandleFunc("/api/v1/evaluations/providers/", h.HandleGetProvider)

	// System metrics endpoint
	router.HandleFunc("/api/v1/metrics/system", h.HandleGetSystemMetrics)

	// OpenAPI documentation endpoints
	router.HandleFunc("/openapi.yaml", h.HandleOpenAPI)
	router.HandleFunc("/docs", h.HandleDocs)

	// Prometheus metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// Wrap router with metrics middleware
	return metrics.Middleware(router)
}

// SetupRoutes exposes the route setup for testing
func (s *Server) SetupRoutes() http.Handler {
	return s.setupRoutes()
}

// GetPort returns the server port
func (s *Server) GetPort() string {
	return s.port
}

// SetPort sets the server port (for testing)
func (s *Server) SetPort(port string) {
	s.port = port
}

func (s *Server) Start() error {
	handler := s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:         ":" + s.port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Server starting on port %s\n", s.port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	fmt.Println("Shutting down server gracefully...")
	return s.httpServer.Shutdown(ctx)
}
