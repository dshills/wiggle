package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (ai *OpenAI) GenEmbed(ctx context.Context, txt string) ([]float32, error) {
	const embedEP = "/v1/embeddings"
	ep, err := url.JoinPath(ai.baseURL, embedEP)
	if err != nil {
		return nil, err
	}
	req := embedReq{Model: ai.model, Input: txt, EncodingFormat: "float"}
	js, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(js))
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
		return nil, fmt.Errorf("%v %v", httpResp.StatusCode, httpResp.Status)
	}

	resp := embedResp{}
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no data returned")
	}
	if len(resp.Data) > 1 {
		return nil, fmt.Errorf("expected single vectore got %v", len(resp.Data))
	}

	return resp.Data[0].Embedding, nil
}

type embedReq struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	EncodingFormat string `json:"encoding_format"`
}
type embedResp struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}
