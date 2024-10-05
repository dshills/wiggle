package nlib

import (
	"context"
	"io"
	"time"

	"github.com/dshills/wiggle/node"
)

// SimpleStringReaderNode is an InputNode that reads string data from an io.Reader
// and processes it within the Wiggle framework. It extends EmptyNode for base functionality
// and interacts with the StateManager for signal handling and node lifecycle events.
var _ node.InputNode = (*SimpleStringReaderNode)(nil)

// SimpleStringReaderNode represents a node that reads data from an io.Reader and processes it as string input.
type SimpleStringReaderNode struct {
	EmptyNode
	reader io.Reader
}

// NewSimpleStringReaderNode creates a new instance of SimpleStringReaderNode. It sets up the reader,
// options, and StateManager, and initializes the input channel for signal reception. A go routine
// is launched to process incoming signals or to handle node termination via StateManager.
func NewSimpleStringReaderNode(r io.Reader, mgr node.StateManager, options node.Options) *SimpleStringReaderNode {
	n := SimpleStringReaderNode{reader: r}
	n.SetOptions(options)
	n.SetStateManager(mgr)
	n.MakeInputCh()

	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.LogInfo("Received Signal")
				n.processSignal(sig)
			case <-n.StateManager().Register():
				n.LogInfo("Received Done")
				return
			}
		}
	}()

	return &n
}

// processSignal handles the processing of incoming signals. It first applies signal preprocessing,
// reads the string data from the node's io.Reader, and then sends the processed signal to connected
// nodes. It ensures that any error encountered during reading or processing is logged and that the signal
// is marked as failed. A 2-second timeout is applied during the sending of the signal to prevent blocking.
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
	sig.Result = &Carrier{TextData: string(byts)}

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

// SetReader allows updating the io.Reader for the SimpleStringReaderNode at runtime.
// It replaces the current reader with a new one.
func (n *SimpleStringReaderNode) SetReader(r io.Reader) {
	n.reader = r
}
