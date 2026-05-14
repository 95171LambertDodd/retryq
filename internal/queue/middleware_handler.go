package queue

import (
	"encoding/json"
	"net/http"
)

// MiddlewareConfig holds JSON-serialisable configuration for the retry middleware.
type MiddlewareConfig struct {
	MaxRetries     int     `json:"max_retries"`
	InitialDelayMs int64   `json:"initial_delay_ms"`
	MaxDelayMs     int64   `json:"max_delay_ms"`
	Multiplier     float64 `json:"multiplier"`
}

// MiddlewareConfigHandler returns an HTTP handler that exposes the current
// middleware backoff configuration as JSON.
func MiddlewareConfigHandler(cfg BackoffConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := MiddlewareConfig{
			MaxRetries:     cfg.MaxRetries,
			InitialDelayMs: cfg.InitialDelay.Milliseconds(),
			MaxDelayMs:     cfg.MaxDelay.Milliseconds(),
			Multiplier:     cfg.Multiplier,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode config", http.StatusInternalServerError)
		}
	}
}
