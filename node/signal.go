package node

// Signal represents the core data structure passed between nodes in a processing chain.
// It contains the data being processed, contextual information, metadata, response data,
// and a history of transformations. Signals enable the flow of information across nodes,
// allowing each node to modify, route, and act on the data while keeping track of its
// progression throughout the workflow.
type Signal struct {
	NodeID  string
	Task    DataCarrier
	Context ContextManager
	Result  DataCarrier
	Schema  string
	Meta    []Meta
	History HistoryManager
	Err     string
	Status  string
}

func (s *Signal) AddContext(forID string, data DataCarrier) {
	if s.Context != nil {
		s.Context.SetContext(forID, data)
	}
}

func (s *Signal) AddHistory() {
	if s.History != nil {
		s.History.AddHistory(*s)
	}
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
