package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/nicikess/out-run-management-service/internal/config"
	"github.com/nicikess/out-run-management-service/internal/ports/http/handlers"
	custommiddleware "github.com/nicikess/out-run-management-service/internal/ports/http/middleware"
	"github.com/nicikess/out-run-management-service/internal/service/run"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServer(cfg config.HTTPConfig, runService *run.Service) *Server {
	logger, _ := zap.NewProduction()

	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(custommiddleware.Logger(logger))
	router.Use(custommiddleware.Auth)

	// Initialize handlers
	runHandler := handlers.NewRunHandler(runService, logger)

	// Routes
	router.Route("/api/v1", func(r chi.Router) {
		// Runs
		r.Post("/runs", runHandler.StartRun)
		r.Get("/runs/{runId}", runHandler.GetRun)
		r.Get("/runs/active", runHandler.GetActiveRun)
		r.Put("/runs/{runId}/pause", runHandler.PauseRun)
		r.Put("/runs/{runId}/resume", runHandler.ResumeRun)
		r.Put("/runs/{runId}/end", runHandler.EndRun)
	})

	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%s", cfg.Port),
			Handler:      router,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger: logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", zap.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
