package node

import "fmt"

type SimpleNodeHooks struct {
	logger Logger
}

// SimpleNodeHooks provides hooks for pre-processing and post-processing logic.
// It allows for custom actions, such as logging or validation, to be executed
// before and after a node processes a signal, giving more control over the node's lifecycle.
func NewSimpleNodeHooks(logger Logger) *SimpleNodeHooks {
	return &SimpleNodeHooks{logger: logger}
}

func (h *SimpleNodeHooks) BeforeAction(signal Signal) Result {
	h.logger.Log(fmt.Sprintf("Before processing Node %s", signal.NodeID))
	return Result{Value: "Pre-processing complete", Error: nil}
}

func (h *SimpleNodeHooks) AfterAction(signal Signal) Result {
	h.logger.Log(fmt.Sprintf("After processing Node %s", signal.NodeID))
	return Result{Value: "Post-processing complete", Error: nil}
}
