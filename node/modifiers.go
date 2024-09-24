package node

import "time"

// Logger is responsible for logging messages during the execution of nodes.
// It provides a simple interface to track the flow of data, errors, and
// other significant events in the system, aiding in debugging and monitoring
// the behavior of nodes in a chain of tasks.
type Logger interface {
	Log(string)
}

// ErrorHandler defines how errors are managed during node execution. It provides
// a mechanism to handle errors encountered in a node's processing logic and
// determines whether the workflow should continue or halt based on the nature
// of the error. This interface helps to create robust error management strategies
// across a chain of nodes.
type ErrorHandler interface {
	HandleError(Signal, error) bool
}

// ResourceManager handles the allocation and management of resources
// during the execution of nodes. It provides mechanisms such as rate limiting
// and resource throttling to prevent overuse of external systems (e.g., APIs or databases).
// This interface ensures efficient and controlled use of resources across the node chain.
type ResourceManager interface {
	RateLimit(Signal) error
}

// ContextManager is responsible for managing the contextual information
// passed between nodes during execution. It provides methods to update
// and retrieve context from signals, ensuring that relevant data is available
// and consistent as it flows through the node chain. This interface helps
// maintain continuity and relevance in processing workflows.
type ContextManager interface {
	UpdateContext(Signal) (string, error)
	GetContext(Signal) string
}

// NodeHooks provides lifecycle hooks that allow custom logic to be executed
// before and after a node processes a signal. The interface allows for
// pre-processing, such as input validation or logging, and post-processing,
// such as result modification or cleanup, giving more control over the
// execution flow in a chain of nodes.
type NodeHooks interface {
	BeforeAction(Signal) Result
	AfterAction(Signal) Result
}

// StateManager is responsible for tracking and updating the state of a signal
// as it moves through the node chain. It provides methods to update and retrieve
// the current state, ensuring that each node can access the relevant status
// information. This interface helps maintain workflow consistency and manage
// state transitions across complex processes.
type StateManager interface {
	UpdateState(Signal, string) error
	GetState(Signal) string
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

// Coordinator is responsible for managing the synchronization and execution flow
// across multiple nodes. It provides mechanisms such as waiting for the completion
// of tasks, handling timeouts, and coordinating the parallel execution of nodes,
// ensuring that complex workflows proceed smoothly and efficiently.
type Coordinator interface {
	WaitForCompletion(nodes ...Node) error
	CancelOnTimeout(duration time.Duration)
}
