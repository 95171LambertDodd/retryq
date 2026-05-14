package queue

import (
	"sync"
	"time"
)

// RateLimiter controls how many jobs can be dispatched per second.
type RateLimiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
}

// RateLimitConfig holds configuration for the rate limiter.
type RateLimitConfig struct {
	// MaxTokens is the burst capacity (maximum tokens at once).
	MaxTokens float64
	// Rate is the number of tokens replenished per second.
	Rate float64
}

// DefaultRateLimitConfig returns a sensible default rate limit.
var DefaultRateLimitConfig = RateLimitConfig{
	MaxTokens: 10,
	Rate:      5,
}

// NewRateLimiter creates a new token-bucket rate limiter.
func NewRateLimiter(cfg RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		tokens:   cfg.MaxTokens,
		max:      cfg.MaxTokens,
		rate:     cfg.Rate,
		lastTick: time.Now(),
	}
}

// Allow returns true if a token is available, consuming one if so.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastTick).Seconds()
	r.lastTick = now

	r.tokens += elapsed * r.rate
	if r.tokens > r.max {
		r.tokens = r.max
	}

	if r.tokens >= 1 {
		r.tokens--
		return true
	}
	return false
}

// Available returns the current token count (approximate).
func (r *RateLimiter) Available() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tokens
}
