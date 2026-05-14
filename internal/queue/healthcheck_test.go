package queue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_ReturnsOK(t *testing.T) {
	w := newTestWorker()
	handler := HealthHandler(w)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var status HealthStatus
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if status.Status != "ok" {
		t.Errorf("expected status 'ok', got %q", status.Status)
	}
	if status.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if status.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
}

func TestHealthHandler_MethodNotAllowed(t *testing.T) {
	w := newTestWorker()
	handler := HealthHandler(w)

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHealthHandler_ContentType(t *testing.T) {
	w := newTestWorker()
	handler := HealthHandler(w)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestHealthHandler_QueueSizeAndWorkers(t *testing.T) {
	w := newTestWorker()
	handler := HealthHandler(w)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	var status HealthStatus
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if status.QueueSize < 0 {
		t.Errorf("queue size should be non-negative, got %d", status.QueueSize)
	}
	if status.Workers <= 0 {
		t.Errorf("worker count should be positive, got %d", status.Workers)
	}
}
