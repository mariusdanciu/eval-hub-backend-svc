package main

import (
  "context"
  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"

  "eval-hub-backend-svc/internal/server"
)

func main() {
  srv := server.NewServer()

  // Start server in a goroutine
  go func() {
    if err := srv.Start(); err != nil && err != http.ErrServerClosed {
      log.Fatal("Server failed to start:", err)
    }
  }()

  // Wait for interrupt signal to gracefully shutdown the server
  quit := make(chan os.Signal, 1)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  <-quit

  log.Println("Shutting down server...")

  // Create a context with timeout for graceful shutdown
  ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  defer cancel()

  if err := srv.Shutdown(ctx); err != nil {
    log.Fatal("Server forced to shutdown:", err)
  }

  log.Println("Server exited")
}
