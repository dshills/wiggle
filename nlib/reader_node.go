package nlib

import (
	"io"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure SimpleStringReaderNode implements the node.Input interface
var _ node.InputNode = (*SimpleStringReaderNode)(nil)

type SimpleStringReaderNode struct {
	EmptyNode
	reader io.Reader
}

func NewSimpleStringReaderNode(r io.Reader, l node.Logger, sm node.StateManager, name string) *SimpleStringReaderNode {
	n := SimpleStringReaderNode{reader: r}
	n.Init(l, sm, name)

	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.processSignal(sig)
			case <-n.doneCh:
				return
			}
		}
	}()

	return &n
}

func (n *SimpleStringReaderNode) processSignal(sig node.Signal) {
	sig = n.PreProcessSignal(sig)

	byts, err := io.ReadAll(n.reader)
	if err != nil {
		n.LogErr(err)
		return
	}
	sig.Response = NewStringData(string(byts))

	sig = n.PostProcesSignal(sig)
	n.SendToConnected(sig)
}

func (n *SimpleStringReaderNode) SetReader(r io.Reader) {
	n.reader = r
}
