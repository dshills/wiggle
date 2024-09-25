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
	errList    []error
	errStrings []string
	doneChs    []chan struct{}
	nodeWaitID string
}

func NewSimpleStateManager() *SimpleStateManager {
	sm := SimpleStateManager{stateMap: make(map[string]node.State)}
	return &sm
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

	if s.nodeWaitID == sig.NodeID {
		s.Complete()
	}
}

func (s *SimpleStateManager) GetState(signal node.Signal) node.State {
	s.mu.Lock()
	defer s.mu.Unlock()
	if state, exists := s.stateMap[signal.NodeID]; exists {
		return state
	}
	return node.State{Status: "unknown"}
}

func (s *SimpleStateManager) Register() chan struct{} {
	ch := make(chan struct{})
	s.doneChs = append(s.doneChs, ch)
	return ch
}

func (s *SimpleStateManager) ShouldFail(err error) bool {
	for _, e := range s.errList {
		if errors.Is(err, e) {
			return true
		}
	}
	for _, es := range s.errStrings {
		if es == err.Error() {
			return true
		}
	}
	return false
}

func (s *SimpleStateManager) Complete() {
	for _, ch := range s.doneChs {
		ch <- struct{}{}
	}
}

func (s *SimpleStateManager) WaitFor(nodeid string) {
	s.nodeWaitID = nodeid
	done := s.Register()
	<-done
}
