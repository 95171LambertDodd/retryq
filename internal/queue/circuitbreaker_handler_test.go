package queue

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCircuitBreakerHandler_ReturnsClosedState(t *testing.T) {
	cb := NewCircuitBreaker(3, 5*time.Second)
	SetCircuitBreaker(cb)
	t.Cleanup(func() { SetCircuitBreaker(nil) })

	req := httptest.NewRequest(http.MethodGet, "/circuit-breaker", nil)
	rec := httptest.NewRecorder()
	CircuitBreakerHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["state"] != "closed" {
		t.Errorf("expected state=closed, got %v", body["state"])
	}
}

func TestCircuitBreakerHandler_MethodNotAllowed(t *testing.T) {
	cb := NewCircuitBreaker(3, 5*time.Second)
	SetCircuitBreaker(cb)
	t.Cleanup(func() { SetCircuitBreaker(nil) })

	req := httptest.NewRequest(http.MethodPost, "/circuit-breaker", nil)
	rec := httptest.NewRecorder()
	CircuitBreakerHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestCircuitBreakerHandler_NotConfigured(t *testing.T) {
	SetCircuitBreaker(nil)

	req := httptest.NewRequest(http.MethodGet, "/circuit-breaker", nil)
	rec := httptest.NewRecorder()
	CircuitBreakerHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rec.Code)
	}
}

func TestCircuitBreakerHandler_ContentType(t *testing.T) {
	cb := NewCircuitBreaker(3, 5*time.Second)
	SetCircuitBreaker(cb)
	t.Cleanup(func() { SetCircuitBreaker(nil) })

	req := httptest.NewRequest(http.MethodGet, "/circuit-breaker", nil)
	rec := httptest.NewRecorder()
	CircuitBreakerHandler().ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}
