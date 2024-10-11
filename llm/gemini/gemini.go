package gemini

import (
	"net/http"

	"github.com/dshills/wiggle/llm"
)

// Compile-time check
var _ llm.LLM = (*Gemini)(nil)

type Gemini struct {
	model      string
	options    Options
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, model, apiKey string, options *Options) *Gemini {
	g := Gemini{
		baseURL:    baseURL,
		model:      model,
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
	if options != nil {
		g.options = *options
	}
	return &g
}

func (g *Gemini) SetModel(model string) {
	g.model = model
}

func (g *Gemini) Model() string {
	return g.model
}
