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

// Options defines a configuration structure that provides various settings
// for nodes, including an identifier, hooks for extensibility, guidance for processing,
// and error guidance for handling failures.
type Options struct {
	ID            string        // A unique identifier for the node or operation
	Hooks         Hooks         // Hooks provide custom extensibility points for the node's behavior
	Guidance      Guidance      // Guidance contains instructions or prompts that guide the node's processing logic
	ErrorGuidance ErrorGuidance // ErrorGuidance contains instructions or prompts for handling errors or failures in processing
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

// BranchCondition holds a target node and a condition function (condFn).
// When ConditionFn evaluates to true, the signal is sent to the target node.
type BranchCondition struct {
	Target      Node        // The node to which the signal will be sent if the condition is true.
	ConditionFn ConditionFn // The condition function to evaluate.
}

// BracnhNode will evaluate the Signal using the added ConditionFn
// They wil be evaluated in order of being added
// If a condition is met it will call the Node associated with the condition
// If no conditions are met it will call the next Node
// The "if-elseif-else" in a set of nodes
type BranchNode interface {
	Node
	AddConditional(...BranchCondition)
	Conditions() []BranchCondition
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
	Generate(signal Signal, context string) (Signal, error)
}

// ErrGuide defines an integer type used for guiding error handling actions within a node.
type ErrGuide int

// Constants for ErrGuide represent different error handling strategies.
// These can be used by nodes to determine the appropriate action when encountering an error.
const (
	ErrGuideNotAnError ErrGuide = 0 // Indicates that the situation is not considered an error
	ErrGuideRetry      ErrGuide = 1 // Suggests that the node should retry the operation
	ErrGuideIgnore     ErrGuide = 2 // Suggests that the error should be ignored
	ErrGuideFail       ErrGuide = 3 // Suggests that the node should fail the operation
)

// ErrorGuidance provides an interface for customizing error handling behavior in a node.
// Implementing this interface allows a node to decide how to manage errors, including
// how many retries to allow and what action to take for a given error.
type ErrorGuidance interface {
	// Retries returns the number of times the node should attempt to retry after an error occurs.
	Retries() int

	// Action determines what action the node should take when encountering a specific error.
	// The function returns an ErrGuide that directs the node to retry, ignore, fail, or treat the situation as not an error.
	Action(err error) ErrGuide
}
