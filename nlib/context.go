package nlib

import (
	"fmt"
	"sync"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.ContextManager = (*SimpleContextManager)(nil)

type SimpleContextManager struct {
	context map[string]node.DataCarrier
	m       sync.RWMutex
}

// SimpleContextManager manages the contextual information within a signal.
// It updates the context based on metadata or predefined logic, ensuring that
// each node has access to the relevant context during processing.
func NewSimpleContextManager() *SimpleContextManager {
	return &SimpleContextManager{context: make(map[string]node.DataCarrier)}
}

func (c *SimpleContextManager) SetContext(id string, data node.DataCarrier) {
	c.m.Lock()
	defer c.m.Unlock()
	c.context[id] = data
}

func (c *SimpleContextManager) RemoveContext(id string) {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.context, id)
}

func (c *SimpleContextManager) GetContext(id string) (node.DataCarrier, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	data, ok := c.context[id]
	if !ok {
		return nil, fmt.Errorf("Not found")
	}
	return data, nil

}
