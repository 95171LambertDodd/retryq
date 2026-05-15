package queue

// currentDeadLetterPath holds the active dead-letter log file path.
// It is set via SetDeadLetterPath and read by DLQHandler.
var currentDeadLetterPath string

// SetDeadLetterPath configures the file path used for dead-letter logging.
// The directory must already exist. Passing an empty string disables file logging.
func SetDeadLetterPath(path string) {
	currentDeadLetterPath = path
}
