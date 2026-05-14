# retryq

A configurable retry queue middleware for HTTP services with exponential backoff and dead-letter logging.

---

## Installation

```bash
go get github.com/yourusername/retryq
```

---

## Usage

```go
package main

import (
    "net/http"
    "github.com/yourusername/retryq"
)

func main() {
    queue := retryq.New(retryq.Config{
        MaxRetries:  5,
        BaseDelay:   500 * time.Millisecond,
        MaxDelay:    30 * time.Second,
        DeadLetterLog: "dead_letters.log",
    })

    handler := queue.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    http.ListenAndServe(":8080", handler)
}
```

Failed requests are automatically retried with exponential backoff. Requests that exhaust all retry attempts are written to the configured dead-letter log for later inspection or replay.

---

## Configuration

| Field           | Type            | Description                          |
|-----------------|-----------------|--------------------------------------|
| `MaxRetries`    | `int`           | Maximum number of retry attempts     |
| `BaseDelay`     | `time.Duration` | Initial delay between retries        |
| `MaxDelay`      | `time.Duration` | Maximum delay cap for backoff        |
| `DeadLetterLog` | `string`        | File path for dead-letter logging    |

---

## License

This project is licensed under the [MIT License](LICENSE).