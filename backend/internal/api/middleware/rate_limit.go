package middleware

import (
	"sync"
	"time"

	"easy-arbitra/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type clientBucket struct {
	tokens     float64
	lastRefill time.Time
}

func RateLimit(rps float64, burst float64) gin.HandlerFunc {
	if rps <= 0 {
		rps = 20
	}
	if burst <= 0 {
		burst = 40
	}

	var mu sync.Mutex
	clients := make(map[string]*clientBucket)

	return func(c *gin.Context) {
		key := c.ClientIP()
		now := time.Now()

		mu.Lock()
		bucket, ok := clients[key]
		if !ok {
			bucket = &clientBucket{tokens: burst, lastRefill: now}
			clients[key] = bucket
		}

		elapsed := now.Sub(bucket.lastRefill).Seconds()
		bucket.tokens += elapsed * rps
		if bucket.tokens > burst {
			bucket.tokens = burst
		}
		bucket.lastRefill = now

		if bucket.tokens < 1 {
			mu.Unlock()
			response.TooManyRequests(c, "rate limit exceeded")
			c.Abort()
			return
		}

		bucket.tokens -= 1
		mu.Unlock()

		c.Next()
	}
}
