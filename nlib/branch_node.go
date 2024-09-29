package nlib

import (
	"fmt"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure SimpleBranchNode implements the node.BranchNode
var _ node.BranchNode = (*SimpleBranchNode)(nil)

// branchCondition holds a target node and a condition function (condFn).
// When condFn evaluates to true, the signal is sent to the target node.
type branchCondition struct {
	target node.Node        // The node to which the signal will be sent if the condition is true.
	condFn node.ConditionFn // The condition function to evaluate.
}

// SimpleBranchNode represents a branching node that routes signals to different nodes
// based on multiple conditions. It checks each condition in sequence, sending the signal
// to the first node whose condition is met.
type SimpleBranchNode struct {
	EmptyNode                    // Inherits base node functionality.
	conditions []branchCondition // List of conditions and their associated target nodes.
}

// NewSimpleBranchNode creates a new SimpleBranchNode with the given logger, state manager, and name.
// It initializes the node and starts a goroutine to listen for incoming signals.
func NewSimpleBranchNode(l node.Logger, sm node.StateManager, name string) *SimpleBranchNode {
	n := SimpleBranchNode{}
	n.Init(l, sm, name) // Initialize the node with logger, state manager, and name.

	// Goroutine to listen for incoming signals and process them.
	go func() {
		for {
			select {
			case sig := <-n.inCh: // Receive a signal from the input channel.
				n.processSignal(sig) // Process the signal.
			case <-n.DoneCh(): // If the done channel is closed, exit the loop.
				return
			}
		}
	}()

	return &n
}

// processSignal processes the incoming signal, checking conditions in sequence.
// If a condition is met, the signal is sent to the corresponding target node.
// If no conditions are met, the signal is sent to the connected nodes.
func (n *SimpleBranchNode) processSignal(sig node.Signal) {
	sig = n.PreProcessSignal(sig) // Run pre-processing hooks on the signal.

	// No specific processing here, but post-processing is handled next.
	sig = n.PostProcesSignal(sig)
	sig.Status = StatusInProcess

	// Iterate over the conditions to find a match.
	for _, cond := range n.conditions {
		// If the condition function evaluates to true, send the signal to the target node.
		if cond.condFn(sig) {
			sig = PrepareSignalForNext(sig)
			n.LogInfo(fmt.Sprintf("Sending to %s", cond.target.ID())) // Log the routing action.
			sig.NodeID = cond.target.ID()
			cond.target.InputCh() <- sig // Send the signal to the target node.
			return
		}
	}
	sig.Status = StatusSuccess

	// If no conditions are met, send the signal to the next connected node.
	n.SendToConnected(sig)
}

// AddConditional adds a new condition and target node to the branch.
// If the condition function evaluates to true, the signal is sent to the target node.
func (n *SimpleBranchNode) AddConditional(target node.Node, fn node.ConditionFn) {
	// Ensure both the target node and condition function are non-nil.
	if target == nil || fn == nil {
		return
	}

	// Create a new branchCondition and add it to the list of conditions.
	bc := branchCondition{target: target, condFn: fn}
	n.conditions = append(n.conditions, bc)
}
