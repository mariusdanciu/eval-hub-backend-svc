package handlers

import (
  "encoding/json"
  "net/http"
  "time"
)

type Handlers struct{}

func New() *Handlers {
  return &Handlers{}
}

func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
    "status":    "healthy",
    "timestamp": time.Now().UTC().Format(time.RFC3339),
  })
}

func (h *Handlers) HandleStatus(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]interface{}{
    "service":   "eval-hub-backend-svc",
    "version":   "1.0.0",
    "status":    "running",
    "timestamp": time.Now().UTC().Format(time.RFC3339),
  })
}
