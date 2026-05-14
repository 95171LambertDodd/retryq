package queue

import (
	"testing"
)

func TestMetrics_InitialValuesAreZero(t *testing.T) {
	ResetMetrics()
	m := GetMetrics()
	snap := m.Snapshot()

	if snap.Enqueued != 0 || snap.Succeeded != 0 || snap.Failed != 0 ||
		snap.Retried != 0 || snap.DeadLetter != 0 {
		t.Errorf("expected all metrics to be zero after reset, got %+v", snap)
	}
}

func TestMetrics_CountersIncrementCorrectly(t *testing.T) {
	ResetMetrics()
	m := GetMetrics()

	m.Enqueued.Add(3)
	m.Succeeded.Add(2)
	m.Failed.Add(1)
	m.Retried.Add(5)
	m.DeadLetter.Add(1)

	snap := m.Snapshot()

	if snap.Enqueued != 3 {
		t.Errorf("expected Enqueued=3, got %d", snap.Enqueued)
	}
	if snap.Succeeded != 2 {
		t.Errorf("expected Succeeded=2, got %d", snap.Succeeded)
	}
	if snap.Failed != 1 {
		t.Errorf("expected Failed=1, got %d", snap.Failed)
	}
	if snap.Retried != 5 {
		t.Errorf("expected Retried=5, got %d", snap.Retried)
	}
	if snap.DeadLetter != 1 {
		t.Errorf("expected DeadLetter=1, got %d", snap.DeadLetter)
	}
}

func TestMetrics_SnapshotIsIsolated(t *testing.T) {
	ResetMetrics()
	m := GetMetrics()

	snap1 := m.Snapshot()
	m.Enqueued.Add(10)
	snap2 := m.Snapshot()

	if snap1.Enqueued != 0 {
		t.Errorf("snap1 should not be affected by later increments, got %d", snap1.Enqueued)
	}
	if snap2.Enqueued != 10 {
		t.Errorf("snap2 should reflect increment, got %d", snap2.Enqueued)
	}
}

func TestGetMetrics_ReturnsSameInstance(t *testing.T) {
	ResetMetrics()
	a := GetMetrics()
	b := GetMetrics()
	if a != b {
		t.Error("GetMetrics should return the same instance")
	}
}
