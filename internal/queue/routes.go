package queue

import "net/http"

// RegisterRoutes attaches the retryq management endpoints to the given ServeMux.
// The following routes are registered:
//
//	GET  /retryq/health    — liveness and queue depth check
//	GET  /retryq/metrics   — snapshot of retry/success/dead-letter counters
//	GET  /retryq/config    — current middleware backoff configuration
func RegisterRoutes(mux *http.ServeMux, w *Worker, m *RetryMiddleware) {
	mux.Handle("/retryq/health", HealthHandler(w))
	mux.Handle("/retryq/metrics", MetricsHandler())
	mux.Handle("/retryq/config", MiddlewareConfigHandler(m))
}
