package nlib

import (
	"fmt"
	"strings"

	"github.com/dshills/wiggle/node"
)

// Ensure that SimpleGuidance implements the node.Guidance interface
var _ node.Guidance = (*SimpleGuidance)(nil)

// SimpleGuidance provides a structure for generating prompts for large language models (LLMs).
// It stores metadata such as the role, task, audience, goal, steps, format, and tone of the response.
type SimpleGuidance struct {
	Role           string `json:"role"`            // The role or persona for the LLM to assume
	Task           string `json:"task"`            // The task or instruction the LLM needs to perform
	TargetAudience string `json:"target_audience"` // The intended audience for the output
	Goal           string `json:"goal"`            // The specific goal or outcome of the task
	Steps          string `json:"steps"`           // The steps or process the LLM should follow
	OutputFormat   string `json:"output_format"`   // The desired format for the output
	Tone           string `json:"tone"`            // The tone the LLM should use in the response
}

// NewSimpleGuidance creates a new instance of SimpleGuidance with default values.
// This is a constructor function that returns a pointer to the newly created SimpleGuidance.
func NewSimpleGuidance() *SimpleGuidance {
	return &SimpleGuidance{}
}

// Generate constructs a new signal by creating a prompt from the guidance metadata.
// It retrieves additional context from the signal, if available, and includes that in the prompt.
// If no context is found, the prompt is generated without it. The generated prompt is assigned
// to the signal's Task and returned for further processing.
func (g *SimpleGuidance) Generate(sig node.Signal, context string) (node.Signal, error) {
	// Generate a prompt using both the signal's task and the retrieved context
	prompt := g.prompt(sig.Task.String(), context)
	sig.Task = &Carrier{TextData: prompt}
	return sig, nil
}

// prompt generates the prompt text based on the guidance metadata and any provided context.
// It constructs the prompt using various tags like role, task, audience, goal, etc.
// The context is included as part of the prompt if it's available.
func (g *SimpleGuidance) prompt(input, context string) string {
	builder := strings.Builder{}

	// Add role to the prompt if it's provided
	if g.Role != "" {
		builder.WriteString(fmt.Sprintf("<role>%s</role>\n", g.Role))
	}

	// Add task to the prompt if it's provided
	if g.Task != "" {
		builder.WriteString(fmt.Sprintf("<task>%s</task>\n", g.Task))
	}

	// Add target audience to the prompt if it's provided
	if g.TargetAudience != "" {
		builder.WriteString(fmt.Sprintf("<target audience>%s</target audience>\n", g.TargetAudience))
	}

	// Add goal to the prompt if it's provided
	if g.Goal != "" {
		builder.WriteString(fmt.Sprintf("<goal>%s</goal>\n", g.Goal))
	}

	// Add steps to the prompt if they're provided
	if g.Steps != "" {
		builder.WriteString(fmt.Sprintf("<steps>%s</steps>\n", g.Steps))
	}

	// Add output format to the prompt if it's provided
	if g.OutputFormat != "" {
		builder.WriteString(fmt.Sprintf("<format>%s</format>\n", g.OutputFormat))
	}

	// Add tone to the prompt if it's provided
	if g.Tone != "" {
		builder.WriteString(fmt.Sprintf("<tone>%s</tone>\n", g.Tone))
	}

	// Add context to the prompt if it's available
	if context != "" {
		builder.WriteString(fmt.Sprintf("<context>%s</context>", context))
	}

	// Add the input (task) to the prompt
	if input != "" {
		builder.WriteString(fmt.Sprintf("<input>%s</input>", input))
	}

	// Return the complete prompt as a string
	return builder.String()
}
