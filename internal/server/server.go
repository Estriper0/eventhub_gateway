package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Estriper0/EventHub/internal/config"
	"github.com/Estriper0/EventHub/internal/handlers"
	"github.com/Estriper0/EventHub/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	HttpServer    *http.Server
	Router        *gin.Engine
	Logger        *slog.Logger
	Config        *config.Config
	EventHandlers *handlers.Event
}

func New(logger *slog.Logger, config *config.Config, eventHandlers *handlers.Event) *Server {
	switch config.Env {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router,
	}
	return &Server{
		HttpServer:    server,
		Router:        router,
		Logger:        logger,
		Config:        config,
		EventHandlers: eventHandlers,
	}
}

func (s *Server) RegisterRoutes() {
	s.Router.Use(cors.Default())
	s.Router.Use(middleware.Recovery(s.Logger))
	s.Router.Use(middleware.UUIDMiddlewate())
	s.Router.Use(middleware.LoggerMiddleware(s.Logger))

	events := s.Router.Group("events")
	events.GET("/", s.EventHandlers.GetAll)
	events.GET("/:id", s.EventHandlers.GetById)
	events.POST("/", s.EventHandlers.Create)
	events.DELETE("/:id", s.EventHandlers.DeleteById)
	events.PATCH("/", s.EventHandlers.Update)
}

func (s *Server) Run() {
	go func() {
		if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	s.Logger.Info(fmt.Sprintf("Server started on %s", s.HttpServer.Addr))
}

func (s *Server) Shutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.Logger.Info("Received shutdown signal. Initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.Config.Timeout))
	defer cancel()

	if err := s.HttpServer.Shutdown(ctx); err != nil {
		return err
	}

	s.Logger.Info("Server shutdown gracefully")
	return nil
}
