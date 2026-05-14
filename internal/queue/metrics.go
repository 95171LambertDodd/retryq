package queue

import (
	"sync/atomic"
)

// Metrics tracks runtime statistics for the retry queue.
type Metrics struct {
	Enqueued   atomic.Int64
	Succeeded  atomic.Int64
	Failed     atomic.Int64
	Retried    atomic.Int64
	DeadLetter atomic.Int64
}

// globalMetrics is the package-level metrics instance.
var globalMetrics = &Metrics{}

// GetMetrics returns the global metrics instance.
func GetMetrics() *Metrics {
	return globalMetrics
}

// ResetMetrics resets all counters to zero. Useful in tests.
func ResetMetrics() {
	globalMetrics = &Metrics{}
}

// Snapshot returns a point-in-time copy of the current metrics.
func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		Enqueued:   m.Enqueued.Load(),
		Succeeded:  m.Succeeded.Load(),
		Failed:     m.Failed.Load(),
		Retried:    m.Retried.Load(),
		DeadLetter: m.DeadLetter.Load(),
	}
}

// MetricsSnapshot is an immutable copy of Metrics at a point in time.
type MetricsSnapshot struct {
	Enqueued   int64
	Succeeded  int64
	Failed     int64
	Retried    int64
	DeadLetter int64
}
