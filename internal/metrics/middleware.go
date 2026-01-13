package metrics

import (
  "net/http"
  "strconv"
  "time"
)

// Middleware wraps an http.Handler to collect Prometheus metrics
func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // Track in-flight requests
    HTTPRequestInFlight.Inc()
    defer HTTPRequestInFlight.Dec()

    // Create a response writer wrapper to capture status code
    rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

    // Call the next handler
    next.ServeHTTP(rw, r)

    // Calculate duration
    duration := time.Since(start).Seconds()

    // Extract method and endpoint
    method := r.Method
    endpoint := r.URL.Path
    status := strconv.Itoa(rw.statusCode)

    // Record metrics
    HTTPRequestDuration.WithLabelValues(method, endpoint, status).Observe(duration)
    HTTPRequestTotal.WithLabelValues(method, endpoint, status).Inc()
  })
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
  http.ResponseWriter
  statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
  rw.statusCode = code
  rw.ResponseWriter.WriteHeader(code)
}
