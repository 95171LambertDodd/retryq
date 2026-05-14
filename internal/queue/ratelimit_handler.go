package queue

import (
	"encoding/json"
	"net/http"
)

// rateLimitState holds a reference to the active rate limiter for inspection.
var rateLimitState *RateLimiter

// SetRateLimiter registers the active rate limiter for the handler to expose.
func SetRateLimiter(rl *RateLimiter) {
	rateLimitState = rl
}

// RateLimitHandler returns the current rate limiter status as JSON.
//
// GET /retryq/ratelimit
//
// Response:
//
//	{
//	  "available_tokens": 4.75,
//	  "max_tokens": 10,
//	  "rate_per_second": 5
//	}
func RateLimitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if rateLimitState == nil {
		http.Error(w, "rate limiter not configured", http.StatusServiceUnavailable)
		return
	}

	payload := map[string]float64{
		"available_tokens": rateLimitState.Available(),
		"max_tokens":       rateLimitState.max,
		"rate_per_second":  rateLimitState.rate,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payload)
}
