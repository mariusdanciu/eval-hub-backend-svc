package handlers

import (
  "net/http"
  "net/http/httptest"
  "os"
  "path/filepath"
  "testing"
)

func TestHandleOpenAPI(t *testing.T) {
  h := New()

  // Ensure the OpenAPI file exists for testing
  apiPath := filepath.Join("..", "..", "api", "openapi.yaml")
  if _, err := os.Stat(apiPath); os.IsNotExist(err) {
    // Try alternative path
    apiPath = "api/openapi.yaml"
    if _, err := os.Stat(apiPath); os.IsNotExist(err) {
      t.Skip("OpenAPI spec file not found, skipping test")
    }
  }

  t.Run("GET request returns OpenAPI spec", func(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
    w := httptest.NewRecorder()

    h.HandleOpenAPI(w, req)

    if w.Code != http.StatusOK {
      t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }

    contentType := w.Header().Get("Content-Type")
    if contentType != "application/yaml" && contentType != "application/json" {
      t.Errorf("Expected Content-Type application/yaml or application/json, got %s", contentType)
    }

    if len(w.Body.Bytes()) == 0 {
      t.Error("Response body is empty")
    }

    // Check if response contains OpenAPI keywords
    body := w.Body.String()
    if !contains(body, "openapi") && !contains(body, "OpenAPI") {
      t.Error("Response does not appear to be an OpenAPI specification")
    }
  })

  t.Run("POST request returns method not allowed", func(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/openapi.yaml", nil)
    w := httptest.NewRecorder()

    h.HandleOpenAPI(w, req)

    if w.Code != http.StatusMethodNotAllowed {
      t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
    }
  })

  t.Run("JSON content type when Accept header is application/json", func(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
    req.Header.Set("Accept", "application/json")
    w := httptest.NewRecorder()

    h.HandleOpenAPI(w, req)

    contentType := w.Header().Get("Content-Type")
    if contentType != "application/json" {
      t.Errorf("Expected Content-Type application/json, got %s", contentType)
    }
  })
}

func TestHandleDocs(t *testing.T) {
  h := New()

  t.Run("GET request returns HTML documentation", func(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/docs", nil)
    w := httptest.NewRecorder()

    h.HandleDocs(w, req)

    if w.Code != http.StatusOK {
      t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }

    contentType := w.Header().Get("Content-Type")
    if contentType != "text/html; charset=utf-8" {
      t.Errorf("Expected Content-Type text/html; charset=utf-8, got %s", contentType)
    }

    body := w.Body.String()
    if !contains(body, "swagger-ui") && !contains(body, "SwaggerUI") {
      t.Error("Response does not appear to be Swagger UI HTML")
    }

    if !contains(body, "openapi.yaml") {
      t.Error("Response does not reference openapi.yaml")
    }
  })

  t.Run("POST request returns method not allowed", func(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/docs", nil)
    w := httptest.NewRecorder()

    h.HandleDocs(w, req)

    if w.Code != http.StatusMethodNotAllowed {
      t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
    }
  })
}

func contains(s, substr string) bool {
  return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
    (len(s) > len(substr) && (s[:len(substr)] == substr || 
    s[len(s)-len(substr):] == substr || 
    findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
  for i := 0; i <= len(s)-len(substr); i++ {
    if s[i:i+len(substr)] == substr {
      return true
    }
  }
  return false
}
