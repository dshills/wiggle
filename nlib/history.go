package nlib

import "github.com/dshills/wiggle/node"

// Compile-time check
var _ node.HistoryManager = (*SimpleHistoryManager)(nil)

type SimpleHistoryManager struct {
	signals []node.Signal
}

func NewSimpleHistoryManager() *SimpleHistoryManager {
	return &SimpleHistoryManager{}
}

func (hx *SimpleHistoryManager) AddHistory(sig node.Signal) {
	hx.signals = append(hx.signals, sig)
}

func (hx *SimpleHistoryManager) CompressHistory() error {
	return nil
}

func (hx *SimpleHistoryManager) GetHistory() []node.Signal {
	return hx.signals
}

func (hx *SimpleHistoryManager) GetHistoryByID(id string) ([]node.Signal, error) {
	sigList := []node.Signal{}
	for _, sig := range hx.signals {
		if sig.NodeID == id {
			sigList = append(sigList, sig)
		}
	}
	return sigList, nil
}
