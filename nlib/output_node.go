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

func NewOutputStringNode(w io.Writer, l node.Logger, sm node.StateManager, name string) *OutputStringNode {
	n := OutputStringNode{w: w}
	n.SetID(name)
	n.SetLogger(l)
	n.SetStateManager(sm)
	n.MakeInputCh()

	go func() {
		for {
			select {
			case sig := <-n.inCh:
				sig = n.PreProcessSignal(sig)

				_, err := n.w.Write([]byte(sig.Data.String()))
				if err != nil {
					n.LogErr(err)
				}

				n.PostProcesSignal(sig)

			case <-n.doneCh:
				return
			}
		}
	}()

	return &n
}
