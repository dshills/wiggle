package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dshills/wiggle/llm"
)

// Compile-time check
var _ llm.LLM = (*OpenAI)(nil)

type OpenAI struct {
	baseURL string
	model   string
	apiKey  string
	options *Options
}

func New(baseURL, model, apiKey string, options *Options) llm.LLM {
	return &OpenAI{model: model, baseURL: baseURL, apiKey: apiKey, options: options}
}

func (ai *OpenAI) AvailableModels() ([]llm.Model, error) {
	const modelEP = "/v1/models"
	ep, err := url.JoinPath(ai.baseURL, modelEP)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ai.apiKey))

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 300 {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("OpenAI: Chat: %v %v\n%v", httpResp.StatusCode, httpResp.Status, string(body))
	}

	mods := models{}
	if err := json.NewDecoder(httpResp.Body).Decode(&mods); err != nil {
		return nil, err
	}
	return mods.AsModels(), nil
}

func (ai *OpenAI) SetModel(model string) {
	ai.model = model
}

func (ai *OpenAI) Model() string {
	return ai.model
}

type models struct {
	Object string `json:"object"`
	Data   []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

func (m models) AsModels() []llm.Model {
	mods := []llm.Model{}
	for _, m := range m.Data {
		llmMod := llm.Model{
			Name: m.ID,
		}
		mods = append(mods, llmMod)
	}
	return mods
}
