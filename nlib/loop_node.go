package nlib

import (
	"context"
	"fmt"
	"time"

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
func NewSimpleLoopNode(start node.Node, condFn node.ConditionFn, l node.Logger, sm node.StateManager, name string) *SimpleLoopNode {
	n := SimpleLoopNode{startNode: start, condFn: condFn}
	n.Init(l, sm, name) // Initialize the node with logger, state manager, and name.

	// Goroutine to listen for incoming signals and process them.
	go func() {
		for {
			select {
			case sig := <-n.InputCh(): // Process signal when received.
				n.processSignal(sig)
			case <-n.DoneCh(): // Exit loop when done.
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
		// If the condition is not met, send the signal back to the start node.
		if n.startNode != nil {
			sig = PrepareSignalForNext(sig)
			sig.NodeID = n.startNode.ID()
			n.LogInfo(fmt.Sprintf("Sending to %s", n.startNode.ID()))
			n.startNode.InputCh() <- sig // Send the signal back to the start node.
		}
	}
	sig.Status = StatusSuccess

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure the context is cancelled once we're done

	if err := n.SendToConnected(ctx, sig); err != nil {
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
