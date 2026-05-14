package queue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddlewareConfigHandler_ReturnsConfig(t *testing.T) {
	cfg := BackoffConfig{
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
		MaxRetries:   5,
	}

	req := httptest.NewRequest(http.MethodGet, "/middleware/config", nil)
	rec := httptest.NewRecorder()

	MiddlewareConfigHandler(cfg)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result MiddlewareConfig
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.MaxRetries != 5 {
		t.Errorf("expected MaxRetries=5, got %d", result.MaxRetries)
	}
	if result.InitialDelayMs != 500 {
		t.Errorf("expected InitialDelayMs=500, got %d", result.InitialDelayMs)
	}
	if result.MaxDelayMs != 30000 {
		t.Errorf("expected MaxDelayMs=30000, got %d", result.MaxDelayMs)
	}
	if result.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", result.Multiplier)
	}
}

func TestMiddlewareConfigHandler_MethodNotAllowed(t *testing.T) {
	cfg := DefaultBackoffConfig()

	req := httptest.NewRequest(http.MethodPost, "/middleware/config", nil)
	rec := httptest.NewRecorder()

	MiddlewareConfigHandler(cfg)(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestMiddlewareConfigHandler_ContentType(t *testing.T) {
	cfg := DefaultBackoffConfig()

	req := httptest.NewRequest(http.MethodGet, "/middleware/config", nil)
	rec := httptest.NewRecorder()

	MiddlewareConfigHandler(cfg)(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}
