package nlib

import (
	"context"
	"fmt"
	"time"

	"github.com/dshills/wiggle/llm"
	"github.com/dshills/wiggle/node"
)

// AINode implements the node.Node interface and integrates with a large language model (LLM).
// It is responsible for processing incoming signals, sending them to the LLM, and handling responses.
var _ node.Node = (*AINode)(nil)

// AINode represents a node that uses a large language model (LLM) to process signals.
// It embeds EmptyNode for base functionality and integrates with the LLM through the lm field.
type AINode struct {
	EmptyNode         // Provides base node functionality like logging, state management, etc.
	lm        llm.LLM // The large language model (LLM) used for processing the node's signals
}

// NewAINode creates a new AINode with the specified LLM, state manager, and options.
// It sets up the node by configuring options, state management, and input channel.
// A goroutine is started to listen for incoming signals and process them using the LLM.
func NewAINode(lm llm.LLM, sm node.StateManager, options node.Options) node.Node {
	n := AINode{lm: lm} // Initialize the AINode with the provided LLM
	n.SetOptions(options)
	n.SetStateManager(sm)
	n.MakeInputCh()

	// Start a goroutine to listen for signals or handle termination events
	go func() {
		for {
			select {
			case sig := <-n.InputCh():
				n.LogInfo("Received Signal")
				n.processSignal(sig) // Process the received signal
			case <-n.StateManager().Register():
				n.LogInfo("Received Done")
				return // Terminate the goroutine when done
			}
		}
	}()

	return &n
}

// processSignal handles the signal processing for the AINode. It preprocesses the signal,
// sends it to the LLM for processing, and handles the response. If any error occurs during
// processing, the signal is marked as failed. The function also logs the total time taken to process the signal.
func (n *AINode) processSignal(sig node.Signal) {
	start := time.Now()
	var err error
	var ctx = context.TODO()

	// Preprocess the signal before sending it to the LLM
	sig, err = n.PreProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err)
		return
	}

	sig.Status = StatusInProcess // Set signal status to in process

	// Optionally generate guidance (modify the signal) before sending to the LLM
	if guide := n.Guidance(); guide != nil {
		context := ""
		if ctxmgr := n.stateMgr.ContextManager(); ctxmgr != nil {
			data, err := ctxmgr.GetContext(n.ID())
			if err == nil && data != nil {
				context = data.String()
			}
		}
		sig, err = guide.Generate(sig, context)
		if err != nil {
			n.LogErr(err) // Log error in guidance generation
		}
	}

	n.LogInfo(fmt.Sprintf("Sending to llm %s", n.lm.Model())) // Log the LLM model being used

	// Call the LLM to process the signal
	sig, err = n.CallLLM(ctx, sig)
	if err != nil {
		n.Fail(sig, err) // Mark the signal as failed
		return
	}
	sig.Status = StatusSuccess // Mark the signal as successful after LLM processing

	// Postprocess the signal after successful LLM interaction
	sig, err = n.PostProcessSignal(sig)
	if err != nil {
		n.Fail(sig, err) // Mark the signal as failed if postprocessing fails
		return
	}

	// Send the processed signal to connected nodes
	if err := n.SendToConnected(ctx, sig); err != nil {
		n.Fail(sig, err) // Mark the signal as failed if sending fails
		return
	}

	// Log the total time taken for processing the signal
	n.LogInfo(fmt.Sprintf("%v completed in %v", n.lm.Model(), time.Since(start)))
}

// CallLLM sends the signal data to the LLM for processing and returns the modified signal.
// It creates a message list from the signal's task data and sends it to the LLM via its Chat method.
// If successful, the response is stored in the signal's Result field.
func (n *AINode) CallLLM(ctx context.Context, sig node.Signal) (node.Signal, error) {
	// Create a message list with the signal's task data as the user message
	msgList := llm.MessageList{llm.UserMsg(sig.Task.String())}

	// Call the LLM to process the message list and return a response
	msg, err := n.lm.Chat(ctx, msgList)
	if err != nil {
		return sig, err // Return the signal and error if the LLM call fails
	}

	// Set the LLM's response as the result in the signal
	sig.Result = &Carrier{TextData: msg.Content}

	return sig, nil // Return the signal with the LLM's response
}
