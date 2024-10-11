package mistral

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (m *Mistral) GenEmbed(ctx context.Context, str string) ([]float32, error) {
	const embEP = "/v1/embeddings"
	ep, err := url.JoinPath(m.baseURL, embEP)
	if err != nil {
		return nil, err
	}

	req := embedRequest{
		Input:          str,
		Model:          m.model,
		EncodingFormat: "float",
	}
	jsReq, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(jsReq))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))

	httpResp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Mistral: client.Do: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 300 {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("Mistral: Embeddings: %v %v\n%v", httpResp.StatusCode, httpResp.Status, string(body))
	}

	resp := embedResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("Mistral: JSON: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no content")
	}

	return resp.Data[0].Embedding, nil
}

type embedRequest struct {
	Input          string `json:"input,omitempty"`
	Model          string `json:"model,omitempty"`
	EncodingFormat string `json:"encoding_format,omitempty"`
}

type embedResponse struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}
