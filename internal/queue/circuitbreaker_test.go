package queue

import (
	"testing"
	"time"
)

func TestCircuitBreaker_InitialStateIsClosed(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)
	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", cb.State())
	}
}

func TestCircuitBreaker_OpensAfterMaxFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 10*time.Second)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != StateOpen {
		t.Fatalf("expected StateOpen after %d failures, got %v", 3, cb.State())
	}
}

func TestCircuitBreaker_DoesNotAllowWhenOpen(t *testing.T) {
	cb := NewCircuitBreaker(1, 10*time.Second)
	cb.RecordFailure()
	if cb.Allow() {
		t.Fatal("expected Allow() to return false when circuit is open")
	}
}

func TestCircuitBreaker_TransitionsToHalfOpenAfterCooldown(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	if !cb.Allow() {
		t.Fatal("expected Allow() to return true after cooldown (half-open)")
	}
	if cb.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", cb.State())
	}
}

func TestCircuitBreaker_ClosesOnSuccessFromHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(1, 50*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(60 * time.Millisecond)
	cb.Allow() // transition to half-open
	cb.RecordSuccess()
	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", cb.State())
	}
}

func TestCircuitBreaker_DefaultsAppliedOnInvalidConfig(t *testing.T) {
	cb := NewCircuitBreaker(0, 0)
	if cb.maxFailures != 5 {
		t.Errorf("expected default maxFailures=5, got %d", cb.maxFailures)
	}
	if cb.cooldown != 30*time.Second {
		t.Errorf("expected default cooldown=30s, got %v", cb.cooldown)
	}
}

func TestCircuitBreaker_AllowsWhenClosed(t *testing.T) {
	cb := NewCircuitBreaker(5, time.Second)
	if !cb.Allow() {
		t.Fatal("expected Allow() to return true when circuit is closed")
	}
}
