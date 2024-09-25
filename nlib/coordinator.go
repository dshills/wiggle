package nlib

/*
// Compile-time check
var _ node.Coordinator = (*SimpleCoordinator)(nil)

// SimpleCoordinator manages the execution flow across multiple nodes by sending
// signals through each node's input channel. It synchronizes node processing, waits
// for completion using a timeout mechanism, and handles errors if any node encounters issues.
// This ensures that workflows proceed smoothly or fail gracefully.
type SimpleCoordinator struct {
	timeout time.Duration
}

func NewSimpleCoordinator(timeout time.Duration) *SimpleCoordinator {
	return &SimpleCoordinator{timeout: timeout}
}

func (c *SimpleCoordinator) WaitForCompletion(nodes ...node.Node) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(nodes))

	// Start processing for each node by sending a signal to its input channel
	for _, nd := range nodes {
		wg.Add(1)
		go func(n node.Node) {
			defer wg.Done()

			// Create a signal for this node
			//signal := node.Signal{NodeID: n.ID(), Data: MessageData{Message: "Trigger processing"}}

			// Send the signal to the node's input channel
			select {
			case n.InputCh() <- signal:
				// Wait for processing
				// Assume that nodes communicate any errors through their InputCh processing
				// If node-specific error handling is required, it could be added here
			case <-time.After(c.timeout):
				errCh <- fmt.Errorf("timeout on node %s", n.ID())
			}
		}(nd)
	}

	// Wait for all nodes to finish
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	// Handle timeout or errors
	select {
	case <-doneCh:
		// All nodes completed successfully
		return nil
	case err := <-errCh:
		// Return the first error that occurred
		return err
	case <-time.After(c.timeout):
		// Timeout exceeded
		return fmt.Errorf("timeout exceeded while waiting for nodes to complete")
	}
}

func (c *SimpleCoordinator) CancelOnTimeout(duration time.Duration) {
	c.timeout = duration
}

// MessageData implements the DataCarrier interface for carrying simple string messages.
type MessageData struct {
	Message string
}

// ToMessageList returns a message list with the content of the MessageData.
func (m MessageData) ToMessageList() llm.MessageList {
	return llm.MessageList{llm.Message{Content: m.Message}}
}

// ToJSON converts the message data to a JSON string.
func (m MessageData) ToJSON() string {
	return fmt.Sprintf(`{"message": "%s"}`, m.Message)
}

// ToVector is not applicable for MessageData, so we return nil here.
func (m MessageData) ToVector() []float32 {
	return nil
}
*/
