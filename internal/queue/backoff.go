package queue

import (
	"math"
	"math/rand"
	"time"
)

// BackoffConfig holds parameters for exponential backoff calculation.
type BackoffConfig struct {
	// BaseDelay is the initial delay before the first retry.
	BaseDelay time.Duration
	// MaxDelay caps the computed delay regardless of attempt count.
	MaxDelay time.Duration
	// Multiplier is the exponential growth factor (e.g. 2.0).
	Multiplier float64
	// Jitter adds randomness to avoid thundering-herd problems.
	Jitter bool
}

// DefaultBackoffConfig returns a sensible default configuration.
func DefaultBackoffConfig() BackoffConfig {
	return BackoffConfig{
		BaseDelay:  1 * time.Second,
		MaxDelay:   5 * time.Minute,
		Multiplier: 2.0,
		Jitter:     true,
	}
}

// Next computes the delay duration for the given attempt number (1-indexed).
func (c BackoffConfig) Next(attempt int) time.Duration {
	delay := float64(c.BaseDelay) * math.Pow(c.Multiplier, float64(attempt-1))
	if delay > float64(c.MaxDelay) {
		delay = float64(c.MaxDelay)
	}
	if c.Jitter {
		// Add up to ±20% jitter.
		jitter := delay * 0.2 * (rand.Float64()*2 - 1)
		delay += jitter
		if delay < 0 {
			delay = 0
		}
	}
	return time.Duration(delay)
}

// Schedule sets NextRetryAt on the job based on the current attempt count.
func Schedule(j *Job, cfg BackoffConfig) {
	j.NextRetryAt = time.Now().UTC().Add(cfg.Next(j.Attempts))
}
