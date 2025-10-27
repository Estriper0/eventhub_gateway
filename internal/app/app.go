package app

import (
	"fmt"
	"log/slog"

	"github.com/Estriper0/EventHub/internal/config"
	"github.com/Estriper0/EventHub/internal/handlers"
	"github.com/Estriper0/EventHub/internal/repositories/database/event_repository"
	"github.com/Estriper0/EventHub/internal/server"
	"github.com/Estriper0/EventHub/internal/service/event_service"
	"github.com/Estriper0/EventHub/pkg/db"
)

type App struct {
	logger *slog.Logger
	config *config.Config
	server *server.Server
}

func New(logger *slog.Logger, config *config.Config) *App {
	db := db.GetDB(&config.DB)

	eventRepo := event_repository.New(db)
	eventService := event_service.New(eventRepo, logger)
	eventHandlers := handlers.NewEvents(logger, config, eventService)

	server := server.New(logger, config, eventHandlers)
	server.RegisterRoutes()

	return &App{
		logger: logger,
		config: config,
		server: server,
	}
}

func (a *App) Run() {
	a.server.Run()

	if err := a.server.Shutdown(); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to shutdown server: %v", err))
	}
}
