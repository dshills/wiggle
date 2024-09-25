package nlib

import (
	"errors"
	"sync"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.StateManager = (*SimpleStateManager)(nil)

// SimpleStateManager tracks and updates the state of a signal as it moves through the node chain.
// It ensures that nodes can access and update the signal's current state, allowing for
// consistent state management across complex workflows.
type SimpleStateManager struct {
	stateMap   map[string]node.State // Mapping of NodeID to state
	mu         sync.Mutex
	fail       chan struct{}
	errList    []error
	errStrings []string
	doneChs    []chan struct{}
}

func NewSimpleStateManager() *SimpleStateManager {
	return &SimpleStateManager{
		stateMap: make(map[string]node.State),
	}
}

func (s *SimpleStateManager) UpdateState(sig node.Signal) {
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.stateMap[sig.NodeID]
	if !ok {
		st = node.State{}
	}
	st.Completed++
	if sig.Err != nil {
		st.Failures++
	}
	st.Status = sig.Status
	s.stateMap[sig.NodeID] = st
}

func (s *SimpleStateManager) GetState(signal node.Signal) node.State {
	s.mu.Lock()
	defer s.mu.Unlock()
	if state, exists := s.stateMap[signal.NodeID]; exists {
		return state
	}
	return node.State{Status: "unknown"}
}

func (nm *SimpleStateManager) Register() chan struct{} {
	ch := make(chan struct{})
	nm.doneChs = append(nm.doneChs, ch)
	return ch
}

func (nm *SimpleStateManager) ShouldFail(err error) bool {
	for _, e := range nm.errList {
		if errors.Is(err, e) {
			return true
		}
	}
	for _, es := range nm.errStrings {
		if es == err.Error() {
			return true
		}
	}
	return false
}

func (nm *SimpleStateManager) Complete() {
	for _, ch := range nm.doneChs {
		ch <- struct{}{}
	}
}
