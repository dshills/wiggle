package nlib

import (
	"context"
	"fmt"
	"time"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure SimpleBranchNode implements the node.BranchNode
var _ node.BranchNode = (*SimpleBranchNode)(nil)

// SimpleBranchNode represents a branching node that routes signals to different nodes
// based on multiple conditions. It checks each condition in sequence, sending the signal
// to the first node whose condition is met.
type SimpleBranchNode struct {
	EmptyNode                         // Inherits base node functionality.
	conditions []node.BranchCondition // List of conditions and their associated target nodes.
}

// NewSimpleBranchNode creates a new SimpleBranchNode with the given logger, state manager, and name.
// It initializes the node and starts a goroutine to listen for incoming signals.
func NewSimpleBranchNode(mgr node.StateManager, options node.Options) *SimpleBranchNode {
	n := SimpleBranchNode{}
	n.SetOptions(options)
	n.SetStateManager(mgr)
	n.MakeInputCh()

	// Goroutine to listen for incoming signals and process them.
	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.LogInfo("Received Signal")
				n.ProcessSignal(sig)
			case <-n.StateManager().Register():
				n.LogInfo("Received Done")
				return
			}
		}
	}()

	return &n
}

// ProcessSignal processes the incoming signal, checking conditions in sequence.
// If a condition is met, the signal is sent to the corresponding target node.
// If no conditions are met, the signal is sent to the connected nodes.
func (n *SimpleBranchNode) ProcessSignal(sig node.Signal) {
	var err error
	sig, err = n.PreProcessSignal(sig) // Run pre-processing hooks on the signal.
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

	// Iterate over the conditions to find a match.
	for _, cond := range n.conditions {
		if cond.ConditionFn(sig) {
			n.LogInfo(fmt.Sprintf("Sending to %s", cond.Target.ID()))
			newSig := NewSignalFromSignal(cond.Target.ID(), sig)
			cond.Target.InputCh() <- newSig
			return
		}
	}
	sig.Status = StatusSuccess

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure the context is cancelled once we're done
	// If no conditions are met, send the signal to the next connected node.
	if err := n.SendToConnected(ctx, sig); err != nil {
		n.Fail(sig, err)
		return
	}
}

// AddConditional adds a new condition and target node to the branch.
// If the condition function evaluates to true, the signal is sent to the target node.
func (n *SimpleBranchNode) AddConditional(conditions ...node.BranchCondition) {
	// Create a new branchCondition and add it to the list of conditions.
	n.conditions = append(n.conditions, conditions...)
}

func (n *SimpleBranchNode) Conditions() []node.BranchCondition {
	return n.conditions
}
