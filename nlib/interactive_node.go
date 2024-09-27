package nlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dshills/wiggle/node"
)

type InteractiveNode struct {
	EmptyNode
}

func NewInteractiveNode(l node.Logger, sm node.StateManager, name string) *InteractiveNode {
	n := InteractiveNode{}
	n.SetLogger(l)
	n.SetStateManager(sm)
	n.SetID(name)
	n.MakeInputCh()
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

			n.PostProcesSignal(sig)

		case <-n.doneCh: // If the done channel is closed, exit the function
			return
		}
	}
}
