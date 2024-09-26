package main

import (
	"log"
	"os"

	"github.com/dshills/wiggle/llm/openai"
	"github.com/dshills/wiggle/nlib"
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

	// Define output writer
	writer := os.Stdout

	// Create Nodes
	firstNode := nlib.NewAINode(lm, logger, stateMgr)
	firstNode.SetID("AI Node")
	outNode := nlib.NewOutputStringNode(writer, logger, stateMgr)
	outNode.SetID("Output Node")
	firstNode.Connect(outNode)

	// Send it
	firstNode.InputCh() <- nlib.NewDefaultSignal(firstNode, "Why is the sky blue?")

	// Wait for the output node to print the result
	stateMgr.WaitFor(outNode)
}
