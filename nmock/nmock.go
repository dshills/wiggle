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

func (m *MockGuidance) Generate(sig node.Signal) (node.Signal, error) {
	args := m.Called(sig)
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

// MockStateManager is a testing mock for StateManager
type MockStateManager struct {
	mock.Mock
	logger      node.Logger
	coordinator node.Coordinator
	resource    node.ResourceManager
}

func (m *MockStateManager) SetLogger(l node.Logger) {
	m.Called(l)
	m.logger = l
}

func (m *MockStateManager) Log(message string) {
	m.Called(message)
}

func (m *MockStateManager) UpdateState(sig node.Signal) {
	m.Called(sig)
}

func (m *MockStateManager) GetState(node.Signal) node.State {
	return node.State{}
}

func (m *MockStateManager) Complete() {
	m.Called()
}

func (m *MockStateManager) ResourceManager() node.ResourceManager {
	args := m.Called()
	return args.Get(0).(node.ResourceManager)
}

func (m *MockStateManager) SetResourceManager(mgr node.ResourceManager) {
	m.resource = mgr
}

func (m *MockStateManager) Coordinator() node.Coordinator {
	return m.coordinator
}

func (m *MockStateManager) Register() chan struct{} {
	return nil
}

func (m *MockStateManager) SetCoordinator(cor node.Coordinator) {
	m.coordinator = cor
}

func (m *MockStateManager) WaitFor(node.Node) {
}

// MockResourceManager is a testing mock for resource manager
type MockResourceManager struct {
	mock.Mock
}

func (m *MockResourceManager) RateLimit(sig node.Signal) error {
	args := m.Called(sig)
	return args.Error(0)
}
