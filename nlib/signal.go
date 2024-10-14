package nlib

import (
	"fmt"

	"github.com/dshills/wiggle/node"
)

// SignalToLog is a helper function to generate a string
// to use in logging a Signal
func SignalToLog(sig node.Signal) string {
	return fmt.Sprintf("{ NodeID: %s, data: %v, Response: %v, Err: %s, Status: %s }", sig.NodeID, sig.Task, sig.Result, sig.Err, sig.Status)
}

func NewSignalFromSignal(toID, fromID string, sig node.Signal) node.Signal {
	return node.Signal{
		NodeID:     toID,
		FromNodeID: fromID,
		Task:       sig.Result,
		Meta:       sig.Meta,
	}
}
