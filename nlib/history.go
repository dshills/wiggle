package nlib

import "github.com/dshills/wiggle/node"

// Compile-time check to ensure SimpleHistoryManager implements the node.HistoryManager interface.
var _ node.HistoryManager = (*SimpleHistoryManager)(nil)

// SimpleHistoryManager is a basic implementation of the HistoryManager interface.
// It stores a list of signals, keeping track of the signal history as they are processed by nodes.
type SimpleHistoryManager struct {
	signals []node.Signal // Slice to store the history of signals
}

// NewSimpleHistoryManager initializes and returns a new instance of SimpleHistoryManager.
func NewSimpleHistoryManager() *SimpleHistoryManager {
	return &SimpleHistoryManager{}
}

// AddHistory appends the given signal to the list of signal history.
func (hx *SimpleHistoryManager) AddHistory(sig node.Signal) {
	hx.signals = append(hx.signals, sig) // Add the signal to the history slice
}

// CompressHistory currently does nothing, but could be extended to implement history compression.
func (hx *SimpleHistoryManager) CompressHistory() error {
	return nil // Placeholder for future implementation
}

// GetHistory returns the full list of signals in the history.
func (hx *SimpleHistoryManager) GetHistory() []node.Signal {
	return hx.signals // Return the complete signal history
}

// GetHistoryByID returns the signals that match the provided Node ID.
// It filters through the history and collects signals with the specified ID.
func (hx *SimpleHistoryManager) Filter(nodeid string) []node.Signal {
	sigList := []node.Signal{}
	for _, sig := range hx.signals {
		if sig.NodeID == nodeid {
			sigList = append(sigList, sig)
		}
	}
	return sigList
}
