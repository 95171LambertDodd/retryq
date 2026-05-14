package queue

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthStatus represents the current health state of the retry queue.
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	QueueSize int       `json:"queue_size"`
	Workers   int       `json:"workers"`
	Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

// HealthHandler returns an HTTP handler that reports the health of the retry queue.
// It responds with a JSON payload containing queue depth, worker count, and uptime.
func HealthHandler(w *Worker) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status := HealthStatus{
			Status:    "ok",
			Timestamp: time.Now().UTC(),
			QueueSize: w.QueueSize(),
			Workers:   w.WorkerCount(),
			Uptime:    time.Since(startTime).Round(time.Second).String(),
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(rw).Encode(status)
	}
}
