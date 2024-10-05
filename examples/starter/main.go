package main

import (
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

	// Define output writer
	writer := os.Stdout

	// Create Nodes
	firstNode := nlib.NewAINode(lm, stateMgr, node.Options{ID: "AI-Node"})
	outNode := nlib.NewOutputStringNode(writer, stateMgr, node.Options{ID: "Output Node"})

	// Connect
	firstNode.Connect(outNode)

	// Send it
	firstNode.InputCh() <- nlib.NewDefaultSignal(firstNode, "Why is the sky blue?")

	// Wait for the output node to print the result
	stateMgr.WaitFor(outNode)
}
