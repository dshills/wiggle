package nlib

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dshills/wiggle/node"
)

const (
	StatusSuccess   = "success"
	StatusInProcess = "in-process"
	StatusFail      = "fail"
)

// Compile-time check that EmptyNode implements the node.Node interface
var _ node.Node = (*EmptyNode)(nil)

// EmptyNode is a boilerplate implementation of the node.Node interface
type EmptyNode struct {
	nodes    []node.Node
	id       string
	guide    node.Guidance
	hooks    node.Hooks
	stateMgr node.StateManager
	errGuide node.ErrorGuidance
	mu       sync.RWMutex
	inputCh  chan node.Signal
}

// Connect attaches nodes to the EmptyNode
func (n *EmptyNode) Connect(nn ...node.Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.nodes = append(n.nodes, nn...)
}

// ID returns the ID of the EmptyNode
func (n *EmptyNode) ID() string {
	return n.id
}

// SetID sets the ID for the EmptyNode
func (n *EmptyNode) SetID(id string) {
	n.id = id
}

// InputCh returns the input channel for the EmptyNode
func (n *EmptyNode) InputCh() chan node.Signal {
	return n.inputCh
}

func (n *EmptyNode) SetInputCh(inCh chan node.Signal) {
	n.inputCh = inCh
}

func (n *EmptyNode) SetOptions(options node.Options) {
	n.hooks = options.Hooks
	n.guide = options.Guidance
	n.errGuide = options.ErrorGuidance
	n.id = options.ID
	if n.id == "" {
		var err error
		n.id, err = GenerateUUID()
		if err != nil {
			n.LogErr(err)
			n.id = "UUID-FAILED"
		}
	}
}

// SetStateManager sets the StateManager for managing the state of the node
func (n *EmptyNode) SetStateManager(mgr node.StateManager) {
	n.stateMgr = mgr
}

// Helper functions

func (n *EmptyNode) Guidance() node.Guidance {
	return n.guide
}

// Hooks returns the Hooks associated with the EmptyNode
func (n *EmptyNode) Hooks() node.Hooks {
	return n.hooks
}

// StateManager returns the StateManager associated with the EmptyNode
func (n *EmptyNode) StateManager() node.StateManager {
	return n.stateMgr
}

func (n *EmptyNode) MakeInputCh() {
	n.inputCh = make(chan node.Signal)
}

// Return connected nodes
func (n *EmptyNode) Nodes() []node.Node {
	return n.nodes
}

// RunBeforeHook executes the before-action hooks for the signal
func (n *EmptyNode) RunBeforeHook(sig node.Signal) (node.Signal, error) {
	if n.hooks != nil {
		return n.hooks.BeforeAction(sig)
	}
	return sig, nil
}

// RunAfterHook executes the after-action hooks for the signal
func (n *EmptyNode) RunAfterHook(sig node.Signal) (node.Signal, error) {
	if n.hooks != nil {
		return n.hooks.AfterAction(sig)
	}
	return sig, nil
}

// LogErr logs an error message using the logger
func (n *EmptyNode) LogErr(err error) {
	n.log("error", n.id, err.Error())
}

// LogInfo logs an informational message using the logger
func (n *EmptyNode) LogInfo(msg string) {
	n.log("info", n.id, msg)
}

// LogDebug logs a debug message using the logger
func (n *EmptyNode) LogDebug(msg string) {
	n.log("debug", n.id, msg)
}

// log handles the actual logging of messages with a specified severity
func (n *EmptyNode) log(severity, id, msg string) {
	n.stateMgr.Log(fmt.Sprintf("{ \"severity\": %q, \"id\": %q, \"msg\": %q }", severity, id, msg))
}

func (n *EmptyNode) Fail(sig node.Signal, err error) {
	n.LogErr(err)
	sig.Err = err.Error()
	sig.Status = StatusFail
	n.StateManager().UpdateState(sig)
	n.StateManager().Complete()
}

func (n *EmptyNode) PreProcessSignal(sig node.Signal) (node.Signal, error) {
	if resMgr := n.stateMgr.ResourceManager(); resMgr != nil {
		// Rate limiting check with exponential backoff
		for retries := 0; retries < 3; retries++ {
			if err := resMgr.RateLimit(sig); err == nil {
				break
			}
			time.Sleep(time.Duration(retries*retries) * time.Second) // Exponential backoff
		}
		if err := resMgr.RateLimit(sig); err != nil {
			return sig, fmt.Errorf("exceeded rate limit, could not recover")
		}
	}

	// Run any registered before-action hooks
	return n.RunBeforeHook(sig)
}

func (n *EmptyNode) PostProcessSignal(sig node.Signal) (node.Signal, error) {
	// Run any registered after-action hooks
	sig, err := n.RunAfterHook(sig)
	if err != nil {
		return sig, err
	}

	// Update the state of the signal after processing
	n.stateMgr.UpdateState(sig)

	return sig, nil
}

// SendToConnected sends a signal to all connected nodes using the provided context for timeout control
func (n *EmptyNode) SendToConnected(ctx context.Context, sig node.Signal) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for _, conNode := range n.nodes {
		n.LogInfo(fmt.Sprintf("Sending to %s", conNode.ID()))
		newSig := NewSignalFromSignal(conNode.ID(), n.ID(), sig)

		select {
		case <-ctx.Done():
			fmt.Println("ctx.Done")
			err := fmt.Errorf("context timeout or cancellation while sending signal to node %s: %v", conNode.ID(), ctx.Err())
			n.LogErr(err)
			return err
		case conNode.InputCh() <- newSig:
		}
	}
	return nil
}
