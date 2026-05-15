package queue

import (
	"encoding/json"
	"net/http"
	"sync"
)

var (
	globalCB   *CircuitBreaker
	globalCBMu sync.RWMutex
)

// SetCircuitBreaker sets the global circuit breaker used by the handler.
func SetCircuitBreaker(cb *CircuitBreaker) {
	globalCBMu.Lock()
	defer globalCBMu.Unlock()
	globalCB = cb
}

// CircuitBreakerHandler returns an HTTP handler that exposes the current
// state of the circuit breaker as a JSON response.
//
// GET /circuit-breaker
//
//	{
//	  "state":    "closed",
//	  "failures": 0,
//	  "max_failures": 5
//	}
func CircuitBreakerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		globalCBMu.RLock()
		cb := globalCB
		globalCBMu.RUnlock()

		if cb == nil {
			http.Error(w, "circuit breaker not configured", http.StatusServiceUnavailable)
			return
		}

		cb.mu.Lock()
		payload := map[string]interface{}{
			"state":        cb.state.String(),
			"failures":     cb.failures,
			"max_failures": cb.maxFailures,
		}
		cb.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(payload)
	})
}
