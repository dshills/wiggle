package nlib

import (
	"sync"
	"time"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.IntegratorNode = (*SimpleIntegratorNode)(nil)

type SimpleIntegratorNode struct {
	EmptyNode
	integratorFunc node.IntegratorFn
	childNodes     []node.Node
	resultCh       chan string
	waitGroup      *sync.WaitGroup
}

func NewSimpleIntegratorNode(integratorFunc node.IntegratorFn, l node.Logger, sm node.StateManager, name string) *SimpleIntegratorNode {
	n := SimpleIntegratorNode{
		integratorFunc: integratorFunc,
		resultCh:       make(chan string, 10), // Buffered to store results from child nodes
		waitGroup:      &sync.WaitGroup{},
	}
	n.Init(l, sm, name)

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

func (n *SimpleIntegratorNode) SetIntegratorFunc(integratorFunc node.IntegratorFn) {
	n.integratorFunc = integratorFunc
}

func (n *SimpleIntegratorNode) SetChildNodes(nodes ...node.Node) {
	n.childNodes = nodes
}

func (n *SimpleIntegratorNode) Wait() {
	n.waitGroup.Wait()
}

func (n *SimpleIntegratorNode) processSignal(sig node.Signal) {
	var results []string

	sig = n.PreProcessSignal(sig)

	// Collect results from child nodes
	for _, child := range n.childNodes {
		n.waitGroup.Add(1)
		go func(child node.Node) {
			defer n.waitGroup.Done()

			childSignal := sig // Clone the signal for the child node
			child.InputCh() <- childSignal

			select {
			case result := <-n.resultCh:
				results = append(results, result)
			case <-time.After(5 * time.Second): // Timeout per node's response
			}
		}(child)
	}

	// Wait for all child nodes to finish
	n.waitGroup.Wait()

	// Combine the results using the integrator function
	finalResult, err := n.integratorFunc(results)
	if err != nil {
		n.LogErr(err)
	}
	sig.Response = NewStringData(finalResult)

	sig = n.PostProcesSignal(sig)
	n.SendToConnected(sig)
}
