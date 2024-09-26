package nlib

import (
	"errors"
	"sync"
	"time"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure SimpleCoordinator implements the node.Coordinator interface.
var _ node.Coordinator = (*SimpleCoordinator)(nil)

// SimpleCoordinator is responsible for managing the synchronization and execution flow across multiple nodes.
type SimpleCoordinator struct {
	timeout time.Duration // Duration to wait for node completion before timing out
}

// NewSimpleCoordinator creates a new instance of SimpleCoordinator.
func NewSimpleCoordinator() *SimpleCoordinator {
	return &SimpleCoordinator{}
}

// WaitForCompletion waits for all the given nodes to signal their completion.
func (c *SimpleCoordinator) WaitForCompletion(nodes ...node.Node) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(nodes))

	// Add all nodes to the wait group
	wg.Add(len(nodes))

	for _, n := range nodes {
		go func(n node.Node) {
			defer wg.Done()
			select {
			case <-n.InputCh(): // Assuming nodes signal their completion via input channel
				// Node has completed processing
			case <-time.After(c.timeout): // Timeout if the node takes too long
				errCh <- errors.New("timeout waiting for node " + n.ID())
			}
		}(n)
	}

	// Wait for all nodes to complete or return an error on timeout
	wg.Wait()
	close(errCh)

	// Return any error encountered during execution
	if len(errCh) > 0 {
		return <-errCh
	}
	return nil
}

// CancelOnTimeout cancels the waiting process if it exceeds the given timeout duration.
func (c *SimpleCoordinator) CancelOnTimeout(duration time.Duration) {
	c.timeout = duration
}
