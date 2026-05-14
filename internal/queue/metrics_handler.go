package queue

import (
	"encoding/json"
	"net/http"
)

// metricsResponse is the JSON shape returned by the metrics HTTP handler.
type metricsResponse struct {
	Enqueued   int64 `json:"enqueued"`
	Succeeded  int64 `json:"succeeded"`
	Failed     int64 `json:"failed"`
	Retried    int64 `json:"retried"`
	DeadLetter int64 `json:"dead_letter"`
}

// MetricsHandler returns an http.HandlerFunc that serves a JSON snapshot
// of the current queue metrics. Intended to be mounted at a debug or
// health endpoint (e.g. /metrics or /_retryq/metrics).
//
// Example:
//
//	http.Handle("/_retryq/metrics", queue.MetricsHandler())
func MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snap := GetMetrics().Snapshot()

		resp := metricsResponse{
			Enqueued:   snap.Enqueued,
			Succeeded:  snap.Succeeded,
			Failed:     snap.Failed,
			Retried:    snap.Retried,
			DeadLetter: snap.DeadLetter,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode metrics", http.StatusInternalServerError)
		}
	}
}
