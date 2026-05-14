// Package queue provides the core data structures and scheduling logic
// for retryq's HTTP retry queue.
//
// # Overview
//
// A [Job] represents a single outbound HTTP request that failed and needs
// to be retried. Each job tracks its attempt history, current status, and
// the timestamp of its next scheduled retry.
//
// # Backoff
//
// Retry delays are computed by [BackoffConfig.Next] using exponential backoff:
//
//	delay = baseDelay * multiplier^(attempt-1)
//
// Delays are capped at MaxDelay. When Jitter is enabled, up to ±20% random
// noise is added to spread load across retry windows.
//
// Use [Schedule] to automatically update a job's NextRetryAt field after
// each failed attempt.
//
// # Job Lifecycle
//
//	pending → retrying → (success) done
//	                   → (exhausted) dead
//
// Dead jobs are handed off to the dead-letter logger for inspection.
package queue
