package node

// Signal represents the core data structure passed between nodes in a processing chain.
// It contains the data being processed, contextual information, metadata, response data,
// and a history of transformations. Signals enable the flow of information across nodes,
// allowing each node to modify, route, and act on the data while keeping track of its
// progression throughout the workflow.
type Signal struct {
	Err    string
	Meta   []Meta
	NodeID string
	Result DataCarrier
	Status string
	Task   DataCarrier
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

// DataCarrier provides an abstraction for handling different types of data
// within a signal. It allows for conversion of the data into various formats,
// such as string, JSON, or vectors, ensuring flexibility in how data
// is passed between nodes and processed in different stages of the workflow.
type DataCarrier interface {
	JSON() []byte
	String() string
	Vector() [][]float32
	Base64() []string
	ImageURLs() []string
}
