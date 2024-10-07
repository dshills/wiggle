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
type StateManager interface {
	GetState(Signal) State
	UpdateState(Signal)

	ContextManager() ContextManager
	Coordinator() Coordinator
	HistoryManager() HistoryManager
	Logger() Logger
	ResourceManager() ResourceManager

	SetContextManager(ContextManager)
	SetCoordinator(Coordinator)
	SetHistoryManager(HistoryManager)
	SetLogger(Logger)
	SetResourceManager(ResourceManager)

	Register() chan struct{}
	Complete()
	WaitFor(Node)

	Log(string)
	GetContext(key string) (DataCarrier, error)
	SetContext(key string, data DataCarrier)
	RemoveContext(key string)

	AddHistory(Signal)
	GetHistory() []Signal
	FilterHistory(nodeid string) []Signal
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

// ContextManager is responsible for managing the contextual information
// passed between nodes during execution. It provides methods to update
// and retrieve context from signals, ensuring that relevant data is available
// and consistent as it flows through the node chain. This interface helps
// maintain continuity and relevance in processing workflows.
type ContextManager interface {
	GetContext(key string) (DataCarrier, error)
	RemoveContext(key string)
	SetContext(key string, data DataCarrier)
}

// HistoryManager is responsible for managing the history of signals as they pass
// through nodes. It provides methods to add entries, retrieve, and optionally
// compress or truncate the history, allowing nodes to track the progression of
// a signal and maintain a record of its transformations throughout the workflow.
type HistoryManager interface {
	AddHistory(Signal)             // Adds a new entry to history
	CompressHistory() error        // Compress or truncate history
	GetHistory() []Signal          // Retrieve full history
	Filter(nodeid string) []Signal // Get specific history
}
