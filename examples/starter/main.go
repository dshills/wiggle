package main

import (
	"log"
	"os"

	"github.com/dshills/wiggle/llm/openai"
	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/node"
)

func main() {
	// Setup LLM
	baseURL := os.Getenv("OPENAI_API_URL")
	apiKey := os.Getenv("OPENAI_API_KEY")
	model := "gpt-4o"
	lm := openai.New(baseURL, model, apiKey, nil)

	// Create a Logger
	logger := nlib.NewSimpleLogger(log.Default())

	// Create State Manager
	stateMgr := nlib.NewSimpleStateManager()

	// Create Context Manager
	contextMgr := nlib.NewSimpleContextManager()

	// Create History Manager
	historyMgr := nlib.NewSimpleHistoryManager()

	// Create Nodes
	firstNode := nlib.NewAINode(lm, logger, stateMgr)
	firstNode.SetID("AI Node")
	outNode := nlib.NewOutputStringNode(os.Stdout, logger, stateMgr)
	outNode.SetID("Output Node")
	firstNode.Connect(outNode)

	// Create the initial Signal with our task
	task := nlib.NewStringData("Why is the sky blue?")
	initialSig := node.NewSignal(firstNode.ID(), contextMgr, historyMgr, task)

	// Send it
	firstNode.InputCh() <- initialSig

	// Wait for the output node to print the result
	stateMgr.WaitFor(outNode.ID())
}
