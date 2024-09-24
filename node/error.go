package node

import "fmt"

// SimpleErrorHandler logs errors encountered in nodes and checks for severity.
// Based on the error, it either allows the workflow to continue or halts it
// if the error is considered severe. This helps manage fault tolerance in workflows.
type SimpleErrorHandler struct {
	severeErrors []string // list of severe errors that should halt the workflow
	logger       Logger
}

func NewSimpleErrorHandler(severeErrors []string, logger Logger) *SimpleErrorHandler {
	return &SimpleErrorHandler{
		severeErrors: severeErrors,
		logger:       logger,
	}
}

func (h *SimpleErrorHandler) HandleError(signal Signal, err error) bool {
	h.logger.Log(fmt.Sprintf("Error in Node %s: %v", signal.NodeID, err))

	// Check if the error is severe and halt the workflow
	for _, severeError := range h.severeErrors {
		if err.Error() == severeError {
			h.logger.Log("Severe error encountered, halting workflow.")
			return false
		}
	}

	// Otherwise, continue processing
	return true
}
