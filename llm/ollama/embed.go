package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (o *Ollama) GenEmbed(ctx context.Context, txt string) ([]float32, error) {
	const embedEP = "/api/embed"
	ep, err := url.JoinPath(o.baseURL, embedEP)
	if err != nil {
		return nil, err
	}
	req := embedReq{Model: o.model, Input: txt}
	js, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, ep, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
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
	if len(resp.Embeddings) == 0 {
		return nil, nil
	}

	return resp.Embeddings[0], nil
}

type embedReq struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embedResp struct {
	Model           string      `json:"model"`
	Embeddings      [][]float32 `json:"embeddings"`
	TotalDuration   int         `json:"total_duration"`
	LoadDuration    int         `json:"load_duration"`
	PromptEvalCount int         `json:"prompt_eval_count"`
}
