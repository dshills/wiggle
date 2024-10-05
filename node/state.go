package node

import "time"

type State struct {
	Completed int
	Failures  int
	Status    string
}

// StateManager is responsible for tracking and updating the state of a signal
// as it moves through the node chain. It provides methods to update and retrieve
// the current state, ensuring that each node can access the relevant status
// information. This interface helps maintain workflow consistency and manage
// state transitions across complex processes.
// Defines how state and errors are managed during node execution. It provides
// a mechanism to handle errors encountered in a node's processing logic and
// determines whether the workflow should continue or halt based on the nature
// of the error. This interface helps to create robust error management strategies
// across a chain of nodes.
type StateManager interface {
	Complete()
	Coordinator() Coordinator
	GetState(Signal) State
	Log(string)
	Register() chan struct{}
	ResourceManager() ResourceManager
	SetCoordinator(Coordinator)
	SetLogger(Logger)
	SetResourceManager(ResourceManager)
	UpdateState(Signal)
	WaitFor(Node)
}

// Coordinator is responsible for managing the synchronization and execution flow
// across multiple nodes. It provides mechanisms such as waiting for the completion
// of tasks, handling timeouts, and coordinating the parallel execution of nodes,
// ensuring that complex workflows proceed smoothly and efficiently.
type Coordinator interface {
	WaitForCompletion(nodes ...Node) error
	CancelOnTimeout(duration time.Duration)
}

// Logger is responsible for logging messages during the execution of nodes.
// It provides a simple interface to track the flow of data, errors, and
// other significant events in the system, aiding in debugging and monitoring
// the behavior of nodes in a chain of tasks.
type Logger interface {
	Log(string)
}

// ResourceManager handles the allocation and management of resources
// during the execution of nodes. It provides mechanisms such as rate limiting
// and resource throttling to prevent overuse of external systems (e.g., APIs or databases).
// This interface ensures efficient and controlled use of resources across the node chain.
type ResourceManager interface {
	RateLimit(Signal) error
}
