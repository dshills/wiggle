package nlib

import (
	"fmt"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.Node = (*EmptyNode)(nil)

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

func (n *EmptyNode) Connect(nn node.Node) {
	n.nodes = append(n.nodes, nn)
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

func (n *EmptyNode) RateLimit(sig node.Signal) error {
	if n.resourceMgr != nil {
		return n.resourceMgr.RateLimit(sig)
	}
	return nil
}

func (n *EmptyNode) SetStateManager(mgr node.StateManager) {
	n.stateMgr = mgr
}

func (n *EmptyNode) UpdateState(sig node.Signal) {
	if n.stateMgr != nil {
		n.stateMgr.UpdateState(sig)
	}
}
