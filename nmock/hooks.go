package nmock

import (
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/mock"
)

// Compile-time check
var _ node.Hooks = (*MockHooks)(nil)

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
