package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		request_id, _ := c.Get("RequestID")
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
