package middleware

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("Panic recovered: %v", err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong!",
				})
			}
		}()
		c.Next()
	}
}
