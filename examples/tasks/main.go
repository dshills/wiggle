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

	guide := nlib.SimpleGuidance{
		Role:           "You are an expert at dividing Google Go projects into smaller tasks. You are a golang software expert. You never write code",
		Task:           "Take the input and divide it into 10 smaller coding tasks. Each task should be something a developer can begin coding immediately",
		TargetAudience: "Principal level software engineers",
		Goal:           "10 tasks that can be coded",
		Steps:          []string{"Review the entire input", "Consider possible coding tasks", "Format into the 10 tasks"},
		OutputFormat:   "Numbered list in markdown",
		Tone:           "professional software engineer",
	}

	taskNode := nlib.NewAINode(lm, stateMgr, node.Options{ID: "Task Node", Guidance: &guide})

	guide = nlib.SimpleGuidance{
		Role:           "You are an expert JSON constructor. You never write code",
		Task:           "Take a list of ten prograaming tasks formatted as markdown and convert it to JSON. Do not include anything other than the JSON in the output",
		TargetAudience: "Principal level software engineers",
		Goal:           "Properly formed JSON containing the 10 programming tasks",
		Steps:          []string{"Review the entire input", "Identify the ten tasks", "Format into JSON based on the OutputFormat"},
		OutputFormat: `
		[
			{"task_name": "name of task", "task_steps": ["step1", "step2"], "notes": ["note", "note"]}
		]
		`,
		Tone: "professional software engineer",
	}

	jsonNode := nlib.NewAINode(lm, stateMgr, node.Options{ID: "JSON Node", Guidance: &guide})
	outNode := nlib.NewOutputStringNode(writer, stateMgr, node.Options{ID: "Output Node"})

	// Connect
	taskNode.Connect(jsonNode)
	jsonNode.Connect(outNode)

	sig := node.Signal{
		NodeID: taskNode.ID(),
		Task:   &nlib.Carrier{TextData: "Write a complete rope algorithm package. Include methods for using a character index as well as line and column values"},
	}

	taskNode.InputCh() <- sig

	stateMgr.WaitFor(outNode)
}
