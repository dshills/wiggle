package nlib

import (
	"fmt"
	"sync"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure SimpleStateManager implements the node.StateManager interface
var _ node.StateManager = (*SimpleStateManager)(nil)

// SimpleStateManager is a basic implementation of the StateManager interface.
// It manages the state of signals and tracks errors and completion for nodes.
type SimpleStateManager struct {
	stateMap    map[string]node.State
	mu          sync.Mutex
	doneChs     []chan struct{}
	nodeWaitID  string
	waitCh      chan struct{}
	logger      node.Logger
	resourceMgr node.ResourceManager
	coordinator node.Coordinator
	historyMgr  node.HistoryManager
	contextMgr  node.ContextManager
}

// NewSimpleStateManager creates and returns a new instance of SimpleStateManager.
func NewSimpleStateManager(l node.Logger) *SimpleStateManager {
	sm := SimpleStateManager{
		stateMap:   make(map[string]node.State),
		logger:     l,
		historyMgr: NewSimpleHistoryManager(),
		contextMgr: NewSimpleContextManager(),
	}
	return &sm
}

// GetState returns the current state of the specified signal.
// If no state exists for the NodeID, it returns a default state with "unknown" status.
func (s *SimpleStateManager) GetState(signal node.Signal) node.State {
	s.mu.Lock() // Lock to ensure safe access to stateMap
	defer s.mu.Unlock()
	if state, exists := s.stateMap[signal.NodeID]; exists {
		return state // Return the found state
	}
	return node.State{Status: "unknown"} // Return default state if none found
}

// UpdateState updates the state of a signal for the corresponding NodeID.
// It increments the Completed or Failures counters and updates the status.
func (s *SimpleStateManager) UpdateState(sig node.Signal) {
	s.mu.Lock() // Lock to ensure safe modification of stateMap
	defer s.mu.Unlock()
	st, ok := s.stateMap[sig.NodeID] // Get the state associated with the signal's NodeID
	if !ok {
		st = node.State{} // Initialize if no state exists for the NodeID
	}
	st.Completed++     // Increment the completion counter
	if sig.Err != "" { // Increment the failure counter if there's an error
		st.Failures++
	}
	st.Status = sig.Status      // Update the status of the signal
	s.stateMap[sig.NodeID] = st // Store the updated state

	// If waiting on this NodeID, signal completion
	if s.nodeWaitID == sig.NodeID {
		s.Complete()
	}
}

func (s *SimpleStateManager) ContextManager() node.ContextManager {
	return s.contextMgr
}

func (s *SimpleStateManager) Coordinator() node.Coordinator {
	return s.coordinator
}

func (s *SimpleStateManager) HistoryManager() node.HistoryManager {
	return s.historyMgr
}

func (s *SimpleStateManager) Logger() node.Logger {
	return s.logger
}

func (s *SimpleStateManager) ResourceManager() node.ResourceManager {
	return s.resourceMgr
}

// Register creates a channel for signaling completion and adds it to the list of done channels.
func (s *SimpleStateManager) Register() chan struct{} {
	ch := make(chan struct{}, 2)      // Create a new completion channel
	s.doneChs = append(s.doneChs, ch) // Add it to the list of done channels
	return ch                         // Return the new channel
}

func (s *SimpleStateManager) SetContextManager(con node.ContextManager) {
	s.contextMgr = con
}

func (s *SimpleStateManager) SetCoordinator(cor node.Coordinator) {
	s.coordinator = cor
}

func (s *SimpleStateManager) SetHistoryManager(hx node.HistoryManager) {
	s.historyMgr = hx
}

func (s *SimpleStateManager) SetLogger(logger node.Logger) {
	s.logger = logger
}

func (s *SimpleStateManager) SetResourceManager(resMgr node.ResourceManager) {
	s.resourceMgr = resMgr
}

// Complete signals completion to all registered channels.
func (s *SimpleStateManager) Complete() {
	if s.logger != nil {
		s.logger.Log(fmt.Sprintf("{ \"severity\": %q, \"id\": %q, \"msg\": %q }", "info", "STATEMANAGER", "Complete"))
	}
	for _, ch := range s.doneChs { // Iterate through all completion channels
		ch <- struct{}{} // Signal completion
	}
	if s.waitCh != nil {
		s.waitCh <- struct{}{}
	}
}

// WaitFor sets the node to wait for its UpdateState call
// if Node is nil wait forever
func (s *SimpleStateManager) WaitFor(n node.Node) {
	if n != nil {
		s.nodeWaitID = n.ID() // Store the NodeID to wait on
	}
	s.waitCh = make(chan struct{})
	<-s.waitCh
}

func (s *SimpleStateManager) Log(msg string) {
	if s.logger != nil {
		s.logger.Log(msg)
	}
}

func (s *SimpleStateManager) GetContext(key string) (node.DataCarrier, error) {
	if s.contextMgr == nil {
		return nil, fmt.Errorf("missing ContextManager")
	}
	return s.contextMgr.GetContext(key)
}

func (s *SimpleStateManager) SetContext(key string, data node.DataCarrier) {
	if s.contextMgr == nil {
		return
	}
	s.contextMgr.SetContext(key, data)
}

func (s *SimpleStateManager) RemoveContext(key string) {
	if s.contextMgr == nil {
		return
	}
	s.contextMgr.RemoveContext(key)
}

func (s *SimpleStateManager) AddHistory(n node.Signal) {
	if s.historyMgr == nil {
		return
	}
	s.historyMgr.AddHistory(n)
}

func (s *SimpleStateManager) GetHistory() []node.Signal {
	if s.historyMgr == nil {
		return nil
	}
	return s.historyMgr.GetHistory()
}

func (s *SimpleStateManager) FilterHistory(nodeid string) []node.Signal {
	if s.historyMgr == nil {
		return nil
	}
	return s.historyMgr.Filter(nodeid)
}
