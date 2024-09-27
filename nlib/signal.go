package nlib

import "github.com/dshills/wiggle/node"

func NewDefaultSignal(targetNode node.Node, task string) node.Signal {
	// Create Context Manager
	contextMgr := NewSimpleContextManager()

	// Create History Manager
	historyMgr := NewSimpleHistoryManager()

	data := NewStringData(task)

	return node.NewSignal(targetNode.ID(), contextMgr, historyMgr, data)
}
