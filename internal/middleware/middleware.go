package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		request_id, _ := c.Get("RequestID")
		logger.Info(
			"Start of request",
			slog.Any("request_id", request_id),
			slog.String("url", c.FullPath()),
			slog.String("method", c.Request.Method),
		)
		t := time.Now()
		c.Next()
		status := c.Writer.Status()
		if status >= 500 {
			logger.Error(
				"End of request with server error",
				slog.Any("request_id", request_id),
				slog.String("uri", c.FullPath()),
				slog.Float64("time", time.Since(t).Seconds()),
				slog.Int("status", status),
			)
			return
		} else if status >= 400 {
			logger.Warn(
				"End of request with user error",
				slog.Any("request_id", request_id),
				slog.String("uri", c.FullPath()),
				slog.Float64("time", time.Since(t).Seconds()),
				slog.Int("status", status),
			)
			return
		}
		logger.Info(
			"End of request",
			slog.Any("request_id", request_id),
			slog.String("uri", c.FullPath()),
			slog.Float64("time", time.Since(t).Seconds()),
			slog.Int("status", status),
		)
	}
}

func UUIDMiddlewate() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Something went wrong!",
				})
			}
		}()
		c.Next()
	}
}
