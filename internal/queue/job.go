package queue

import (
	"net/http"
	"time"
)

// Status represents the current state of a retry job.
type Status string

const (
	StatusPending  Status = "pending"
	StatusRetrying Status = "retrying"
	StatusFailed   Status = "failed"
	StatusDead     Status = "dead"
)

// Job holds all metadata for a single HTTP request that needs to be retried.
type Job struct {
	ID          string
	Method      string
	URL         string
	Headers     http.Header
	Body        []byte
	MaxRetries  int
	Attempts    int
	Status      Status
	LastError   string
	CreatedAt   time.Time
	NextRetryAt time.Time
}

// NewJob creates a new Job with default values.
func NewJob(id, method, url string, headers http.Header, body []byte, maxRetries int) *Job {
	return &Job{
		ID:          id,
		Method:      method,
		URL:         url,
		Headers:     headers,
		Body:        body,
		MaxRetries:  maxRetries,
		Attempts:    0,
		Status:      StatusPending,
		CreatedAt:   time.Now().UTC(),
		NextRetryAt: time.Now().UTC(),
	}
}

// IsExhausted returns true when the job has exceeded its maximum retry count.
func (j *Job) IsExhausted() bool {
	return j.Attempts >= j.MaxRetries
}
