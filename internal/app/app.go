package app

import (
	"log/slog"

	"github.com/Estriper0/eventhub_gateway/internal/config"
	"github.com/Estriper0/eventhub_gateway/internal/server"
)

type App struct {
	logger *slog.Logger
	config *config.Config
	server *server.Server
}

func New(logger *slog.Logger, config *config.Config) *App {
	server := server.New(logger, config)

	return &App{
		logger: logger,
		config: config,
		server: server,
	}
}

func (a *App) Run() {
	a.logger.Info("Start application")
	a.server.Run()
}

func (a *App) Stop() {
	a.server.Stop()
	a.logger.Info("Stop application")
}
