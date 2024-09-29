package nlib

import (
	"fmt"
	"time"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.Node = (*EmptyNode)(nil)

// EmptyNode is the boiler plate code for a node.Node
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
}

func (n *EmptyNode) Init(l node.Logger, mgr node.StateManager, id string) {
	n.SetLogger(l)
	n.SetStateManager(mgr)
	n.SetID(id)
	n.MakeInputCh()
}

func (n *EmptyNode) Connect(nn ...node.Node) {
	n.nodes = append(n.nodes, nn...)
}

func (n *EmptyNode) SendToConnected(sig node.Signal) {
	sig = PrepareSignalForNext(sig)
	for _, conNode := range n.nodes {
		n.LogInfo(fmt.Sprintf("Sending to %s", conNode.ID()))
		sig.NodeID = conNode.ID()
		conNode.InputCh() <- sig // Send the signal to each node's input channel
	}
}

func (n *EmptyNode) ID() string {
	return n.id
}

func (n *EmptyNode) MakeInputCh() {
	n.inCh = make(chan node.Signal)
}

func (n *EmptyNode) InputCh() chan node.Signal {
	return n.inCh
}

func (n *EmptyNode) ShouldFail(err error) bool {
	if n.stateMgr != nil {
		return n.stateMgr.ShouldFail(err)
	}
	return false
}

func (n *EmptyNode) SetGuidance(guide node.Guidance) {
	n.guide = guide
}

func (n *EmptyNode) GenGuidance(sig node.Signal) (node.Signal, error) {
	if n.guide != nil {
		return n.guide.Generate(sig)
	}
	return sig, nil
}

func (n *EmptyNode) SetHooks(hooks node.Hooks) {
	n.hooks = hooks
}

func (n *EmptyNode) Hooks() node.Hooks {
	return n.hooks
}

func (n *EmptyNode) RunBeforeHook(sig node.Signal) (node.Signal, error) {
	if n.hooks != nil {
		return n.hooks.BeforeAction(sig)
	}
	return sig, nil
}

func (n *EmptyNode) RunAfterHook(sig node.Signal) (node.Signal, error) {
	if n.hooks != nil {
		return n.hooks.AfterAction(sig)
	}
	return sig, nil
}

func (n *EmptyNode) SetID(id string) {
	n.id = id
}

func (n *EmptyNode) SetLogger(logger node.Logger) {
	n.logger = logger
}

func (n *EmptyNode) Logger() node.Logger {
	return n.logger
}

func (n *EmptyNode) LogErr(err error) {
	n.log(fmt.Sprintf("[ERROR] %s %v", n.id, err))
}

func (n *EmptyNode) LogInfo(msg string) {
	n.log(fmt.Sprintf("[INFO] %s %s", n.id, msg))
}

func (n *EmptyNode) LogError(msg string) {
	n.log(fmt.Sprintf("[ERROR] %s %s", n.id, msg))
}

func (n *EmptyNode) LogDebug(msg string) {
	n.log(fmt.Sprintf("[DEBUG] %s %s", n.id, msg))
}

func (n *EmptyNode) log(msg string) {
	if n.logger != nil {
		n.logger.Log(msg)
	}
}

func (n *EmptyNode) SetResourceManager(mgr node.ResourceManager) {
	n.resourceMgr = mgr
}

func (n *EmptyNode) ResourceManager() node.ResourceManager {
	return n.resourceMgr
}

func (n *EmptyNode) RateLimit(sig node.Signal) error {
	if n.resourceMgr != nil {
		return n.resourceMgr.RateLimit(sig)
	}
	return nil
}

func (n *EmptyNode) SetStateManager(mgr node.StateManager) {
	n.stateMgr = mgr
	n.doneCh = n.stateMgr.Register()
}

func (n *EmptyNode) DoneCh() chan struct{} {
	return n.doneCh
}

func (n *EmptyNode) StateManager() node.StateManager {
	return n.stateMgr
}

func (n *EmptyNode) UpdateState(sig node.Signal) {
	if n.stateMgr != nil {
		n.stateMgr.UpdateState(sig)
	}
}

func (n *EmptyNode) PreProcessSignal(sig node.Signal) node.Signal {
	n.LogInfo("Received signal")
	// Rate limiting check
	if n.RateLimit(sig) != nil {
		time.Sleep(1 * time.Second) // If rate-limited, pause for 1 second
	}

	// Run any registered before-action hooks
	sig, err := n.RunBeforeHook(sig)
	if err != nil {
		n.LogErr(err) // Log any errors from the before-hook
	}
	return sig
}

func (n *EmptyNode) PostProcesSignal(sig node.Signal) node.Signal {
	// Run any registered after-action hooks
	sig, err := n.RunAfterHook(sig)
	if err != nil {
		n.LogErr(err) // Log any errors from the after-hook
	}

	// Update the state of the signal after processing
	n.UpdateState(sig)

	return sig
}
