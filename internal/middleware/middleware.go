package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
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

func UUIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		requestID := uuid.New().String()

		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Something went wrong!",
				})
			}
		}()
		c.Next()
	}
}

func RateLimiterMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(limit, burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func JWTAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["uuid"])
			c.Set("email", claims["email"])
		}

		c.Next()
	}
}
