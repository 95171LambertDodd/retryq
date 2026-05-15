package queue

// ExportSetDeadLetterPathDirect sets the internal path variable directly for tests
// without touching the file system, allowing DLQHandler tests to control the path.
func ExportSetDeadLetterPathDirect(p string) {
	currentDeadLetterPath = p
}
