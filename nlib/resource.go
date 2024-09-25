package nlib

import (
	"fmt"
	"time"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.ResourceManager = (*SimpleResourceManager)(nil)

// SimpleResourceManager implements rate limiting to control the number of requests
// a node can process within a specified time. It ensures efficient use of system
// resources and prevents overload by throttling excessive requests.
type SimpleResourceManager struct {
	maxRequestsPerSecond int
	requests             chan struct{}
}

func NewSimpleResourceManager(maxRequestsPerSecond int) *SimpleResourceManager {
	manager := &SimpleResourceManager{
		maxRequestsPerSecond: maxRequestsPerSecond,
		requests:             make(chan struct{}, maxRequestsPerSecond),
	}

	// Refill the requests channel every second
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			for i := 0; i < maxRequestsPerSecond; i++ {
				manager.requests <- struct{}{}
			}
		}
	}()

	return manager
}

func (r *SimpleResourceManager) RateLimit(signal node.Signal) error {
	select {
	case <-r.requests:
		// Continue processing
		return nil
	default:
		return fmt.Errorf("rate limit exceeded for Node %s", signal.NodeID)
	}
}
