package nlib

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.IntegratorNode = (*SimpleIntegratorNode)(nil)

type SimpleIntegratorNode struct {
	EmptyNode
	integratorFunc node.IntegratorFn
	groups         []node.Group
	grpSignals     map[string][]node.Signal
	sigM           sync.RWMutex
}

func NewSimpleIntegratorNode(integratorFunc node.IntegratorFn, l node.Logger, sm node.StateManager, name string) *SimpleIntegratorNode {
	n := SimpleIntegratorNode{
		integratorFunc: integratorFunc,
	}
	n.Init(l, sm, name)

	go func() {
		select {
		case sig := <-n.inCh:
			n.processSignal(sig)
		case <-n.DoneCh():
			return
		}
	}()

	return &n
}

func (n *SimpleIntegratorNode) AddGroup(group node.Group) {
	n.groups = append(n.groups, group)
}

func (n *SimpleIntegratorNode) SetIntegratorFunc(integratorFunc node.IntegratorFn) {
	n.integratorFunc = integratorFunc
}

func (n *SimpleIntegratorNode) processSignal(sig node.Signal) {
	if n.integratorFunc == nil {
		n.LogErr(fmt.Errorf("missing required integratorFunc, failing"))
		n.StateManager().Complete()
		return
	}
	var err error
	sig, err = n.PreProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	grp := n.inGroup(sig)
	if grp != nil {
		n.addSignalToGroup(sig, grp.BatchID)
		n.processGroup(sig, grp.BatchID)
		return
	}

	sig, err = n.PostProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure the context is cancelled once we're done
	if err := n.SendToConnected(ctx, sig); err != nil {
		n.Fail(sig, err)
		return
	}
}

func (n *SimpleIntegratorNode) inGroup(sig node.Signal) *node.Group {
	for i := range n.groups {
		if len(FilterMetaKey(sig, n.groups[i].BatchID)) > 0 {
			return &n.groups[i]
		}
	}
	return nil
}

func (n *SimpleIntegratorNode) addSignalToGroup(sig node.Signal, batchID string) {
	n.sigM.Lock()
	defer n.sigM.Unlock()
	grp, ok := n.grpSignals[batchID]
	if !ok {
		grp = []node.Signal{}
	}
	grp = append(grp, sig)
	n.grpSignals[batchID] = grp
}

func (n *SimpleIntegratorNode) processGroup(sig node.Signal, batchID string) {
	if !n.isGroupComplete(batchID) {
		return // Don't have all the pieces yet
	}
	n.sigM.RLock()
	defer n.sigM.RUnlock()

	sig.Status = StatusInProcess
	results := []string{}
	for _, gsig := range n.grpSignals[batchID] {
		results = append(results, gsig.Result.String())
		if sig.History != nil {
			sig.History.AddHistory(gsig)
		}
	}
	final, err := n.integratorFunc(results)
	if err != nil {
		n.LogErr(err)
		sig.Err = err.Error()
		sig.Status = StatusFail
		n.StateManager().Complete()
		return
	}
	sig.Result = NewStringData(final)
	sig.Status = StatusSuccess
	sig, err = n.PostProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure the context is cancelled once we're done

	if err := n.SendToConnected(ctx, sig); err != nil {
		n.Fail(sig, err)
		return
	}
}

func (n *SimpleIntegratorNode) isGroupComplete(batchID string) bool {
	groupDef := node.Group{}
	found := false
	for _, def := range n.groups {
		if def.BatchID == batchID {
			groupDef = def
			found = true
			break
		}
	}
	if !found {
		return false
	}
	n.sigM.RLock()
	defer n.sigM.RUnlock()

	grp := n.grpSignals[batchID]
	if len(grp) == len(groupDef.TaskIDs) {
		return true
	}
	return true
}
