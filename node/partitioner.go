package node

import "sync"

type SimplePartitionerNode struct {
	id            string
	logger        Logger
	errorHandler  ErrorHandler
	partitionFunc PartitionerFn
	childNodes    []Node
	inputCh       chan Signal
	waitGroup     *sync.WaitGroup
}

func NewSimplePartitionerNode(id string, partitionFunc PartitionerFn) *SimplePartitionerNode {
	return &SimplePartitionerNode{
		id:            id,
		partitionFunc: partitionFunc,
		inputCh:       make(chan Signal),
		waitGroup:     &sync.WaitGroup{},
	}
}

func (n *SimplePartitionerNode) ID() string {
	return n.id
}

func (n *SimplePartitionerNode) SetID(id string) {
	n.id = id
}

func (n *SimplePartitionerNode) InputCh() chan Signal {
	return n.inputCh
}

func (n *SimplePartitionerNode) SetLogger(logger Logger) {
	n.logger = logger
}

func (n *SimplePartitionerNode) SetErrorHandler(handler ErrorHandler) {
	n.errorHandler = handler
}

func (n *SimplePartitionerNode) SetPartitionFunc(partitionFunc PartitionerFn) {
	n.partitionFunc = partitionFunc
}

func (n *SimplePartitionerNode) SetChildNodes(nodes ...Node) {
	n.childNodes = nodes
}

func (n *SimplePartitionerNode) Wait() {
	n.waitGroup.Wait()
}

func (n *SimplePartitionerNode) processSignal(signal Signal) {
	// Partition the signal's data
	parts, err := n.partitionFunc(signal.Data.ToMessageList()[0].Content)
	if err != nil {
		if n.errorHandler != nil {
			n.errorHandler.HandleError(signal, err)
		}
		return
	}

	// Send each partitioned part to child nodes for processing
	for _, part := range parts {
		for _, child := range n.childNodes {
			n.waitGroup.Add(1)
			go func(child Node, part string) {
				defer n.waitGroup.Done()
				newSignal := signal // Clone signal to avoid mutation issues
				newSignal.Data = MessageData{Message: part}
				child.InputCh() <- newSignal
			}(child, part)
		}
	}
}

func (n *SimplePartitionerNode) Run() {
	for signal := range n.inputCh {
		n.processSignal(signal)
	}
}
