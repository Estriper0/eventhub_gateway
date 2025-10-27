package main

import (
	"github.com/Estriper0/EventHub/internal/app"
	"github.com/Estriper0/EventHub/internal/config"
	"github.com/Estriper0/EventHub/pkg/logger"
)

func main() {
	config := config.New()

	logger := logger.GetLogger(config.Env)

	app := app.New(logger, config)

	app.Run()
}
