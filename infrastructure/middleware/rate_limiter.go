package middleware

import (
	"sync"
	"time"

	"gin-boilerplate/infrastructure/errors"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	rate     int
	interval time.Duration
	clients  map[string]*client
	mutex    sync.RWMutex
}

type client struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		rate:     rate,
		interval: interval,
		clients:  make(map[string]*client),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if the request is allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	c, exists := rl.clients[clientID]

	if !exists {
		rl.clients[clientID] = &client{
			tokens:   rl.rate - 1,
			lastSeen: now,
		}
		return true
	}

	// Calculate tokens to add based on time elapsed
	elapsed := now.Sub(c.lastSeen)
	tokensToAdd := int(elapsed / rl.interval * time.Duration(rl.rate))

	if tokensToAdd > 0 {
		c.tokens += tokensToAdd
		if c.tokens > rl.rate {
			c.tokens = rl.rate
		}
		c.lastSeen = now
	}

	if c.tokens > 0 {
		c.tokens--
		return true
	}

	return false
}

// cleanup removes old clients
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for clientID, c := range rl.clients {
			if now.Sub(c.lastSeen) > time.Hour {
				delete(rl.clients, clientID)
			}
		}
		rl.mutex.Unlock()
	}
}

// RateLimiterMiddleware creates a rate limiting middleware
func RateLimiterMiddleware(rateLimiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := getClientID(c)
		
		if !rateLimiter.Allow(clientID) {
			errors.AbortWithError(c, errors.NewRateLimitError("Rate limit exceeded"))
			return
		}

		c.Next()
	}
}

// getClientID extracts client identifier for rate limiting
func getClientID(c *gin.Context) string {
	// Try to get user ID from context first (for authenticated users)
	if userID, exists := GetUserID(c); exists {
		return "user:" + userID
	}

	// Fall back to IP address
	return "ip:" + c.ClientIP()
}
