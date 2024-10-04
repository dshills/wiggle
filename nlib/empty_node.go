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
	nodes       []node.Node
	id          string
	inCh        chan node.Signal
	guide       node.Guidance
	hooks       node.Hooks
	logger      node.Logger
	resourceMgr node.ResourceManager
	stateMgr    node.StateManager
	doneCh      chan struct{}
	errGuide    node.ErrorGuidance
	mu          sync.RWMutex
	closeOnce   sync.Once
	waitGroup   *sync.WaitGroup
}

// Init initializes the EmptyNode with a logger, state manager, and ID
func (n *EmptyNode) Init(l node.Logger, mgr node.StateManager, id string) {
	n.SetID(id)
	n.SetLogger(l)
	n.SetStateManager(mgr)
	n.MakeInputCh(5)
}

// Close safely closes the input and done channels of the EmptyNode, ensuring it is only done once
func (n *EmptyNode) Close() {
	n.closeOnce.Do(func() {
		n.LogInfo("Closing")
		if n.doneCh != nil {
			close(n.doneCh)
			n.doneCh = nil
		}
		if n.inCh != nil {
			close(n.inCh)
			n.inCh = nil
		}
	})
}

// Connect attaches nodes to the EmptyNode
func (n *EmptyNode) Connect(nn ...node.Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.nodes = append(n.nodes, nn...)
}

// SendToConnected sends a signal to all connected nodes using the provided context for timeout control
func (n *EmptyNode) SendToConnected(ctx context.Context, sig node.Signal) error {
	sig = PrepareSignalForNext(sig)

	n.mu.RLock()
	defer n.mu.RUnlock()

	for _, conNode := range n.nodes {
		n.LogInfo(fmt.Sprintf("Sending to %s", conNode.ID()))
		sig.NodeID = conNode.ID()

		select {
		case <-ctx.Done():
			fmt.Println("ctx.Done")
			err := fmt.Errorf("context timeout or cancellation while sending signal to node %s: %v", conNode.ID(), ctx.Err())
			n.LogErr(err)
			return err
		case conNode.InputCh() <- sig:
		}
	}
	return nil
}

// SetErrorGuidance sets the ErrorGuidance for the EmptyNode
func (n *EmptyNode) SetErrorGuidance(errGuide node.ErrorGuidance) {
	n.errGuide = errGuide
}

// ErrorAction determines the action to take based on the provided error
func (n *EmptyNode) ErrorAction(err error) node.ErrGuide {
	if n.errGuide == nil {
		return node.ErrGuideFail
	}
	return n.errGuide.Action(err)
}

// ErrorRetries returns the number of retries allowed for handling errors
func (n *EmptyNode) ErrorRetries() int {
	if n.errGuide == nil {
		return 0
	}
	return n.errGuide.Retries()
}

func (n *EmptyNode) Fail(sig node.Signal, err error) {
	n.LogErr(err)
	sig.Err = err.Error()
	sig.Status = StatusFail
	n.UpdateState(sig)
	n.StateManager().Complete()
}

// SetID sets the ID for the EmptyNode
func (n *EmptyNode) SetID(id string) {
	n.id = id
}

// ID returns the ID of the EmptyNode
func (n *EmptyNode) ID() string {
	return n.id
}

// MakeInputCh initializes the input channel with the provided buffer size
func (n *EmptyNode) MakeInputCh(size int) {
	if size > 0 {
		n.inCh = make(chan node.Signal, size)
		return
	}
	n.inCh = make(chan node.Signal)
}

// InputCh returns the input channel for the EmptyNode
func (n *EmptyNode) InputCh() chan node.Signal {
	return n.inCh
}

// SetGuidance sets the Guidance for signal generation
func (n *EmptyNode) SetGuidance(guide node.Guidance) {
	n.guide = guide
}

// GenGuidance generates guidance for the given signal using the registered guidance
func (n *EmptyNode) GenGuidance(sig node.Signal) (node.Signal, error) {
	if n.guide != nil {
		return n.guide.Generate(sig)
	}
	return sig, nil
}

