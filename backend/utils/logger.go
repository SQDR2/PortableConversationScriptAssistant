package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

var logDir string

// InitLogDir creates the logs/ directory under the current working directory.
// If creation fails, it prints to stderr but does not prevent app startup.
func InitLogDir() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[crash-logging] failed to get working directory: %v\n", err)
		return
	}

	logDir = filepath.Join(cwd, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "[crash-logging] failed to create log directory %s: %v\n", logDir, err)
		logDir = "" // disable file logging
	}
}

// LogError writes an error event to an independent log file in the logs/ directory.
// Each call creates a new file named {YYYY-MM-DD_HH-mm-ss}_{errorType}.log.
// The file contains: timestamp, error type, error message, and full stack trace.
func LogError(errorType string, err interface{}) {
	if logDir == "" {
		// Log directory not initialized or creation failed, fallback to stderr
		fmt.Fprintf(os.Stderr, "[crash-logging] %s: %v\n%s\n", errorType, err, debug.Stack())
		return
	}

	now := time.Now()
	filename := fmt.Sprintf("%s_%s.log", now.Format("2006-01-02_15-04-05"), errorType)
	filePath := filepath.Join(logDir, filename)

	content := fmt.Sprintf(
		"Timestamp: %s\nError Type: %s\nError: %v\n\nStack Trace:\n%s\n",
		now.Format("2006-01-02 15:04:05.000"),
		errorType,
		err,
		debug.Stack(),
	)

	if writeErr := os.WriteFile(filePath, []byte(content), 0644); writeErr != nil {
		// Last resort: print to stderr
		fmt.Fprintf(os.Stderr, "[crash-logging] failed to write log file %s: %v\nOriginal error: %s: %v\n", filePath, writeErr, errorType, err)
	}
}
