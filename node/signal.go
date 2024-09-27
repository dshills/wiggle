package node

import (
	"fmt"
)

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
	Err      string
	Status   string
}

func (s *Signal) AsLog() string {
	return fmt.Sprintf("{ NodeID: %s, data: %v, Response: %v, Err: %s, Status: %s }", s.NodeID, s.Data, s.Response, s.Err, s.Status)
}

func (s *Signal) SetContext(forID string, data DataCarrier) {
	if s.Context != nil {
		s.Context.SetContext(forID, data)
	}
}

// Before sending to the next node the response from the current
// Node is set as the data for the next Node
// It will save the history of the current node in the HistoryManager
func (s *Signal) PrepareForNext() {
	if s.History != nil {
		s.History.AddHistory(*s)
	}
	s.Data = s.Response
	s.Response = nil
}

// When sending to the next Node we set the NodeID of the target Node
func (s *Signal) ChangeTarget(targetID string) {
	s.NodeID = targetID
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
		NodeID:  id,
		Context: cm,
		History: hx,
		Meta:    meta,
		Data:    task,
	}
}
