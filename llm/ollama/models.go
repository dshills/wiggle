package ollama

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/dshills/wiggle/llm"
)

func (o *Ollama) AvailableModels() ([]llm.Model, error) {
	const modelEP = "/api/tags"
	ep, err := url.JoinPath(o.baseURL, modelEP)
	if err != nil {
		return nil, err
	}
	// nolint
	httpResp, err := http.Get(ep)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	mods := models{}
	if err := json.NewDecoder(httpResp.Body).Decode(&mods); err != nil {
		return nil, err
	}
	return mods.AsModels(), nil
}

type models struct {
	Models []struct {
		Name       string `json:"name"`
		ModifiedAt string `json:"modified_at"`
		Size       int64  `json:"size"`
		Digest     string `json:"digest"`
		Details    struct {
			Format            string      `json:"format"`
			Family            string      `json:"family"`
			Families          interface{} `json:"families"`
			ParameterSize     string      `json:"parameter_size"`
			QuantizationLevel string      `json:"quantization_level"`
		} `json:"details"`
	} `json:"models"`
}

func (m models) AsModels() []llm.Model {
	mods := []llm.Model{}
	for _, m := range m.Models {
		llmMod := llm.Model{
			Name:         m.Name,
			Size:         m.Size,
			Format:       m.Details.Format,
			Family:       m.Details.Family,
			Parameters:   m.Details.ParameterSize,
			Quantization: m.Details.QuantizationLevel,
		}
		mods = append(mods, llmMod)
	}
	return mods
}
