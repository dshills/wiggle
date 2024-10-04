package nlib

import (
	"fmt"

	"github.com/dshills/wiggle/node"
)

// NewDefaultSignal is a simplified version of NewSignal
// It creates a ContextManager and HistoryManager and assumes
// string data
func NewDefaultSignal(targetNode node.Node, task string) node.Signal {
	// Create Context Manager
	contextMgr := NewSimpleContextManager()

	// Create History Manager
	historyMgr := NewSimpleHistoryManager()

	data := NewStringData(task)

	return NewSignal(targetNode.ID(), contextMgr, historyMgr, data)
}

// NewSignal will return a new Signal. This is typically used to generate the
// initial Signal at the start of processing
func NewSignal(id string, cm node.ContextManager, hx node.HistoryManager, task node.DataCarrier, meta ...node.Meta) node.Signal {
	return node.Signal{
		NodeID:  id,
		Context: cm,
		History: hx,
		Meta:    meta,
		Task:    task,
	}
}

// PrepareSignalForNext is generally used when a node is
// preparing to send it's Signal to the next node
// It will store the history of the current Node
// swap the Result into the Task
func PrepareSignalForNext(sig node.Signal) node.Signal {
	sig.AddHistory()
	sig.Task = sig.Result
	sig.Result = nil
	return sig
}

// SignalToLog is a helper function to generate a string
// to use in logging a Signal
func SignalToLog(sig node.Signal) string {
	return fmt.Sprintf("{ NodeID: %s, data: %v, Response: %v, Err: %s, Status: %s }", sig.NodeID, sig.Task, sig.Result, sig.Err, sig.Status)
}

// SignalFromSignal is generally used by a partitioner node to create
// Signals for each child node
func SignalFromSignal(sig node.Signal, task node.DataCarrier) node.Signal {
	newSignal := sig
	newSignal.Task = task
	newSignal.Result = nil
	return newSignal
}
