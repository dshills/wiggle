package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dshills/wiggle/llm/openai"
	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/node"
	"github.com/dshills/wiggle/schema"
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

	sig := node.Signal{
		NodeID: "task-node",
		Task:   &nlib.Carrier{TextData: "Write a complete rope algorithm package. Include methods for using a character index as well as line and column values"},
	}

	outSchema, err := schema.FromSchemaJSON([]byte(schemaStr))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	guide := makeGuidance(outSchema)
	taskNode := nlib.NewAINode(lm, stateMgr, node.Options{ID: "task-node", Guidance: guide})
	validateNode := nlib.NewJSONValidatorNode(stateMgr, outSchema, node.Options{ID: "validator-node"})
	outNode := nlib.NewOutputStringNode(writer, stateMgr, node.Options{ID: "Output Node"})

	// Connect
	taskNode.Connect(validateNode)
	validateNode.Connect(outNode)

	taskNode.InputCh() <- sig

	stateMgr.WaitFor(outNode)
}
