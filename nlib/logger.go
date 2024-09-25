package nlib

import (
	"log"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.Logger = (*SimpleLogger)(nil)

type SimpleLogger struct {
	l *log.Logger
}

func NewSimpleLogger(l *log.Logger) *SimpleLogger {
	return &SimpleLogger{l: l}
}

func (l *SimpleLogger) Log(msg string) {
	l.l.Println(msg)
}
