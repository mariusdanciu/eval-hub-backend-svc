package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/config"
	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/logging"
	"github.ibm.com/julpayne/eval-hub-backend-svc/internal/server"
)

func main() {
	// TODO write fatal errors to the error file and close down the server

	// Create logger once for all requests
	logger, err := logging.NewLogger()
	if err != nil {
		log.Fatal("Failed to create service logger:", err)
	}

	serviceConfig, err := config.LoadConfig(nil)
	if err != nil {
		log.Fatal("Failed to create service config:", err)
	}

	srv, err := server.NewServer(logger, serviceConfig)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

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
