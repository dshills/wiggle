package node

import (
	"fmt"
	"sync"
	"time"
)

type SimpleIntegratorNode struct {
	id             string
	logger         Logger
	errorHandler   ErrorHandler
	integratorFunc IntegratorFn
	childNodes     []Node
	inputCh        chan Signal
	resultCh       chan string
	waitGroup      *sync.WaitGroup
}

func NewSimpleIntegratorNode(id string, integratorFunc IntegratorFn) *SimpleIntegratorNode {
	return &SimpleIntegratorNode{
		id:             id,
		integratorFunc: integratorFunc,
		inputCh:        make(chan Signal),
		resultCh:       make(chan string, 10), // Buffered to store results from child nodes
		waitGroup:      &sync.WaitGroup{},
	}
}

func (n *SimpleIntegratorNode) ID() string {
	return n.id
}

func (n *SimpleIntegratorNode) SetID(id string) {
	n.id = id
}

func (n *SimpleIntegratorNode) InputCh() chan Signal {
	return n.inputCh
}

func (n *SimpleIntegratorNode) SetLogger(logger Logger) {
	n.logger = logger
}

func (n *SimpleIntegratorNode) SetErrorHandler(handler ErrorHandler) {
	n.errorHandler = handler
}

func (n *SimpleIntegratorNode) SetIntegratorFunc(integratorFunc IntegratorFn) {
	n.integratorFunc = integratorFunc
}

func (n *SimpleIntegratorNode) SetChildNodes(nodes ...Node) {
	n.childNodes = nodes
}

func (n *SimpleIntegratorNode) Wait() {
	n.waitGroup.Wait()
}

func (n *SimpleIntegratorNode) processSignal(signal Signal) {
	var results []string

	// Collect results from child nodes
	for _, child := range n.childNodes {
		n.waitGroup.Add(1)
		go func(child Node) {
			defer n.waitGroup.Done()

			childSignal := signal // Clone the signal for the child node
			child.InputCh() <- childSignal

			select {
			case result := <-n.resultCh:
				results = append(results, result)
			case <-time.After(5 * time.Second): // Timeout per node's response
				if n.errorHandler != nil {
					n.errorHandler.HandleError(signal, fmt.Errorf("timeout in node %s", child.ID()))
				}
			}
		}(child)
	}

	// Wait for all child nodes to finish
	n.waitGroup.Wait()

	// Combine the results using the integrator function
	finalResult, err := n.integratorFunc(results)
	if err != nil && n.errorHandler != nil {
		n.errorHandler.HandleError(signal, err)
	}

	// Log the final result
	if n.logger != nil {
		n.logger.Log(fmt.Sprintf("Final integrated result: %s", finalResult))
	}
}

func (n *SimpleIntegratorNode) Run() {
	for signal := range n.inputCh {
		n.processSignal(signal)
	}
}