// SetHooks sets the Hooks for the EmptyNode
func (n *EmptyNode) SetHooks(hooks node.Hooks) {
	n.hooks = hooks
}

// Hooks returns the Hooks associated with the EmptyNode
func (n *EmptyNode) Hooks() node.Hooks {
	return n.hooks
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

// SetLogger sets the logger for the EmptyNode
func (n *EmptyNode) SetLogger(logger node.Logger) {
	n.logger = logger
}

// Logger returns the logger associated with the EmptyNode
func (n *EmptyNode) Logger() node.Logger {
	return n.logger
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
	if n.logger != nil {
		n.logger.Log(fmt.Sprintf("{ \"severity\": %q, \"id\": %q, \"msg\": %q }", severity, id, msg))
		return
	}
	// Fallback to Stdout
	fmt.Printf("severity=%s id=%s msg=%s", severity, id, msg)
}

// SetResourceManager sets the ResourceManager for rate limiting
func (n *EmptyNode) SetResourceManager(mgr node.ResourceManager) {
	n.resourceMgr = mgr
}

// ResourceManager returns the ResourceManager associated with the EmptyNode
func (n *EmptyNode) ResourceManager() node.ResourceManager {
	return n.resourceMgr
}

// RateLimit checks and applies rate limiting for the signal using the ResourceManager
func (n *EmptyNode) RateLimit(sig node.Signal) error {
	if n.resourceMgr != nil {
		return n.resourceMgr.RateLimit(sig)
	}
	return nil
}

// SetStateManager sets the StateManager for managing the state of the node
func (n *EmptyNode) SetStateManager(mgr node.StateManager) {
	if mgr == nil {
		n.LogErr(fmt.Errorf("StateManager should not be nil"))
		return
	}
	n.stateMgr = mgr
	n.doneCh = n.stateMgr.Register()
	if n.doneCh == nil {
		n.LogErr(fmt.Errorf("failed to register with StateManager"))
	} else {
		n.LogInfo("Registered with StateManager")
	}
}

// DoneCh returns the done channel for the EmptyNode
func (n *EmptyNode) DoneCh() chan struct{} {
	return n.doneCh
}

// StateManager returns the StateManager associated with the EmptyNode
func (n *EmptyNode) StateManager() node.StateManager {
	return n.stateMgr
}

// UpdateState updates the state of the signal using the StateManager
func (n *EmptyNode) UpdateState(sig node.Signal) {
	if n.stateMgr != nil {
		n.stateMgr.UpdateState(sig)
	}
}

// ValidateSignal checks if a signal is valid (e.g., it contains an ID)
func (n *EmptyNode) ValidateSignal(sig node.Signal) error {
	if sig.NodeID == "" {
		return fmt.Errorf("invalid signal missing ID")
	}
	return nil
}

// PreProcessSignal runs before-action hooks and handles rate limiting for the signal
func (n *EmptyNode) PreProcessSignal(sig node.Signal) (node.Signal, error) {
	n.LogInfo("Received signal")
	if err := n.ValidateSignal(sig); err != nil {
		return sig, err
	}

	// Rate limiting check with exponential backoff
	for retries := 0; retries < 3; retries++ {
		if err := n.RateLimit(sig); err == nil {
			break
		}
		time.Sleep(time.Duration(retries*retries) * time.Second) // Exponential backoff
	}

	if err := n.RateLimit(sig); err != nil {
		return sig, fmt.Errorf("exceeded rate limit, could not recover")
	}

	// Run any registered before-action hooks
	return n.RunBeforeHook(sig)
}

// PostProcessSignal runs after-action hooks and updates the signal state
func (n *EmptyNode) PostProcessSignal(sig node.Signal) (node.Signal, error) {
	// Run any registered after-action hooks
	sig, err := n.RunAfterHook(sig)
	if err != nil {
		return sig, err
	}

	// Update the state of the signal after processing
	n.UpdateState(sig)

	return sig, nil
}

func (n *EmptyNode) SetWaitGroup(wg *sync.WaitGroup) {
	n.waitGroup = wg
}
