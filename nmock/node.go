package nmock

import (
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/mock"
)

// Compile-time check
var _ node.Node = (*MockNode)(nil)

// nmock.MockNode
type MockNode struct {
	mock.Mock
	node.Node
}

func (m *MockNode) ID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockNode) InputCh() chan node.Signal {
	args := m.Called()
	return args.Get(0).(chan node.Signal)
}

func (m *MockNode) Connect(nodes ...node.Node) {
	m.Called(nodes)
}

func (m *MockNode) SetID(id string) {
	m.Called(id)
}

func (m *MockNode) SetOptions(options node.Options) {
	m.Called(options)
}

func (m *MockNode) SetStateManager(stateManager node.StateManager) {
	m.Called(stateManager)
}
