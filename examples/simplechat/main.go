package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dshills/wiggle/llm/openai"
	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/node"
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
	stateMgr.SetLogger(logger)

	// Define output writer
	writer := os.Stdout

	// Create Nodes
	options := node.Options{ID: "Input-Node"}
	inputNode := nlib.NewInteractiveNode(stateMgr, options)
	options = node.Options{ID: "AI-Node"}
	firstNode := nlib.NewAINode(lm, stateMgr, options)
	options = node.Options{ID: "Output-Node"}
	outNode := nlib.NewOutputStringNode(writer, stateMgr, options)

	// Connections
	inputNode.Connect(firstNode)
	firstNode.Connect(outNode)
	outNode.Connect(inputNode)

	signal := node.Signal{NodeID: inputNode.ID()}
	// Send it
	inputNode.InputCh() <- signal

	// Wait forever
	stateMgr.WaitFor(nil)

	// Print the history
	for _, hx := range stateMgr.GetHistory() {
		fmt.Println(nlib.SignalToLog(hx))
	}
}
