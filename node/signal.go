package node

// Signal represents the core data structure passed between nodes in a processing chain.
// It contains the data being processed, contextual information, metadata, response data,
// and a history of transformations. Signals enable the flow of information across nodes,
// allowing each node to modify, route, and act on the data while keeping track of its
// progression throughout the workflow.
type Signal struct {
	NodeID   string
	Data     DataCarrier
	Response DataCarrier
	Context  ContextManager
	Meta     []Meta
	History  HistoryManager
	Err      error
	Status   string
}

func (s *Signal) SetContext(forID string, data DataCarrier) {
	if s.Context != nil {
		s.Context.SetContext(forID, data)
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

// NewSignal will return a new Signal. This is typically used to generate the
// initial Signal at the start of processing
func NewSignal(id string, cm ContextManager, hx HistoryManager, task DataCarrier, meta ...Meta) Signal {
	return Signal{
		NodeID:   id,
		Context:  cm,
		History:  hx,
		Meta:     meta,
		Response: task, // Nodes read tasks from the previous response
	}
}

// SignalFromSignal will save the original Signal in history
// Convert the response to the incomming data and update the
// NodeID. This is commonly used when a Node receives a signal
// from the previous Node
func SignalFromSignal(id string, sig Signal) Signal {
	if sig.History != nil {
		sig.History.AddHistory(sig)
	}
	sig.Data = sig.Response
	sig.NodeID = id
	sig.Response = nil
	return sig
}
