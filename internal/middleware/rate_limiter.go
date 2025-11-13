package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/Estriper0/eventhub_gateway/internal/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter *rate.Limiter
}

var clients = make(map[string]*Client)
var mu sync.Mutex

func getClientLimiter(ip string, config *config.Config) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if client, ok := clients[ip]; ok {
		return client.limiter
	}

	limiter := rate.NewLimiter(rate.Every(time.Minute), config.RequestPerMinute)
	clients[ip] = &Client{limiter: limiter}
	return limiter
}

func RateLimiterMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getClientLimiter(c.Request.RemoteAddr, config)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
