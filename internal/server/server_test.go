package server

import (
  "context"
  "net/http"
  "net/http/httptest"
  "os"
  "testing"
  "time"
)

func TestNewServer(t *testing.T) {
  t.Run("creates server with default port", func(t *testing.T) {
    os.Unsetenv("PORT")
    srv := NewServer()

    if srv == nil {
      t.Fatal("NewServer() returned nil")
    }

    if srv.port != "8080" {
      t.Errorf("Expected default port 8080, got %s", srv.port)
    }
  })

  t.Run("creates server with custom port from environment", func(t *testing.T) {
    os.Setenv("PORT", "9000")
    defer os.Unsetenv("PORT")

    srv := NewServer()

    if srv.port != "9000" {
      t.Errorf("Expected port 9000, got %s", srv.port)
    }
  })
}

func TestServerSetupRoutes(t *testing.T) {
  srv := NewServer()
  handler := srv.setupRoutes()

  if handler == nil {
    t.Fatal("setupRoutes() returned nil handler")
  }

  // Test that routes are registered
  testCases := []struct {
    method string
    path   string
    status int
  }{
    {http.MethodGet, "/api/v1/health", http.StatusOK},
    {http.MethodGet, "/api/v1/status", http.StatusOK},
    {http.MethodGet, "/metrics", http.StatusOK},
    {http.MethodGet, "/openapi.yaml", http.StatusOK},
    {http.MethodGet, "/docs", http.StatusOK},
    // Evaluation endpoints
    {http.MethodPost, "/api/v1/evaluations/jobs", http.StatusAccepted},
    {http.MethodGet, "/api/v1/evaluations/jobs", http.StatusOK},
    {http.MethodGet, "/api/v1/evaluations/jobs/test-id", http.StatusOK},
    {http.MethodDelete, "/api/v1/evaluations/jobs/test-id", http.StatusOK},
    {http.MethodGet, "/api/v1/evaluations/jobs/test-id/summary", http.StatusOK},
    // Benchmarks
    {http.MethodGet, "/api/v1/evaluations/benchmarks", http.StatusOK},
    // Collections
    {http.MethodGet, "/api/v1/evaluations/collections", http.StatusOK},
    {http.MethodPost, "/api/v1/evaluations/collections", http.StatusCreated},
    {http.MethodGet, "/api/v1/evaluations/collections/test-collection", http.StatusOK},
    {http.MethodPut, "/api/v1/evaluations/collections/test-collection", http.StatusOK},
    {http.MethodPatch, "/api/v1/evaluations/collections/test-collection", http.StatusOK},
    {http.MethodDelete, "/api/v1/evaluations/collections/test-collection", http.StatusOK},
    // Providers
    {http.MethodGet, "/api/v1/evaluations/providers", http.StatusOK},
    {http.MethodGet, "/api/v1/evaluations/providers/test-provider", http.StatusOK},
    // System metrics
    {http.MethodGet, "/api/v1/metrics/system", http.StatusOK},
    // Error cases
    {http.MethodPost, "/api/v1/health", http.StatusMethodNotAllowed},
    {http.MethodGet, "/nonexistent", http.StatusNotFound},
  }

  for _, tc := range testCases {
    t.Run(tc.method+" "+tc.path, func(t *testing.T) {
      req := httptest.NewRequest(tc.method, tc.path, nil)
      w := httptest.NewRecorder()

      handler.ServeHTTP(w, req)

      if w.Code != tc.status {
        t.Errorf("Expected status %d for %s %s, got %d", tc.status, tc.method, tc.path, w.Code)
      }
    })
  }
}

func TestServerShutdown(t *testing.T) {
  t.Run("shutdown returns nil when server is nil", func(t *testing.T) {
    srv := &Server{
      httpServer: nil,
    }

    ctx := context.Background()
    err := srv.Shutdown(ctx)

    if err != nil {
      t.Errorf("Expected nil error when server is nil, got %v", err)
    }
  })

  t.Run("shutdown works with running server", func(t *testing.T) {
    srv := NewServer()
    srv.port = "0" // Use random port for testing

    // Start server in background
    errChan := make(chan error, 1)
    go func() {
      errChan <- srv.Start()
    }()

    // Wait a bit for server to start
    time.Sleep(100 * time.Millisecond)

    // Shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    err := srv.Shutdown(ctx)
    if err != nil {
      t.Errorf("Shutdown failed: %v", err)
    }

    // Wait for server to stop
    select {
    case err := <-errChan:
      if err != nil && err != http.ErrServerClosed {
        t.Errorf("Server error: %v", err)
      }
    case <-time.After(3 * time.Second):
      t.Error("Server did not stop within timeout")
    }
  })
}
