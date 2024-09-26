package nlib

import "github.com/dshills/wiggle/node"

func NewDefaultSignal(firstNode node.Node, task string) node.Signal {
	// Create Context Manager
	contextMgr := NewSimpleContextManager()

	// Create History Manager
	historyMgr := NewSimpleHistoryManager()

	data := NewStringData(task)

	return node.NewSignal(firstNode.ID(), contextMgr, historyMgr, data)
}
