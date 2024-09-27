package nlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure InteractiveNode implements the node.Node interface.
var _ node.Node = (*InteractiveNode)(nil)

// InteractiveNode is a node that interacts with the user via the command line.
// It waits for user input and processes the user's query as part of the signal.
type InteractiveNode struct {
	EmptyNode // Inherits base node functionality.
}

// NewInteractiveNode creates a new InteractiveNode and starts listening for signals.
// - l: Logger for logging interactions.
// - sm: StateManager for managing the state of the node.
// - name: Name/ID of the node.
func NewInteractiveNode(l node.Logger, sm node.StateManager, name string) *InteractiveNode {
	n := InteractiveNode{}
	n.Init(l, sm, name) // Initialize the node with logger, state manager, and name.

	go n.listen() // Start the signal listener in a separate goroutine.
	return &n
}

// listen listens for incoming signals and processes them by interacting with the user.
// It prompts the user to enter a query and sends that query as the signal's response.
func (n *InteractiveNode) listen() {
	for {
		select {
		case sig := <-n.inCh: // Wait for an incoming signal on the input channel.
			n.LogInfo("Received signal") // Log that a signal has been received.

			sig = n.PreProcessSignal(sig) // Run pre-processing hooks on the signal.

			// Prompt the user for input.
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nEnter your question (type /quit to stop): ")
			query, _ := reader.ReadString('\n')

			// Trim any extra whitespace or newline characters.
			query = strings.TrimSpace(query)

			// Check if the user entered the quit command.
			if query == "/quit" {
				n.stateMgr.Complete() // Mark the node as complete if the user quits.
				return
			}

			// Set the user's input as the signal's response data.
			sig.Response = NewStringData(query)

			// Run post-processing hooks and forward the signal to connected nodes.
			sig = n.PostProcesSignal(sig)
			n.SendToConnected(sig)

		case <-n.doneCh: // Exit the function if the done channel is closed.
			return
		}
	}
}
