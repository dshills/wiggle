package ollama

import (
	"net/http"

	"github.com/dshills/wiggle/llm"
)

type Ollama struct {
	model      string
	options    Options
	baseURL    string
	httpClient *http.Client
}

func New(baseURL, model string, options *Options) llm.LLM {
	o := Ollama{
		baseURL:    baseURL,
		model:      model,
		httpClient: http.DefaultClient,
	}
	if options != nil {
		o.options = *options
	}
	return &o
}

func (o *Ollama) SetModel(model string) {
	o.model = model
}

func (o *Ollama) Model() string {
	return o.model
}
