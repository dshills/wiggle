package anthropic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dshills/wiggle/llm"
)

const (
	ModelSonnet35 = "claude-3-5-sonnet-20240620"
	ModelSonnet3  = "claude-3-sonnet-20240229"
	ModelOpus3    = "claude-3-opus-20240229"
	ModelHaiku3   = "claude-3-haiku-20240307"
)

// Compile-time check
var _ llm.LLM = (*Anthropic)(nil)

type Anthropic struct {
	model      string
	baseURL    string
	httpClient *http.Client
	apiKey     string
	maxTokens  int
}

func New(baseURL, model, apiKey string, maxTokens int) *Anthropic {
	ant := Anthropic{
		baseURL:    baseURL,
		model:      model,
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
		maxTokens:  maxTokens,
	}
	if maxTokens < 10 {
		ant.maxTokens = 1024
	}
	return &ant
}

func (ant *Anthropic) SetModel(model string) {
	ant.model = model
}

func (ant *Anthropic) Model() string {
	return ant.model
}

func (ant *Anthropic) GenEmbed(_ context.Context, _ string) ([]float32, error) {
	// Requires Voyage HTTP API
	return nil, fmt.Errorf("not implemented")
}

func (ant *Anthropic) AvailableModels() ([]llm.Model, error) {
	// Strangly Anthropic does not appear to have an API to get a list of models
	// I'm hardcoding this as of Oct 4, 2024
	models := []llm.Model{}
	mod := llm.Model{Name: ModelSonnet35}
	models = append(models, mod)
	mod = llm.Model{Name: ModelSonnet3}
	models = append(models, mod)
	mod = llm.Model{Name: ModelOpus3}
	models = append(models, mod)
	mod = llm.Model{Name: ModelHaiku3}
	models = append(models, mod)
	return models, nil
}
