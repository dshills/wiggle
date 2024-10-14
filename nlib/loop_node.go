package nlib

import (
	"context"
	"fmt"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure SimpleLoopNode implements the node.LoopNode interface.
var _ node.LoopNode = (*SimpleLoopNode)(nil)

// SimpleLoopNode represents a loop node in a processing chain. It sends the signal
// back to the start node based on a conditional function (condFn), allowing for
// iterative processing of the signal until the condition is met.
type SimpleLoopNode struct {
	EmptyNode                  // Inherits base node functionality.
	startNode node.Node        // Reference to the start node where signals are sent back for re-processing.
	condFn    node.ConditionFn // Condition function that determines if the loop should continue.
}

// NewSimpleLoopNode creates and initializes a SimpleLoopNode.
// - start: The node to which signals will be sent for re-processing.
// - condFn: The function to evaluate if the loop should continue.
// - l: Logger for logging messages.
// - sm: StateManager for managing node state.
// - name: Name/ID of the node.
func NewSimpleLoopNode(start node.Node, condFn node.ConditionFn, mgr node.StateManager, options node.Options) *SimpleLoopNode {
	n := SimpleLoopNode{startNode: start, condFn: condFn}
	n.SetOptions(options)
	n.SetStateManager(mgr)
	n.MakeInputCh()

	// Goroutine to listen for incoming signals and process them.
	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.LogInfo("Received signal")
				n.processSignal(sig)
			case <-n.StateManager().Register():
				n.LogInfo("Received Done")
				return
			}
		}
	}()

	return &n
}

// processSignal handles the incoming signal, checks the condition, and either sends
// it back to the start node or forwards it to the connected nodes.
func (n *SimpleLoopNode) processSignal(sig node.Signal) {
	var err error

	sig, err = n.PreProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	// No specific processing here, but post-processing is handled next.
	sig, err = n.PostProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}
	sig.Status = StatusInProcess

	// Check if the condition is met (if condFn is not nil).
	if n.condFn == nil || !n.condFn(sig) {
		if n.startNode != nil {
			n.LogInfo(fmt.Sprintf("Sending to %s", n.startNode.ID()))
			newSig := NewSignalFromSignal(n.startNode.ID(), n.ID(), sig)
			n.startNode.InputCh() <- newSig
		}
	}
	sig.Status = StatusSuccess

	if err := n.SendToConnected(context.TODO(), sig); err != nil {
		n.Fail(sig, err)
		return
	}
}

// SetStartNode sets the start node where the signal will be looped back to for re-processing.
func (n *SimpleLoopNode) SetStartNode(start node.Node) {
	n.startNode = start
}

// SetConditionFunc sets the condition function that determines if the signal should be looped back.
func (n *SimpleLoopNode) SetConditionFunc(fn node.ConditionFn) {
	n.condFn = fn
}
