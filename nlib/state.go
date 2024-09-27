package nlib

import (
	"errors"
	"sync"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure SimpleStateManager implements the node.StateManager interface
var _ node.StateManager = (*SimpleStateManager)(nil)

// SimpleStateManager is a basic implementation of the StateManager interface.
// It manages the state of signals and tracks errors and completion for nodes.
type SimpleStateManager struct {
	stateMap   map[string]node.State // Maps NodeID to their respective state
	mu         sync.Mutex            // Mutex to ensure safe concurrent access to stateMap
	errList    []error               // List of known errors for comparison
	errStrings []string              // List of known error strings for comparison
	doneChs    []chan struct{}       // Channels for nodes to signal completion
	nodeWaitID string                // ID of the node being waited on for completion
	waitCh     chan struct{}         // Channel to wait on
}

// NewSimpleStateManager creates and returns a new instance of SimpleStateManager.
func NewSimpleStateManager() *SimpleStateManager {
	sm := SimpleStateManager{stateMap: make(map[string]node.State)}
	return &sm
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

// Register creates a channel for signaling completion and adds it to the list of done channels.
func (s *SimpleStateManager) Register() chan struct{} {
	ch := make(chan struct{})         // Create a new completion channel
	s.doneChs = append(s.doneChs, ch) // Add it to the list of done channels
	return ch                         // Return the new channel
}

// ShouldFail determines if an error should cause failure, checking against
// known error types and error strings.
func (s *SimpleStateManager) ShouldFail(err error) bool {
	for _, e := range s.errList { // Check if the error matches any known errors
		if errors.Is(err, e) {
			return true
		}
	}
	for _, es := range s.errStrings { // Check if the error string matches any known error strings
		if es == err.Error() {
			return true
		}
	}
	return false
}

// Complete signals completion to all registered channels.
func (s *SimpleStateManager) Complete() {
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
	if n == nil {
		s.waitCh = make(chan struct{})
		<-s.waitCh
		return
	}

	s.nodeWaitID = n.ID() // Store the NodeID to wait on
	done := s.Register()  // Register the channel for the node
	<-done                // Wait until the completion signal is received
}
