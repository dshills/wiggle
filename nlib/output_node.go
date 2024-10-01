package nlib

import (
	"context"
	"io"
	"time"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure OutputStringNode implements the node.OutputNode interface
var _ node.OutputNode = (*OutputStringNode)(nil)

// OutputStringNode writes the signal's response data to an io.Writer.
// This node is designed to output the data of a signal to a specified writer (e.g., file, stdout).
type OutputStringNode struct {
	EmptyNode           // Inherits base node functionality.
	writer    io.Writer // The writer where the signal's response will be written.
}

// NewOutputStringNode creates a new OutputStringNode with the specified writer, logger, state manager, and name.
// It starts a goroutine to listen for incoming signals and write their data to the writer.
func NewOutputStringNode(w io.Writer, l node.Logger, sm node.StateManager, name string) *OutputStringNode {
	n := OutputStringNode{writer: w} // Initialize with the provided writer.
	n.Init(l, sm, name)              // Initialize the node with logger, state manager, and name.

	// Goroutine to listen for incoming signals and process them.
	go func() {
		var err error
		for {
			select {
			case sig := <-n.InputCh(): // Receive a signal from the input channel.
				sig, err = n.PreProcessSignal(sig)
				if err != nil {
					n.Fail(sig, err)
					return
				}

				sig.Status = StatusInProcess
				// Write the signal's data (response) to the provided writer.
				_, err := n.writer.Write([]byte(sig.Task.String()))
				if err != nil {
					n.LogErr(err) // Log any errors encountered during writing.
					sig.Err = err.Error()
					sig.Status = StatusFail
					return
				}
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

			case <-n.DoneCh(): // If the done channel is closed, exit the loop.
				return
			}
		}
	}()

	return &n
}

func (n *OutputStringNode) SetWriter(w io.Writer) {
	n.writer = w
}
