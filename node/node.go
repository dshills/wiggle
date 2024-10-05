package node

import (
	"io"
)

// Node represents a generic processing unit in a chain of tasks.
// It processes incoming signals, executes actions (e.g., transforming data, querying a model),
// and forwards the processed signal to connected nodes. The interface allows for modular design,
// enabling different types of nodes (e.g., AI nodes, partitioners, integrators)
// to be chained together, ensuring flexible and scalable workflows.
type Node interface {
	Connect(...Node)
	ID() string
	InputCh() chan Signal
	SetID(string)
	SetOptions(Options)
	SetStateManager(StateManager)
}

type Options struct {
	ID            string
	Hooks         Hooks
	Guidance      Guidance
	ErrorGuidance ErrorGuidance
}

// PartitionerFn is a function type that takes an input string and splits it
// into smaller parts or tasks. It is used by PartitionerNodes to divide
// large or complex data into manageable chunks, enabling parallel processing
// by multiple nodes in the chain.
type PartitionerFn func(string) ([]string, error)

// IntegratorFn is a function type that takes the results of partitioned tasks
// as input and combines them into a single, coherent output. It is used by
// PartitionerNode to aggregate and merge processed data, ensuring that the
// final output is consistent and meaningful.
type IntegratorFn func([]string) (string, error)

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
	SetIntegrationFunc(IntegratorFn)
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

// HookFn is used to transform Signals. It is used with the Hooks interface
type HookFn func(Signal) (Signal, error)

// Hooks provides lifecycle hooks that allow custom logic to be executed
// before and after a node processes a signal. The interface allows for
// pre-processing, such as input validation or logging, and post-processing,
// such as result modification or cleanup, giving more control over the
// execution flow in a chain of nodes.
type Hooks interface {
	BeforeAction(Signal) (Signal, error)
	AfterAction(Signal) (Signal, error)
}

// Guidance provides a mechanism to generate structured guidance or instructions
// for processing signals within the node chain. It interprets the input data,
// contextual information, and metadata to formulate a message or set of instructions
// that can guide further processing by LLMs or other nodes in the workflow.
type Guidance interface {
	// Generate processes the input signal, taking into account its data, context,
	// and metadata, and returns a modified Signal with the guidance for further steps.
	Generate(signal Signal) (Signal, error)
}

type ErrGuide int

const (
	ErrGuideNotAnError ErrGuide = 0
	ErrGuideRetry      ErrGuide = 1
	ErrGuideIgnore     ErrGuide = 2
	ErrGuideFail       ErrGuide = 3
)

// ErrorGuidance provides guidance on how a particular node
// should manage errors. It is not required and default behavior on error
// is to fail
type ErrorGuidance interface {
	Retries() int
	Action(err error) ErrGuide
}
