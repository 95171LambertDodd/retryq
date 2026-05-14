package queue

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Handler is a function that processes a job. It should return an error if the
// job should be retried, or nil if the job was processed successfully.
type Handler func(ctx context.Context, job *Job) error

// Worker polls a queue of jobs and processes them using a Handler.
type Worker struct {
	mu       sync.Mutex
	jobs     []*Job
	handler  Handler
	backoff  BackoffConfig
	pollInterval time.Duration
	logger   *slog.Logger
}

// WorkerConfig holds configuration for a Worker.
type WorkerConfig struct {
	Handler      Handler
	Backoff      BackoffConfig
	PollInterval time.Duration
	Logger       *slog.Logger
}

// NewWorker creates a new Worker with the given configuration.
func NewWorker(cfg WorkerConfig) *Worker {
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 2 * time.Second
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	return &Worker{
		handler:      cfg.Handler,
		backoff:      cfg.Backoff,
		pollInterval: cfg.PollInterval,
		logger:       cfg.Logger,
	}
}

// Enqueue adds a job to the worker's internal queue.
func (w *Worker) Enqueue(job *Job) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.jobs = append(w.jobs, job)
}

// Run starts the worker loop, processing jobs until the context is cancelled.
func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker shutting down")
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *Worker) process(ctx context.Context) {
	w.mu.Lock()
	ready := make([]*Job, 0, len(w.jobs))
	remaining := w.jobs[:0]
	now := time.Now()
	for _, job := range w.jobs {
		if !job.NextRetryAt.After(now) {
			ready = append(ready, job)
		} else {
			remaining = append(remaining, job)
		}
	}
	w.jobs = remaining
	w.mu.Unlock()

	for _, job := range ready {
		w.dispatch(ctx, job)
	}
}

func (w *Worker) dispatch(ctx context.Context, job *Job) {
	err := w.handler(ctx, job)
	if err == nil {
		w.logger.Info("job completed", "job_id", job.ID)
		return
	}
	job.Attempts++
	if job.IsExhausted() {
		w.logger.Error("job exhausted, sending to dead-letter", "job_id", job.ID, "error", err)
		writeDeadLetter(job, err, w.logger)
		return
	}
	Schedule(job, w.backoff)
	w.logger.Warn("job failed, rescheduled", "job_id", job.ID, "attempts", job.Attempts, "next_retry", job.NextRetryAt)
	w.Enqueue(job)
}
