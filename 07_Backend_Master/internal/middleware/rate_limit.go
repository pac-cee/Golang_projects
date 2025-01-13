package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu      sync.RWMutex
	rate    rate.Limit
	burst   int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:    r,
		burst:   b,
	}
}

// getVisitor retrieves or creates a limiter for a visitor
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = limiter
	}

	return limiter
}

// cleanupVisitors removes old visitors
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Hour)
		rl.mu.Lock()
		for ip := range rl.visitors {
			delete(rl.visitors, ip)
		}
		rl.mu.Unlock()
	}
}

// RateLimit middleware implements rate limiting per IP
func RateLimit(requests int, duration time.Duration) gin.HandlerFunc {
	// Convert requests per duration to rate.Limit
	r := rate.Limit(float64(requests) / duration.Seconds())
	limiter := NewRateLimiter(r, requests)

	// Start cleanup routine
	go limiter.cleanupVisitors()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		visitor := limiter.getVisitor(ip)
		if !visitor.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}

// RateLimitByKey middleware implements rate limiting by a specific key
func RateLimitByKey(key string, requests int, duration time.Duration) gin.HandlerFunc {
	r := rate.Limit(float64(requests) / duration.Seconds())
	limiter := NewRateLimiter(r, requests)

	go limiter.cleanupVisitors()

	return func(c *gin.Context) {
		value := c.GetString(key)
		if value == "" {
			value = c.ClientIP() // Fallback to IP if key not found
		}

		visitor := limiter.getVisitor(value)
		if !visitor.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
