package node

import "github.com/dshills/wiggle/llm"

// Signal represents the core data structure passed between nodes in a processing chain.
// It contains the data being processed, contextual information, metadata, response data,
// and a history of transformations. Signals enable the flow of information across nodes,
// allowing each node to modify, route, and act on the data while keeping track of its
// progression throughout the workflow.
type Signal struct {
	NodeID   string
	Data     DataCarrier
	Response llm.Message
	Context  string
	Meta     []Meta
	History  HistoryManager
}

// Result encapsulates the outcome of processing a signal within a node.
// It contains the processed value and any error encountered during execution.
// The Result struct allows nodes to communicate the success or failure of a task
// and pass along the output for further processing in the workflow.
type Result struct {
	Value string
	Error error
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

// HistoryManager is responsible for managing the history of signals as they pass
// through nodes. It provides methods to add entries, retrieve, and optionally
// compress or truncate the history, allowing nodes to track the progression of
// a signal and maintain a record of its transformations throughout the workflow.
type HistoryManager interface {
	AddHistory(Signal, Signal) error // Adds a new entry to history
	CompressHistory(Signal) error    // Compress or truncate history
	GetHistory(Signal) []Signal      // Retrieve full history
}

// DataCarrier provides an abstraction for handling different types of data
// within a signal. It allows for conversion of the data into various formats,
// such as message lists, JSON, or vectors, ensuring flexibility in how data
// is passed between nodes and processed in different stages of the workflow.
type DataCarrier interface {
	ToMessageList() llm.MessageList
	ToJSON() string
	ToVector() []float32
}
