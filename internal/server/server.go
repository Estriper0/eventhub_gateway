package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Estriper0/eventhub_gateway/internal/config"
	"github.com/Estriper0/eventhub_gateway/internal/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
	config     *config.Config
}

func New(logger *slog.Logger, config *config.Config) *Server {
	if config.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	eventHandlers := handlers.NewEvent(logger, config)
	authHandlers := handlers.NewAuth(logger, config)
	
	SetupRoutes(router, eventHandlers, authHandlers, logger, config)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	return &Server{
		httpServer: server,
		logger:     logger,
		config:     config,
	}
}

func (s *Server) Run() {
	s.logger.Info(fmt.Sprintf("Starting server on %s", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.config.Timeout))
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	s.logger.Info("Server shutdown gracefully")
	return nil
}
