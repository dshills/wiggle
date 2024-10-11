package nmock

import (
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/mock"
)

// Compile-time check
var _ node.Guidance = (*MockGuidance)(nil)

// MockGuidance is a testing mock for Guidance
type MockGuidance struct {
	mock.Mock
}

func (m *MockGuidance) Generate(sig node.Signal, context string) (node.Signal, error) {
	args := m.Called(sig, context)
	return args.Get(0).(node.Signal), args.Error(1)
}
