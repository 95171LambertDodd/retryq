package queue

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DLQEntry represents a single dead-letter queue log entry.
type DLQEntry struct {
	JobID     string            `json:"job_id"`
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Attempts  int               `json:"attempts"`
	LastError string            `json:"last_error"`
	Headers   map[string]string `json:"headers,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// DLQHandler returns an HTTP handler that reads and serves dead-letter log entries.
func DLQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	path := deadLetterPath()
	if path == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "dead-letter log not configured"})
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]DLQEntry{})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to read dead-letter log"})
		return
	}

	var entries []DLQEntry
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var entry DLQEntry
		if err := json.Unmarshal([]byte(line), &entry); err == nil {
			entries = append(entries, entry)
		}
	}
	if entries == nil {
		entries = []DLQEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(entries)
}

// deadLetterPath returns the current dead-letter file path via the package-level accessor.
func deadLetterPath() string {
	// Re-use the path set via SetDeadLetterPath stored in the package.
	return filepath.Clean(currentDeadLetterPath)
}
