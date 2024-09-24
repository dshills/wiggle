package node

import "sync"

// SimpleStateManager tracks and updates the state of a signal as it moves through the node chain.
// It ensures that nodes can access and update the signal's current state, allowing for
// consistent state management across complex workflows.
type SimpleStateManager struct {
	stateMap map[string]string // Mapping of NodeID to state
	mu       sync.Mutex
}

func NewSimpleStateManager() *SimpleStateManager {
	return &SimpleStateManager{
		stateMap: make(map[string]string),
	}
}

func (s *SimpleStateManager) UpdateState(signal Signal, state string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stateMap[signal.NodeID] = state
	return nil
}

func (s *SimpleStateManager) GetState(signal Signal) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if state, exists := s.stateMap[signal.NodeID]; exists {
		return state
	}
	return "unknown"
}
