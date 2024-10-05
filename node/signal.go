package node

import "github.com/dshills/wiggle/schema"

// Signal represents the core data structure passed between nodes in a processing chain.
// It contains the data being processed, contextual information, metadata, response data,
// and a history of transformations. Signals enable the flow of information across nodes,
// allowing each node to modify, route, and act on the data while keeping track of its
// progression throughout the workflow.
type Signal struct {
	Context ContextManager
	Err     string
	History HistoryManager
	Meta    []Meta
	NodeID  string
	Result  DataCarrier
	Schema  schema.Schema
	Status  string
	Task    DataCarrier
}

// Meta represents key-value pairs of metadata associated with a signal.
// It is used to store additional information that may be relevant for
// processing, such as configuration settings, model parameters, or
// contextual data, allowing nodes to access and act on this metadata
// as part of the workflow.
type Meta struct {
	Key   string
	Value string
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
	AddHistory(Signal)                          // Adds a new entry to history
	CompressHistory() error                     // Compress or truncate history
	GetHistory() []Signal                       // Retrieve full history
	GetHistoryByID(id string) ([]Signal, error) // Get specific history
}

// DataCarrier provides an abstraction for handling different types of data
// within a signal. It allows for conversion of the data into various formats,
// such as string, JSON, or vectors, ensuring flexibility in how data
// is passed between nodes and processed in different stages of the workflow.
type DataCarrier interface {
	JSON() []byte
	String() string
	Vector() []float32
}
