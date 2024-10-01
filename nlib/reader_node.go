package nlib

import (
	"context"
	"io"
	"time"

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
			case <-n.DoneCh():
				return
			}
		}
	}()

	return &n
}

func (n *SimpleStringReaderNode) processSignal(sig node.Signal) {
	var err error
	sig, err = n.PreProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	sig.Status = StatusInProcess
	byts, err := io.ReadAll(n.reader)
	if err != nil {
		n.LogErr(err)
		sig.Err = err.Error()
		return
	}
	sig.Result = NewStringData(string(byts))

	sig.Status = StatusSuccess

	sig, err = n.PostProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure the context is cancelled once we're done

	if err := n.SendToConnected(ctx, sig); err != nil {
		n.Fail(sig, err)
		return
	}
}

func (n *SimpleStringReaderNode) SetReader(r io.Reader) {
	n.reader = r
}
