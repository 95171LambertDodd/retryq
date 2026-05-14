package queue

import (
	"testing"
	"time"
)

func TestBackoffNext_GrowsExponentially(t *testing.T) {
	cfg := BackoffConfig{
		BaseDelay:  1 * time.Second,
		MaxDelay:   10 * time.Minute,
		Multiplier: 2.0,
		Jitter:     false,
	}

	expected := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
		16 * time.Second,
	}

	for i, want := range expected {
		attempt := i + 1
		got := cfg.Next(attempt)
		if got != want {
			t.Errorf("attempt %d: expected %v, got %v", attempt, want, got)
		}
	}
}

func TestBackoffNext_RespectsMaxDelay(t *testing.T) {
	cfg := BackoffConfig{
		BaseDelay:  1 * time.Second,
		MaxDelay:   5 * time.Second,
		Multiplier: 2.0,
		Jitter:     false,
	}

	for attempt := 1; attempt <= 10; attempt++ {
		got := cfg.Next(attempt)
		if got > cfg.MaxDelay {
			t.Errorf("attempt %d: delay %v exceeds max %v", attempt, got, cfg.MaxDelay)
		}
	}
}

func TestSchedule_SetsNextRetryAt(t *testing.T) {
	job := NewJob("test-1", "POST", "http://example.com", nil, nil, 3)
	job.Attempts = 1

	before := time.Now().UTC()
	cfg := DefaultBackoffConfig()
	cfg.Jitter = false
	Schedule(job, cfg)

	if job.NextRetryAt.Before(before) {
		t.Error("NextRetryAt should be in the future")
	}
	if job.NextRetryAt.After(before.Add(cfg.MaxDelay + time.Second)) {
		t.Error("NextRetryAt is unexpectedly far in the future")
	}
}

func TestJob_IsExhausted(t *testing.T) {
	job := NewJob("test-2", "GET", "http://example.com", nil, nil, 3)

	if job.IsExhausted() {
		t.Error("new job should not be exhausted")
	}

	job.Attempts = 3
	if !job.IsExhausted() {
		t.Error("job with attempts == maxRetries should be exhausted")
	}
}
