package nlib

import (
	"io"

	"github.com/dshills/wiggle/node"
)

// OutputStringNode simply writes the Signal Response to a writer
type OutputStringNode struct {
	EmptyNode
	w io.Writer
}

func NewOutputStringNode(w io.Writer, l node.Logger, sm node.StateManager) *OutputStringNode {
	n := OutputStringNode{w: w}
	n.SetID(generateUUID())
	n.SetLogger(l)
	n.SetStateManager(sm)
	n.MakeInputCh()

	go func() {
		select {
		case sig := <-n.inCh:
			sig = node.SignalFromSignal(n.ID(), sig)
			n.LogInfo("Received signal")
			_, err := n.w.Write([]byte(sig.Data.String()))
			if err != nil {
				n.LogErr(err)
			}
			n.UpdateState(sig)
		case <-n.doneCh:
			return
		}
	}()

	return &n
}
