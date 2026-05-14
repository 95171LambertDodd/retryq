package queue

// QueueSize returns the number of jobs currently pending in the worker's queue.
func (w *Worker) QueueSize() int {
	return len(w.jobs)
}

// WorkerCount returns the number of concurrent worker goroutines configured.
func (w *Worker) WorkerCount() int {
	return w.concurrency
}
