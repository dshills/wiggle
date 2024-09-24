package node

type SimpleContextManager struct {
	defaultContext string
}

// SimpleContextManager manages the contextual information within a signal.
// It updates the context based on metadata or predefined logic, ensuring that
// each node has access to the relevant context during processing.
func NewSimpleContextManager(defaultContext string) *SimpleContextManager {
	return &SimpleContextManager{defaultContext: defaultContext}
}

func (c *SimpleContextManager) UpdateContext(signal Signal) (string, error) {
	// Update context based on signal's metadata
	for _, meta := range signal.Meta {
		if meta.Key == "context-update" {
			signal.Context = meta.Value
			return signal.Context, nil
		}
	}
	// If no update is found, use the default context
	return c.defaultContext, nil
}

func (c *SimpleContextManager) GetContext(signal Signal) string {
	// Return the current context of the signal
	if signal.Context == "" {
		return c.defaultContext
	}
	return signal.Context
}
