package nlib

import (
	"context"
	"fmt"
	"time"

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
func NewAINode(lm llm.LLM, sm node.StateManager, options node.Options) node.Node {
	n := AINode{lm: lm} // Initialize the AINode with the provided LLM
	n.SetOptions(options)
	n.SetStateManager(sm)
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

// listen listens for incoming signals on the input channel and processes them.
func (n *AINode) processSignal(sig node.Signal) {
	start := time.Now()
	var err error
	var ctx = context.TODO() // Initialize a context for managing requests
	sig, err = n.PreProcessSignal(sig)
	if err != nil {
		n.LogErr(err)
		n.StateManager().Complete()
		return
	}

	sig.Status = StatusInProcess

	// Generate guidance (possibly modify the signal) before sending it to the LLM
	if guide := n.Guidance(); guide != nil {
		sig, err = guide.Generate(sig)
		if err != nil {
			n.LogErr(err)
		}
	}

	n.LogInfo(fmt.Sprintf("Sending to llm %s", n.lm.Model())) // Log the data being sent to the LLM

	// Call the LLM to process the signal
	sig, err = n.CallLLM(ctx, sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}
	sig.Status = StatusSuccess

	sig, err = n.PostProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	if err := n.SendToConnected(ctx, sig); err != nil {
		n.Fail(sig, err)
		return
	}
	n.LogInfo(fmt.Sprintf("%v completed in %v", n.lm.Model(), time.Since(start)))
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
