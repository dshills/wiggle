package nlib

import (
	"context"
	"fmt"
	"time"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.PartitionerNode = (*SimplePartitionerNode)(nil)

type SimplePartitionerNode struct {
	EmptyNode
	partitionFunc   node.PartitionerFn
	integrationFunc node.IntegratorFn
	factory         node.Factory
}

func NewSimplePartitionerNode(pfn node.PartitionerFn, ifn node.IntegratorFn, fac node.Factory, mgr node.StateManager, options node.Options) *SimplePartitionerNode {
	n := SimplePartitionerNode{
		partitionFunc:   pfn,
		integrationFunc: ifn,
		factory:         fac,
	}
	n.SetOptions(options)
	n.SetStateManager(mgr)
	n.MakeInputCh()

	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.LogInfo("Received Signal")
				n.processSignal(sig)
			case <-n.StateManager().Register():
				n.LogInfo("Received Done")
				return
			}
		}
	}()

	return &n
}

func (n *SimplePartitionerNode) SetPartitionFunc(partitionFunc node.PartitionerFn) {
	n.partitionFunc = partitionFunc
}

func (n *SimplePartitionerNode) SetIntegrationFunc(integratonFunc node.IntegratorFn) {
	n.integrationFunc = integratonFunc
}

func (n *SimplePartitionerNode) SetNodeFactory(factory node.Factory) {
	n.factory = factory
}

func (n *SimplePartitionerNode) processSignal(sig node.Signal) {
	var err error
	if n.partitionFunc == nil || n.factory == nil || n.integrationFunc == nil {
		err := fmt.Errorf("partition, integrator, and factory functions are required")
		n.Fail(sig, err)
		return
	}
	sig, err = n.PreProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	sig.Status = StatusInProcess

	// Partition the signal's data
	parts, err := n.partitionFunc(sig.Task.String())
	if err != nil {
		n.LogErr(err)
		sig.Err = err.Error()
		return
	}

	// Create a set of Nodes to handle the partitioned data
	nodes := n.factory(len(parts))
	respChan := make(chan node.Signal, len(parts))
	emptyNode := &EmptyNode{inputCh: respChan}
	for i, task := range parts {
		newSig := SignalFromSignal(sig, NewStringData(task))
		newSig.NodeID = nodes[i].ID()
		nodes[i].Connect(emptyNode)
		nodes[i].InputCh() <- newSig
	}
	respList := []string{}
	for i := 0; i < len(parts); i++ {
		recSig := <-respChan
		respList = append(respList, recSig.Task.String())
	}
	response, err := n.integrationFunc(respList)
	if err != nil {
		n.LogErr(err)
		sig.Err = err.Error()
		return
	}
	sig.Result = NewStringData(response)
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
