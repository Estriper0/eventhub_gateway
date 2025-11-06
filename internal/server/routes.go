package server

import (
	"log/slog"

	"github.com/Estriper0/EventHub/internal/config"
	"github.com/Estriper0/EventHub/internal/handlers"
	"github.com/Estriper0/EventHub/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, eventHandlers *handlers.Event, authHandlers *handlers.Auth, logger *slog.Logger, config *config.Config) {
	r.Use(cors.Default())
	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.RateLimiterMiddleware(config))
	r.Use(middleware.UUIDMiddleware())
	r.Use(middleware.LoggerMiddleware(logger))

	events := r.Group("events")
	events.Use(middleware.JWTAuthMiddleware(config.AccessTokenSecret))
	events.GET("/", eventHandlers.GetAll)
	events.GET("/status/:status", eventHandlers.GetAllByStatus)
	events.GET("/creator/:creator", eventHandlers.GetAllByCreator)
	events.GET("/:id/users", eventHandlers.GetAllUsersByEvent)
	events.GET("/:id", eventHandlers.GetById)
	events.POST("/", eventHandlers.Create)
	events.DELETE("/:id", eventHandlers.DeleteById)
	events.PUT("/", eventHandlers.Update)
	events.GET("/me", eventHandlers.GetAllByUser)
	events.POST("/:id/register", eventHandlers.Register)
	events.POST("/:id/cancell", eventHandlers.CancellRegister)

	auth := r.Group("auth")
	auth.POST("/register", authHandlers.Register)
	auth.POST("/login", authHandlers.Login)
	auth.POST("/admin", authHandlers.IsAdmin)
	auth.POST("/refresh", authHandlers.Refresh)
	auth.POST("/logout", authHandlers.Logout)
}
