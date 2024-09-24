package openai

import "github.com/dshills/wiggle/llm"

type Options struct {
	Logprobs          bool   `json:"logprobs,omitempty"`
	MaxTokens         *int   `json:"max_tokens,omitempty"`
	ParallelToolCalls bool   `json:"parallel_tool_calls,omitempty"`
	Temperature       *int   `json:"temperature,omitempty"`
	ToolChoice        string `json:"tool_choice,omitempty"`
	Tools             []Tool `json:"tools,omitempty"`
	TopLogprobs       int    `json:"top_logprobs,omitempty"`
}

type Tool struct {
	Type     string   `json:"type,omitempty"`
	Function Function `json:"function,omitempty"`
}

type Function struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Parameters  Parameter `json:"parameters,omitempty"`
}

type Parameter struct {
	Type       string     `json:"type,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Required   []string   `json:"required,omitempty"`
}

type Properties struct {
	Location Location `json:"location,omitempty"`
	Unit     Unit     `json:"unit,omitempty"`
}

type Location struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

type Unit struct {
	Type string   `json:"type,omitempty"`
	Enum []string `json:"enum,omitempty"`
}

func (o Options) asRequest() chatRequestWithTools {
	return chatRequestWithTools{
		Options: o,
	}
}

type chatRequest struct {
	Model       string        `json:"model,omitempty"`
	Messages    []llm.Message `json:"messages,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Temperature *int          `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
}

type chatRequestWithTools struct {
	Options
	Model    string        `json:"model,omitempty"`
	Messages []llm.Message `json:"messages,omitempty"`
	Stream   bool          `json:"stream,omitempty"`
}
