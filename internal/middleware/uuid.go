package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UUIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		requestID := uuid.New().String()

		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
