package nlib

import (
	"log"

	"github.com/dshills/wiggle/node"
)

// Ensure that SimpleLogger implements the node.Logger interface
var _ node.Logger = (*SimpleLogger)(nil)
var _ node.Logger = (*NoLogger)(nil)

// SimpleLogger is a basic implementation of the node.Logger interface.
// It wraps the standard library's log.Logger to provide logging functionality for nodes.
type SimpleLogger struct {
	l *log.Logger // The standard library logger used for logging messages
}

// NewSimpleLogger creates a new instance of SimpleLogger with the provided log.Logger.
// This is a constructor function that returns a pointer to the newly created SimpleLogger.
func NewSimpleLogger(l *log.Logger) *SimpleLogger {
	return &SimpleLogger{l: l}
}

// Log logs the provided message using the underlying log.Logger.
// It simply prints the message with a newline to the standard log output.
func (l *SimpleLogger) Log(msg string) {
	l.l.Println(msg) // Log the message using the standard logger
}

type NoLogger struct{}

func NewNoLogger() *NoLogger {
	return &NoLogger{}
}

func (l *NoLogger) Log(_ string) {
	// NOOP
}
