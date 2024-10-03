package nmock

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dshills/wiggle/node"
)

// MockStateManager is a mock implementation of the node.StateManager interface
type MockStateManager struct {
	mu            sync.Mutex
	state         map[string]node.State
	doneCh        chan struct{}
	waitForCalled map[string]bool
}

func NewMockStateManager() *MockStateManager {
	return &MockStateManager{
		state:         make(map[string]node.State),
		doneCh:        make(chan struct{}),
		waitForCalled: make(map[string]bool),
	}
}

// Complete is a mock method for the StateManager interface, does nothing
func (m *MockStateManager) Complete() {
	// No-op for test
}

// GetState returns the state of the given signal (mock behavior)
func (m *MockStateManager) GetState(sig node.Signal) node.State {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state[sig.NodeID]
}

// Register returns a channel used to signal when registration is complete
func (m *MockStateManager) Register() chan struct{} {
	return m.doneCh
}

// UpdateState updates the mock state for a given signal
func (m *MockStateManager) UpdateState(sig node.Signal) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state[sig.NodeID] = node.State{ /* Fill this out as needed */ }
}

// WaitFor simulates waiting for a node (used to verify if this method is called)
func (m *MockStateManager) WaitFor(n node.Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.waitForCalled[n.ID()] = true
}

// CheckWaitForCalled checks if WaitFor was called for a specific node ID
func (m *MockStateManager) CheckWaitForCalled(nodeID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.waitForCalled[nodeID]
}

type MockLogger struct {
	entries []string
}

func (l *MockLogger) Log(msg string) {
	l.entries = append(l.entries, msg)
}

func (l *MockLogger) Entries() []string {
	return l.entries
}

var ErrMockRetry = fmt.Errorf("mock error")

type MockErrorGuidance struct {
	RetriesVal int
	RetryErr   error
}

func (eg *MockErrorGuidance) Retries() int {
	return eg.RetriesVal
}

func (eg *MockErrorGuidance) Action(err error) node.ErrGuide {
	if errors.Is(err, ErrMockRetry) {
		return node.ErrGuideRetry
	}
	return node.ErrGuideFail
}
