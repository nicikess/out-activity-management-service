package main

import (
	"context"
	"github.com/nicikess/out-run-management-service/internal/config"
	"github.com/nicikess/out-run-management-service/internal/ports/repository/mongodb"
	"github.com/nicikess/out-run-management-service/internal/service/run"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize MongoDB repository
	repo, err := mongodb.NewRepository(ctx, cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB repository: %v", err)
	}

	// Initialize run service
	runService := run.NewService(repo)

	// Initialize HTTP server
	server := http.NewServer(cfg.HTTP, runService)

	// Start server
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Println("Received shutdown signal")
	case <-ctx.Done():
		log.Println("Shutting down due to error")
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
