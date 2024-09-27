package nlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure InteractiveNode implements the node.Node interface
var _ node.Node = (*InteractiveNode)(nil)

type InteractiveNode struct {
	EmptyNode
}

func NewInteractiveNode(l node.Logger, sm node.StateManager, name string) *InteractiveNode {
	n := InteractiveNode{}
	n.Init(l, sm, name)
	go n.listen()
	return &n
}

// listen listens for incoming signals on the input channel and processes them.
func (n *InteractiveNode) listen() {
	for {
		select {
		case sig := <-n.inCh: // Wait for an incoming signal
			n.LogInfo("Received signal") // Log the receipt of a signal

			sig = n.PreProcessSignal(sig)

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nEnter your question (type /quit to stop): ")
			query, _ := reader.ReadString('\n')

			query = strings.TrimSpace(query)

			if query == "/quit" {
				n.stateMgr.Complete()
			}

			sig.Response = NewStringData(query)

			sig = n.PostProcesSignal(sig)
			n.SendToConnected(sig)

		case <-n.doneCh: // If the done channel is closed, exit the function
			return
		}
	}
}
