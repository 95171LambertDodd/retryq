package queue

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func alwaysOKHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func alwaysErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
}

func newTestWorker() *Worker {
	cfg := DefaultBackoffConfig()
	return NewWorker(cfg, func(j *Job) error { return nil })
}

func TestMiddleware_PassesThroughOnSuccess(t *testing.T) {
	worker := newTestWorker()
	mw := NewRetryMiddleware(alwaysOKHandler(), worker)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestMiddleware_EnqueuesJobOnServerError(t *testing.T) {
	worker := newTestWorker()
	mw := NewRetryMiddleware(alwaysErrorHandler(), worker)

	req := httptest.NewRequest(http.MethodPost, "/fail", bytes.NewBufferString(`{"key":"val"}`))
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}

	// Allow worker goroutine to pick up the job.
	time.Sleep(50 * time.Millisecond)
}

func TestMiddleware_WithCustomBackoff(t *testing.T) {
	worker := newTestWorker()
	customCfg := BackoffConfig{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   1.5,
		MaxRetries:   2,
	}
	mw := NewRetryMiddleware(alwaysErrorHandler(), worker, WithBackoffConfig(customCfg))

	req := httptest.NewRequest(http.MethodGet, "/retry", nil)
	rec := httptest.NewRecorder()
	mw.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestStatusRecorder_DefaultsTo200(t *testing.T) {
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}
	if sr.status != http.StatusOK {
		t.Fatalf("expected default status 200, got %d", sr.status)
	}
}
