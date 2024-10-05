package nmock

import (
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/mock"
)

// Compile-time check
var _ node.Guidance = (*MockGuidance)(nil)
var _ node.ResourceManager = (*MockResourceManager)(nil)
var _ node.Hooks = (*MockHooks)(nil)
var _ node.StateManager = (*MockStateManager)(nil)

// MockGuidance is a testing mock for Guidance
type MockGuidance struct {
	mock.Mock
}

func (m *MockGuidance) Generate(sig node.Signal, context string) (node.Signal, error) {
	args := m.Called(sig, context)
	return args.Get(0).(node.Signal), args.Error(1)
}

// MookHooks is a testing mock for Hooks
type MockHooks struct {
	mock.Mock
}

func (m *MockHooks) BeforeAction(sig node.Signal) (node.Signal, error) {
	args := m.Called(sig)
	return args.Get(0).(node.Signal), args.Error(1)
}

func (m *MockHooks) AfterAction(sig node.Signal) (node.Signal, error) {
	args := m.Called(sig)
	return args.Get(0).(node.Signal), args.Error(1)
}

// MockResourceManager is a testing mock for resource manager
type MockResourceManager struct {
	mock.Mock
}

func (m *MockResourceManager) RateLimit(sig node.Signal) error {
	args := m.Called(sig)
	return args.Error(0)
}

// MockStateManager is a testing mock for StateManager, providing mock behavior
// for methods such as logging, state updates, resource management, coordination, context, and history management.
type MockStateManager struct {
	mock.Mock
}

// SetLogger sets a mock logger for the StateManager.
func (m *MockStateManager) SetLogger(l node.Logger) {
	m.Called(l)
}

// Log records a log message. This is used to mock logging behavior.
func (m *MockStateManager) Log(message string) {
	m.Called(message)
}

// UpdateState mocks the behavior of updating the state of a signal.
func (m *MockStateManager) UpdateState(sig node.Signal) {
	m.Called(sig)
}

// Complete marks the StateManager's work as done. This is a mock of the completion behavior.
func (m *MockStateManager) Complete() {
	m.Called()
}

// ResourceManager returns the mock ResourceManager for managing resource constraints.
func (m *MockStateManager) ResourceManager() node.ResourceManager {
	args := m.Called()
	return args.Get(0).(node.ResourceManager)
}

// SetResourceManager sets the ResourceManager for resource control in this mock.
func (m *MockStateManager) SetResourceManager(mgr node.ResourceManager) {
	m.Called(mgr)
}

// Coordinator returns the mock Coordinator for managing node coordination.
func (m *MockStateManager) Coordinator() node.Coordinator {
	args := m.Called()
	return args.Get(0).(node.Coordinator)
}

// Register mocks the registration of a node and returns a channel for synchronization.
func (m *MockStateManager) Register() chan struct{} {
	args := m.Called()
	return args.Get(0).(chan struct{})
}

// SetCoordinator sets the mock Coordinator in the StateManager.
func (m *MockStateManager) SetCoordinator(cor node.Coordinator) {
	m.Called(cor)
}

// WaitFor mocks waiting for another node to complete its work.
func (m *MockStateManager) WaitFor(n node.Node) {
	m.Called(n)
}

// GetState returns the state of a signal in the mock StateManager.
func (m *MockStateManager) GetState(sig node.Signal) node.State {
	args := m.Called(sig)
	return args.Get(0).(node.State)
}

// ContextManager returns the mock ContextManager for managing signal contexts.
func (m *MockStateManager) ContextManager() node.ContextManager {
	args := m.Called()
	return args.Get(0).(node.ContextManager)
}

// HistoryManager returns the mock HistoryManager for managing signal history.
func (m *MockStateManager) HistoryManager() node.HistoryManager {
	args := m.Called()
	return args.Get(0).(node.HistoryManager)
}

// Logger returns the mock Logger associated with the StateManager.
func (m *MockStateManager) Logger() node.Logger {
	args := m.Called()
	return args.Get(0).(node.Logger)
}

// SetContextManager sets the mock ContextManager for managing contexts.
func (m *MockStateManager) SetContextManager(ctxMgr node.ContextManager) {
	m.Called(ctxMgr)
}

// SetHistoryManager sets the mock HistoryManager for managing signal history.
func (m *MockStateManager) SetHistoryManager(histMgr node.HistoryManager) {
	m.Called(histMgr)
}

// GetContext retrieves the context for a given key. Mocks the context lookup behavior.
func (m *MockStateManager) GetContext(key string) (node.DataCarrier, error) {
	args := m.Called(key)
	return args.Get(0).(node.DataCarrier), args.Error(1)
}

// SetContext sets a context value for a given key. Mocks the behavior of setting context.
func (m *MockStateManager) SetContext(key string, data node.DataCarrier) {
	m.Called(key, data)
}

// RemoveContext removes a context value for a given key. Mocks the behavior of removing context.
func (m *MockStateManager) RemoveContext(key string) {
	m.Called(key)
}

// AddHistory adds a signal to the history. Mocks the behavior of adding history.
func (m *MockStateManager) AddHistory(sig node.Signal) {
	m.Called(sig)
}

// GetHistory retrieves the history of signals. Mocks the history retrieval behavior.
func (m *MockStateManager) GetHistory() []node.Signal {
	args := m.Called()
	return args.Get(0).([]node.Signal)
}

// FilterHistory filters the history of signals by node ID. Mocks the filtering behavior.
func (m *MockStateManager) FilterHistory(nodeid string) []node.Signal {
	args := m.Called(nodeid)
	return args.Get(0).([]node.Signal)
}
