package node

import "io"

// Node represents a generic processing unit in a chain of tasks.
// It processes incoming signals, executes actions (e.g., transforming data, querying a model),
// and forwards the processed signal to connected nodes. The interface allows for modular design,
// enabling different types of nodes (e.g., AI nodes, partitioners, integrators)
// to be chained together, ensuring flexible and scalable workflows.
type Node interface {
	Connect(...Node)
	ID() string
	InputCh() chan Signal
	SetErrorGuidance(ErrorGuidance)
	SetGuidance(Guidance)
	SetHooks(Hooks)
	SetID(string)
	SetLogger(Logger)
	SetResourceManager(ResourceManager)
	SetStateManager(StateManager)
}

// PartitionerFn is a function type that takes an input string and splits it
// into smaller parts or tasks. It is used by PartitionerNodes to divide
// large or complex data into manageable chunks, enabling parallel processing
// by multiple nodes in the chain.
type PartitionerFn func(string) ([]string, error)

// Factory is a function that will return an arbitrary number of Nodes
// It is used by PartitionerNode after breaking a Signal into smaller parts
// It will call the Factory function to create nodes to process the chunks.
// This can be as simple as a single Node or it could return a Set Node that
// is comprised of many Nodes.
type Factory func(count int) []Node

// PartitionerNode splits the input signal into smaller tasks or chunks
// using a specified partition function. These partitions are then distributed
// to child nodes for parallel processing. The interface allows for efficient
// handling of large or complex data by breaking it down into manageable parts
// that can be processed independently.
type PartitionerNode interface {
	Node
	SetPartitionFunc(partitionFunc PartitionerFn)
	SetNodeFactory(Factory)
	SetIntegrator(IntegratorNode)
}

// IntegratorFn is a function type that takes the results of partitioned tasks
// as input and combines them into a single, coherent output. It is used by
// IntegratorNodes to aggregate and merge processed data, ensuring that the
// final output is consistent and meaningful.
type IntegratorFn func([]string) (string, error)

// Group describes a partitioned set of tasks
// When passed to an integrator it will collect the outcome of the
// batch before processing
// Signals are identified by Meta entries that include
// the BatchID and the individual TaskID
type Group struct {
	OriginatorID string
	BatchID      string
	TaskIDs      []string
}

// IntegratorNode gathers and combines the results from a PartitionerNode
// into a single coherent output using a specified integrator function.
// It ensures that the partitioned tasks, once processed, are merged back into
// a unified result, maintaining data consistency and flow across the node chain.
type IntegratorNode interface {
	Node
	SetIntegratorFunc(integratorFunc IntegratorFn)
	AddGroup(Group)
}

// ConditionFn is a function type that takes a Signal and returns
// true or false depending on if the Signal meets some criteria
// true = It met the condition, false = Did not meet the condition
type ConditionFn func(Signal) bool

// LoopNode will send it's Signal to the start Node until
// the ConditionFn returns true then it will call the next node
// Used in conjunction with Hook function it can accumulate multiple runs
// or it can rerun until a specific answer is met or simply a set number of times
// The "for" loop in a set of nodes
type LoopNode interface {
	Node
	SetStartNode(Node)
	SetConditionFunc(ConditionFn)
}

// BracnhNode will evaluate the Signal using the added ConditionFn
// They wil be evaluated in order of being added
// If a condition is met it will call the Node associated with the condition
// If no conditions are met it will call the next Node
// The "if-elseif-else" in a set of nodes
type BranchNode interface {
	Node
	AddConditional(Node, ConditionFn)
}

// OutputNode writes data to a writer
type OutputNode interface {
	Node
	SetWriter(io.Writer)
}

// InputNode reades data from the reader
type InputNode interface {
	Node
	SetReader(io.Reader)
}

// SetNode represents a collection of interconnected nodes forming a processing pipeline.
// It defines the starting node and provides mechanisms for setting the name and
// managing the execution flow of the node chain. The Set interface allows for
// orchestrating complex workflows by organizing nodes into coherent processing units.
type SetNode interface {
	Node
	SetStartNode(Node)
	SetFinalNode(Node)
	SetCoordinator(Coordinator)
}
