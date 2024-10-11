package nlib

import (
	"bytes"
	"text/template"

	"github.com/dshills/wiggle/node"
	"github.com/dshills/wiggle/schema"
)

// Ensure that SimpleGuidance implements the node.Guidance interface
var _ node.Guidance = (*SimpleGuidance)(nil)

// SimpleGuidance provides a structure for generating prompts for large language models (LLMs).
// It stores metadata such as the role, task, audience, goal, steps, format, and tone of the response.
type SimpleGuidance struct {
	Role           string         // The role or persona for the LLM to assume
	Task           string         // The task or instruction the LLM needs to perform
	TargetAudience string         // The intended audience for the output
	Goal           string         // The specific goal or outcome of the task
	Steps          []string       // The steps or process the LLM should follow
	OutputFormat   string         // The desired format for the output
	Tone           string         // The tone the LLM should use in the response
	Schema         *schema.Schema // Optional JSON schema for defining output format
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
	var err error
	tmpl := template.New("basic").Funcs(template.FuncMap{"add": AddFn})
	tmpl, err = tmpl.Parse(BasicTemplate)
	if err != nil {
		return sig, err
	}
	data, err := g.makeTemplateData(sig.Task.String(), context)
	if err != nil {
		return sig, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return sig, err
	}
	sig.Task = &Carrier{TextData: buf.String()}
	return sig, nil
}

func (g *SimpleGuidance) makeTemplateData(task, context string) (BasicTemplateData, error) {
	data := BasicTemplateData{
		Role:           g.Role,
		Task:           g.Task,
		TargetAudience: g.TargetAudience,
		Goal:           g.Goal,
		Steps:          g.Steps,
		OutputFormat:   g.OutputFormat,
		Tone:           g.Tone,
		Context:        context,
		Input:          task,
	}
	if g.Schema != nil {
		sc, err := g.Schema.ToJSON()
		if err != nil {
			return data, err
		}
		data.Schema = sc
	}

	return data, nil
}
