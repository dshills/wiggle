package nlib

import (
	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.Hooks = (*SimpleNodeHooks)(nil)

type SimpleNodeHooks struct {
	before node.HookFn
	after  node.HookFn
}

// SimpleNodeHooks provides hooks for pre-processing and post-processing logic.
// It allows for custom actions, such as logging or validation, to be executed
// before and after a node processes a signal, giving more control over the node's lifecycle.
func NewSimpleNodeHooks(before, after node.HookFn) *SimpleNodeHooks {
	return &SimpleNodeHooks{before: before, after: after}
}

func (h *SimpleNodeHooks) BeforeAction(signal node.Signal) (node.Signal, error) {
	if h.before != nil {
		return h.before(signal)
	}
	return signal, nil
}

func (h *SimpleNodeHooks) AfterAction(signal node.Signal) (node.Signal, error) {
	if h.after != nil {
		return h.after(signal)
	}
	return signal, nil
}
