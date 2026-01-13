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

  // Create logger once for all requests
  logger := NewLogger()

  // Health and status endpoints
  router.HandleFunc("/api/v1/health", h.HandleHealth)
  router.HandleFunc("/api/v1/status", h.HandleStatus)

  // Evaluation jobs endpoints
  router.HandleFunc("/api/v1/evaluations/jobs", func(w http.ResponseWriter, r *http.Request) {
    ctx := handlers.NewExecutionContext(r, logger)
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
    ctx := handlers.NewExecutionContext(r, logger)
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
    ctx := handlers.NewExecutionContext(r, logger)
    h.HandleListBenchmarks(ctx, w, r)
  })

  // Collections endpoints
  router.HandleFunc("/api/v1/evaluations/collections", func(w http.ResponseWriter, r *http.Request) {
    ctx := handlers.NewExecutionContext(r, logger)
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
    ctx := handlers.NewExecutionContext(r, logger)
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
    ctx := handlers.NewExecutionContext(r, logger)
    h.HandleListProviders(ctx, w, r)
  })
  router.HandleFunc("/api/v1/evaluations/providers/", func(w http.ResponseWriter, r *http.Request) {
    ctx := handlers.NewExecutionContext(r, logger)
    h.HandleGetProvider(ctx, w, r)
  })

  // System metrics endpoint
  router.HandleFunc("/api/v1/metrics/system", func(w http.ResponseWriter, r *http.Request) {
    ctx := handlers.NewExecutionContext(r, logger)
    h.HandleGetSystemMetrics(ctx, w, r)
  })

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
