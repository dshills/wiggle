package nlib

import (
	"sync"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.PartitionerNode = (*SimplePartitionerNode)(nil)

type SimplePartitionerNode struct {
	EmptyNode
	partitionFunc node.PartitionerFn
	childNodes    []node.Node
	waitGroup     *sync.WaitGroup
}

func NewSimplePartitionerNode(partitionFunc node.PartitionerFn, l node.Logger, sm node.StateManager, name string) *SimplePartitionerNode {
	n := SimplePartitionerNode{
		partitionFunc: partitionFunc,
		waitGroup:     &sync.WaitGroup{},
	}
	n.Init(l, sm, name)

	go func() {
		for {
			select {
			case sig := <-n.inCh:
				n.processSignal(sig)
			case <-n.doneCh:
				return
			}
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

func (n *SimplePartitionerNode) processSignal(sig node.Signal) {
	sig = n.PreProcessSignal(sig)

	// Partition the signal's data
	parts, err := n.partitionFunc(sig.Data.String())
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
				newSignal := sig // Clone signal to avoid mutation issues
				newSignal.Data = NewStringData(part)
				child.InputCh() <- newSignal
			}(child, part)
		}
	}
	n.waitGroup.Wait()

	sig = n.PostProcesSignal(sig)
	n.SendToConnected(sig)
}
