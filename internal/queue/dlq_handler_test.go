package queue_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/retryq/internal/queue"
)

func writeDLQLines(t *testing.T, path string, entries []queue.DLQEntry) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create dlq file: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			t.Fatalf("encode entry: %v", err)
		}
	}
}

func TestDLQHandler_ReturnsEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dlq.log")
	queue.ExportSetDeadLetterPathDirect(path)

	entries := []queue.DLQEntry{
		{JobID: "abc", URL: "http://example.com", Method: "POST", Attempts: 3, LastError: "timeout", Timestamp: time.Now()},
	}
	writeDLQLines(t, path, entries)

	req := httptest.NewRequest(http.MethodGet, "/dlq", nil)
	rr := httptest.NewRecorder()
	queue.DLQHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result []queue.DLQEntry
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].JobID != "abc" {
		t.Errorf("expected job_id abc, got %s", result[0].JobID)
	}
}

func TestDLQHandler_EmptyWhenFileAbsent(t *testing.T) {
	dir := t.TempDir()
	queue.ExportSetDeadLetterPathDirect(filepath.Join(dir, "nonexistent.log"))

	req := httptest.NewRequest(http.MethodGet, "/dlq", nil)
	rr := httptest.NewRecorder()
	queue.DLQHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var result []queue.DLQEntry
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(result))
	}
}

func TestDLQHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/dlq", nil)
	rr := httptest.NewRecorder()
	queue.DLQHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestDLQHandler_NotConfigured(t *testing.T) {
	queue.ExportSetDeadLetterPathDirect("")

	req := httptest.NewRequest(http.MethodGet, "/dlq", nil)
	rr := httptest.NewRecorder()
	queue.DLQHandler(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rr.Code)
	}
}
