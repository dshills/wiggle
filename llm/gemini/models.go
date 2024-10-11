package gemini

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/dshills/wiggle/llm"
)

func (g *Gemini) AvailableModels() ([]llm.Model, error) {
	const modelEP = "/v1beta/models?key=%%APIKEY%%"
	modEP := strings.ReplaceAll(modelEP, "%%APIKEY%%", g.apiKey)
	ep, err := url.JoinPath(g.baseURL, modEP)
	if err != nil {
		return nil, err
	}
	// nolint
	httpResp, err := http.Get(ep)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	mods := []model{}
	if err := json.NewDecoder(httpResp.Body).Decode(&mods); err != nil {
		return nil, err
	}

	modelList := []llm.Model{}
	for _, mm := range mods {
		modelList = append(modelList, mm.asModel())
	}

	return modelList, nil
}

type model struct {
	Name                       string   `json:"name"`
	BaseModelID                string   `json:"baseModelId"`
	Version                    string   `json:"version"`
	DisplayName                string   `json:"displayName"`
	Description                string   `json:"description"`
	InputTokenLimit            int      `json:"inputTokenLimit"`
	OutputTokenLimit           int      `json:"outputTokenLimit"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
	Temperature                float32  `json:"temperature"`
	MaxTemperature             float32  `json:"maxTemperature"`
	TopP                       float32  `json:"topP"`
	TopK                       int      `json:"topK"`
}

func (m *model) asModel() llm.Model {
	return llm.Model{
		Name:   m.Name,
		Family: m.BaseModelID,
	}
}
