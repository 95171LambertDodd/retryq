package queue

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorker_ProcessesJobSuccessfully(t *testing.T) {
	var called atomic.Int32
	w := NewWorker(WorkerConfig{
		Handler: func(_ context.Context, _ *Job) error {
			called.Add(1)
			return nil
		},
		Backoff:      DefaultBackoffConfig(),
		PollInterval: 20 * time.Millisecond,
	})

	job := NewJob("test-ok", map[string]any{"key": "value"}, 3)
	w.Enqueue(job)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if called.Load() != 1 {
		t.Errorf("expected handler called once, got %d", called.Load())
	}
}

func TestWorker_RetriesOnFailure(t *testing.T) {
	var called atomic.Int32
	w := NewWorker(WorkerConfig{
		Handler: func(_ context.Context, _ *Job) error {
			called.Add(1)
			return errors.New("transient error")
		},
		Backoff: BackoffConfig{
			InitialDelay: 1 * time.Millisecond,
			Multiplier:   1.0,
			MaxDelay:     5 * time.Millisecond,
		},
		PollInterval: 10 * time.Millisecond,
	})

	job := NewJob("test-retry", map[string]any{}, 3)
	w.Enqueue(job)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	if called.Load() < 2 {
		t.Errorf("expected at least 2 handler calls (retries), got %d", called.Load())
	}
}

func TestWorker_DeadLetterOnExhaustion(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "dead_letters_*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	SetDeadLetterPath(tmpFile.Name())
	t.Cleanup(func() { SetDeadLetterPath("dead_letters.jsonl") })

	w := NewWorker(WorkerConfig{
		Handler: func(_ context.Context, _ *Job) error {
			return errors.New("always fails")
		},
		Backoff: BackoffConfig{
			InitialDelay: 1 * time.Millisecond,
			Multiplier:   1.0,
			MaxDelay:     2 * time.Millisecond,
		},
		PollInterval: 10 * time.Millisecond,
	})

	job := NewJob("test-dead", map[string]any{"x": 1}, 1)
	w.Enqueue(job)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	w.Run(ctx)

	data, readErr := os.ReadFile(tmpFile.Name())
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(data) == 0 {
		t.Error("expected dead-letter entry to be written, but file is empty")
	}
}
