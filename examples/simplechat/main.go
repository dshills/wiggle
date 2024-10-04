package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dshills/wiggle/llm/openai"
	"github.com/dshills/wiggle/nlib"
)

const envURL = "OPENAI_API_URL"
const envKey = "OPENAI_API_KEY"
const model = "gpt-4o"

func main() {
	// Setup LLM
	lm := openai.New(os.Getenv(envURL), model, os.Getenv(envKey), nil)

	// Create a Logger
	logger := nlib.NewSimpleLogger(log.Default())

	// Create State Manager
	stateMgr := nlib.NewSimpleStateManager(logger)

	// Define output writer
	writer := os.Stdout

	// Create Nodes
	inputNode := nlib.NewInteractiveNode(logger, stateMgr, "Input Node")
	firstNode := nlib.NewAINode(lm, logger, stateMgr, "AI Node")
	inputNode.Connect(firstNode)
	outNode := nlib.NewOutputStringNode(writer, logger, stateMgr, "Output Node")
	firstNode.Connect(outNode)
	outNode.Connect(inputNode)

	signal := nlib.NewDefaultSignal(inputNode, "")
	// Send it
	inputNode.InputCh() <- signal

	// Wait forever
	stateMgr.WaitFor(nil)

	// Print the history
	for _, hx := range signal.History.GetHistory() {
		fmt.Println(nlib.SignalToLog(hx))
	}
}
