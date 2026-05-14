package queue

import (
	"testing"
	"time"
)

func TestRateLimiter_AllowsUpToBurst(t *testing.T) {
	cfg := RateLimitConfig{MaxTokens: 3, Rate: 0}
	rl := NewRateLimiter(cfg)

	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Fatalf("expected Allow() to return true on call %d", i+1)
		}
	}

	if rl.Allow() {
		t.Fatal("expected Allow() to return false after burst exhausted")
	}
}

func TestRateLimiter_RefillsOverTime(t *testing.T) {
	cfg := RateLimitConfig{MaxTokens: 2, Rate: 100} // 100 tokens/sec
	rl := NewRateLimiter(cfg)

	// Drain all tokens
	rl.Allow()
	rl.Allow()

	if rl.Allow() {
		t.Fatal("expected no tokens immediately after drain")
	}

	// Wait enough for at least one token to refill
	time.Sleep(20 * time.Millisecond)

	if !rl.Allow() {
		t.Fatal("expected at least one token to refill after wait")
	}
}

func TestRateLimiter_TokensDoNotExceedMax(t *testing.T) {
	cfg := RateLimitConfig{MaxTokens: 5, Rate: 1000}
	rl := NewRateLimiter(cfg)

	time.Sleep(20 * time.Millisecond)

	avail := rl.Available()
	if avail > cfg.MaxTokens {
		t.Fatalf("tokens %.2f exceeded max %.2f", avail, cfg.MaxTokens)
	}
}

func TestNewRateLimiter_DefaultConfig(t *testing.T) {
	rl := NewRateLimiter(DefaultRateLimitConfig)
	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
	if rl.max != DefaultRateLimitConfig.MaxTokens {
		t.Fatalf("expected max tokens %.2f, got %.2f", DefaultRateLimitConfig.MaxTokens, rl.max)
	}
	if rl.rate != DefaultRateLimitConfig.Rate {
		t.Fatalf("expected rate %.2f, got %.2f", DefaultRateLimitConfig.Rate, rl.rate)
	}
}

func TestRateLimiter_Available_DecreasesOnAllow(t *testing.T) {
	cfg := RateLimitConfig{MaxTokens: 5, Rate: 0}
	rl := NewRateLimiter(cfg)

	before := rl.Available()
	rl.Allow()
	after := rl.Available()

	if after >= before {
		t.Fatalf("expected available tokens to decrease: before=%.2f after=%.2f", before, after)
	}
}
