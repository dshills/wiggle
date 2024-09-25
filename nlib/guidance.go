package nlib

import (
	"github.com/dshills/wiggle/llm"
	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.Guidance = (*SimpleGuidance)(nil)

// SimpleGuidance interacts with a large language model (LLM) to generate
// instructions for processing signals. It uses the input data, context, and metadata
// to construct a message for the LLM and updates the signal with the response
// and new guidance, which can be used by subsequent nodes in the workflow.
type SimpleGuidance struct {
	llm llm.LLM // Assume llm.LLM is an interface for interacting with a language model
}

func NewSimpleGuidance(llm llm.LLM) *SimpleGuidance {
	return &SimpleGuidance{llm: llm}
}

func (g *SimpleGuidance) Generate(signal node.Signal) (node.Signal, error) {
	return signal, nil
}
