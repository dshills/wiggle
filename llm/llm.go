package llm

import (
	"context"
)

type Model struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	Format       string `json:"format"`
	Family       string `json:"family"`
	Parameters   string `json:"parameters"`
	Quantization string `json:"quantization"`
}

type LLM interface {
	GenerateResponse(string, string) (string, error)
	Chat(ctx context.Context, msgs MessageList) (Message, error)
	GenEmbed(ctx context.Context, txt string) ([]float32, error)
	AvailableModels() ([]Model, error)
	SetModel(model string)
	Model() string
}
