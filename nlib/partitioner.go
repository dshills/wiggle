package nlib

import (
	"sync"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.Node = (*SimplePartitionerNode)(nil)

type SimplePartitionerNode struct {
	EmptyNode
	partitionFunc node.PartitionerFn
	childNodes    []node.Node
	waitGroup     *sync.WaitGroup
}

func NewSimplePartitionerNode(partitionFunc node.PartitionerFn, l node.Logger, sm node.StateManager) *SimplePartitionerNode {
	n := SimplePartitionerNode{
		partitionFunc: partitionFunc,
		waitGroup:     &sync.WaitGroup{},
	}
	n.SetLogger(l)
	n.SetStateManager(sm)
	n.SetID(generateUUID())
	n.MakeInputCh()

	go func() {
		select {
		case sig := <-n.inCh:
			n.processSignal(sig)
		case <-n.doneCh:
			return
		}
	}()

	return &n
}

func (n *SimplePartitionerNode) SetPartitionFunc(partitionFunc node.PartitionerFn) {
	n.partitionFunc = partitionFunc
}

func (n *SimplePartitionerNode) SetChildNodes(nodes ...node.Node) {
	n.childNodes = nodes
}

func (n *SimplePartitionerNode) processSignal(signal node.Signal) {
	// Partition the signal's data
	parts, err := n.partitionFunc(signal.Data.String())
	if err != nil {
		n.LogErr(err)
		return
	}

	// Send each partitioned part to child nodes for processing
	for _, part := range parts {
		for _, child := range n.childNodes {
			n.waitGroup.Add(1)
			go func(child node.Node, part string) {
				defer n.waitGroup.Done()
				newSignal := signal // Clone signal to avoid mutation issues
				newSignal.Data = NewStringData(part)
				child.InputCh() <- newSignal
			}(child, part)
		}
	}
	n.waitGroup.Wait()
	n.UpdateState(signal)
}
