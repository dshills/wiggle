package mistral

import (
	"net/http"

	"github.com/dshills/wiggle/llm"
)

// Compile-time check
var _ llm.LLM = (*Mistral)(nil)

type Mistral struct {
	model      string
	options    Options
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, model, apiKey string, options *Options) *Mistral {
	m := Mistral{
		baseURL:    baseURL,
		model:      model,
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
	if options != nil {
		m.options = *options
	}
	return &m
}

func (m *Mistral) SetModel(model string) {
	m.model = model
}

func (m *Mistral) Model() string {
	return m.model
}
