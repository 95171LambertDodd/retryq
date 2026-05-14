package queue

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
	"time"
)

// DeadLetterEntry represents a job that has exhausted all retry attempts.
type DeadLetterEntry struct {
	JobID     string         `json:"job_id"`
	Payload   map[string]any `json:"payload"`
	Attempts  int            `json:"attempts"`
	LastError string         `json:"last_error"`
	FailedAt  time.Time      `json:"failed_at"`
}

var (
	deadLetterMu   sync.Mutex
	deadLetterPath = "dead_letters.jsonl"
)

// SetDeadLetterPath overrides the default dead-letter log file path.
// Must be called before any workers are started.
func SetDeadLetterPath(path string) {
	deadLetterMu.Lock()
	defer deadLetterMu.Unlock()
	deadLetterPath = path
}

// writeDeadLetter appends a dead-letter entry to the configured JSONL file.
func writeDeadLetter(job *Job, err error, logger *slog.Logger) {
	entry := DeadLetterEntry{
		JobID:     job.ID,
		Payload:   job.Payload,
		Attempts:  job.Attempts,
		LastError: err.Error(),
		FailedAt:  time.Now().UTC(),
	}

	data, jsonErr := json.Marshal(entry)
	if jsonErr != nil {
		logger.Error("failed to marshal dead-letter entry", "job_id", job.ID, "error", jsonErr)
		return
	}

	deadLetterMu.Lock()
	defer deadLetterMu.Unlock()

	f, openErr := os.OpenFile(deadLetterPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if openErr != nil {
		logger.Error("failed to open dead-letter file", "path", deadLetterPath, "error", openErr)
		return
	}
	defer f.Close()

	if _, writeErr := f.Write(append(data, '\n')); writeErr != nil {
		logger.Error("failed to write dead-letter entry", "job_id", job.ID, "error", writeErr)
	}
}
