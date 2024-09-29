package nlib

import (
	"context"
	"fmt"

	"github.com/dshills/wiggle/llm"
	"github.com/dshills/wiggle/node"
)

// Compile-time check to ensure AINode implements the node.Node interface
var _ node.Node = (*AINode)(nil)

// AINode represents an AI processing node that uses an LLM to process incoming signals.
type AINode struct {
	EmptyNode         // Embeds the base functionality of an EmptyNode
	lm        llm.LLM // llm represents the large language model (LLM) used for processing
}

// NewAINode creates a new AINode, initializing the LLM, logger, and state manager,
// and sets up the node to listen for incoming signals asynchronously.
func NewAINode(lm llm.LLM, l node.Logger, sm node.StateManager, name string) node.Node {
	ai := AINode{lm: lm} // Initialize the AINode with the provided LLM
	ai.Init(l, sm, name)

	go func() {
		for {
			select {
			case sig := <-ai.InputCh():
				ai.processSignal(sig)
			case <-ai.DoneCh():
			}
		}
	}()

	return &ai
}

// listen listens for incoming signals on the input channel and processes them.
func (n *AINode) processSignal(sig node.Signal) {
	var err error
	var ctx = context.TODO() // Initialize a context for managing requests
	sig = n.PreProcessSignal(sig)

	sig.Status = StatusInProcess

	// Generate guidance (possibly modify the signal) before sending it to the LLM
	sig, err = n.GenGuidance(sig)
	if err != nil {
		n.LogErr(err) // Log any errors during guidance generation
		sig.Err = err.Error()
	}

	n.LogInfo(fmt.Sprintf("Sending to llm: %v", sig.Task.String())) // Log the data being sent to the LLM

	// Call the LLM to process the signal
	sig, err = n.CallLLM(ctx, sig)
	if err != nil {
		n.LogErr(err) // Log any errors returned by the LLM
		sig.Err = err.Error()
		sig.Status = StatusFail
		return
	}
	sig.Status = StatusSuccess

	sig = n.PostProcesSignal(sig)
	n.SendToConnected(sig)
}

// callLLM sends the signal's data to the LLM for processing and returns the modified signal.
func (n *AINode) CallLLM(ctx context.Context, sig node.Signal) (node.Signal, error) {
	// Create a list of messages for the LLM, using the signal's data as the user message
	msgList := llm.MessageList{llm.UserMsg(sig.Task.String())}

	// Call the LLM with the message list
	msg, err := n.lm.Chat(ctx, msgList)
	if err != nil {
		return sig, err // Return the error if the LLM call fails
	}

	// Set the response data in the signal
	sig.Result = NewStringData(msg.Content)

	return sig, nil // Return the modified signal
}
