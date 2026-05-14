package queue

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// RetryMiddleware wraps an HTTP handler and enqueues failed requests for retry.
type RetryMiddleware struct {
	worker  *Worker
	next    http.Handler
	backoff BackoffConfig
}

// MiddlewareOption configures a RetryMiddleware.
type MiddlewareOption func(*RetryMiddleware)

// WithBackoffConfig sets a custom backoff configuration.
func WithBackoffConfig(cfg BackoffConfig) MiddlewareOption {
	return func(m *RetryMiddleware) {
		m.backoff = cfg
	}
}

// NewRetryMiddleware creates a RetryMiddleware wrapping the given handler.
func NewRetryMiddleware(next http.Handler, worker *Worker, opts ...MiddlewareOption) *RetryMiddleware {
	m := &RetryMiddleware{
		worker:  worker,
		next:    next,
		backoff: DefaultBackoffConfig(),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// ServeHTTP processes the request and enqueues a retry job on failure.
func (m *RetryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))

	rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
	m.next.ServeHTTP(rec, r)

	if rec.status >= 500 {
		job := NewJob(r.Method, r.URL.String(), body, m.backoff.MaxRetries)
		Schedule(job, m.backoff, time.Now())
		m.worker.Enqueue(job)
	}
}

// statusRecorder captures the HTTP status code written by a handler.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}
