package nmock

import (
	"github.com/dshills/wiggle/node"
	"github.com/stretchr/testify/mock"
)

// Compile-time check
var _ node.ResourceManager = (*MockResourceManager)(nil)

// MockResourceManager is a testing mock for resource manager
type MockResourceManager struct {
	mock.Mock
}

func (m *MockResourceManager) RateLimit(sig node.Signal) error {
	args := m.Called(sig)
	return args.Error(0)
}
