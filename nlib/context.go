package nlib

import (
	"fmt"
	"sync"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure SimpleContextManager implements the node.ContextManager interface.
var _ node.ContextManager = (*SimpleContextManager)(nil)

// SimpleContextManager is a basic implementation of the ContextManager interface.
// It manages contextual data for signals, storing and retrieving the relevant data based on node IDs.
type SimpleContextManager struct {
	context map[string]node.DataCarrier // Maps NodeID to its associated contextual data
	m       sync.RWMutex                // Mutex to ensure thread-safe access to the context map
}

// NewSimpleContextManager initializes and returns a new instance of SimpleContextManager.
func NewSimpleContextManager() *SimpleContextManager {
	return &SimpleContextManager{context: make(map[string]node.DataCarrier)} // Initializes the context map
}

// SetContext sets the context for a given NodeID, storing the provided DataCarrier.
func (c *SimpleContextManager) SetContext(id string, data node.DataCarrier) {
	c.m.Lock() // Acquire a write lock to ensure safe modification
	defer c.m.Unlock()
	c.context[id] = data // Store the context data for the given NodeID
}

// RemoveContext removes the context associated with a given NodeID from the map.
func (c *SimpleContextManager) RemoveContext(id string) {
	c.m.Lock() // Acquire a write lock to safely modify the context map
	defer c.m.Unlock()
	delete(c.context, id) // Remove the context data for the specified NodeID
}

// GetContext retrieves the context associated with a given NodeID.
// If no context is found, it returns an error indicating the context is not found.
func (c *SimpleContextManager) GetContext(id string) (node.DataCarrier, error) {
	c.m.RLock() // Acquire a read lock to ensure safe access to the context map
	defer c.m.RUnlock()
	data, ok := c.context[id] // Retrieve the context data for the given NodeID
	if !ok {
		return nil, fmt.Errorf("not found") // Return an error if no context is found for the NodeID
	}
	return data, nil // Return the found context data
}
