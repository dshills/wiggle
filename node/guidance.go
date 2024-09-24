package node

import (
	"fmt"

	"github.com/dshills/wiggle/llm"
)

type SimpleGuidance struct {
	llm llm.LLM // Assume llm.LLM is an interface for interacting with a language model
}

func NewSimpleGuidance(llm llm.LLM) *SimpleGuidance {
	return &SimpleGuidance{llm: llm}
}

func (g *SimpleGuidance) Generate(signal Signal) (Signal, error) {
	// Combine input data and context into a message for the LLM
	inputMessage := fmt.Sprintf("Input: %s\nContext: %s", signal.Data.ToMessageList()[0].Content, signal.Context)

	// Generate a response from the LLM (guidance instructions)
	response, err := g.llm.GenerateResponse(inputMessage, "Provide guidance for processing the input.")
	if err != nil {
		return signal, err
	}

	// Update the signal with new guidance
	signal.Response = llm.Message{Content: response}
	signal.Meta = append(signal.Meta, Meta{Key: "guidance", Value: response})

	return signal, nil
}
