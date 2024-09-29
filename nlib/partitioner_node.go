package nlib

import (
	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.PartitionerNode = (*SimplePartitionerNode)(nil)

type SimplePartitionerNode struct {
	EmptyNode
	partitionFunc node.PartitionerFn
	factory       node.Factory
	integrator    node.IntegratorNode
}

func NewSimplePartitionerNode(partitionFunc node.PartitionerFn, l node.Logger, sm node.StateManager, name string) *SimplePartitionerNode {
	n := SimplePartitionerNode{
		partitionFunc: partitionFunc,
	}
	n.Init(l, sm, name)

	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.processSignal(sig)
			case <-n.DoneCh():
				return
			}
		}
	}()

	return &n
}

func (n *SimplePartitionerNode) SetPartitionFunc(partitionFunc node.PartitionerFn) {
	n.partitionFunc = partitionFunc
}

func (n *SimplePartitionerNode) SetIntegrator(integrator node.IntegratorNode) {
	n.integrator = integrator
}

func (n *SimplePartitionerNode) SetNodeFactory(factory node.Factory) {
	n.factory = factory
}

func (n *SimplePartitionerNode) processSignal(sig node.Signal) {
	if n.partitionFunc == nil || n.factory == nil {
		n.LogError("missing required partition function or node factory function, failing")
		n.StateManager().Complete()
		return
	}
	sig = n.PreProcessSignal(sig)

	sig.Status = StatusInProcess

	// Partition the signal's data
	parts, err := n.partitionFunc(sig.Task.String())
	if err != nil {
		n.LogErr(err)
		sig.Err = err.Error()
		return
	}

	sig.Status = StatusSuccess
	// Create a set of Nodes to handle the partitioned data
	nodes := n.factory(len(parts))

	// Tell the IntegratorNode what's coming
	group := NewGroup(n.ID(), nodes...)
	if n.integrator != nil {
		n.integrator.AddGroup(group)
	}

	// Send it...
	for i, task := range parts {
		newSig := SignalFromSignal(sig, NewStringData(task)) // Create a new Signal based on the current
		newSig = GroupSignal(newSig, group, nodes[i])        // Add the Group meta tracking data
		newSig.NodeID = nodes[i].ID()                        // Set the new Node ID
		nodes[i].InputCh() <- newSig                         // Send to Node
	}

	sig = n.PostProcesSignal(sig)
	n.SendToConnected(sig)
}
