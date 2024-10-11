package mistral

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dshills/wiggle/llm"
)

func (m *Mistral) GenerateResponse(info string, instruct string) (string, error) {
	msgList := llm.MessageList{llm.UserMsg(fmt.Sprintf("%s %s", info, instruct))}
	resp, err := m.Chat(context.TODO(), msgList)
	return resp.Content, err
}

func (m *Mistral) Chat(ctx context.Context, conv llm.MessageList) (llm.Message, error) {
	chatReq := chatRequest{
		Model:    m.model,
		Messages: conv,
	}
	jsReq, err := json.Marshal(&chatReq)
	if err != nil {
		return llm.Message{}, err
	}

	chatResp, err := m.send(ctx, bytes.NewReader(jsReq))
	if err != nil {
		return llm.Message{}, err
	}

	return chatResp.Choices[0].Message, nil
}

func (m *Mistral) send(ctx context.Context, reader io.Reader) (*chatResponse, error) {
	const chatEP = "/v1/chat/completions"
	ep, err := url.JoinPath(m.baseURL, chatEP)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, reader)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))

	httpResp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= 300 {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("Mistral: Chat: %v %v\n%v", httpResp.StatusCode, httpResp.Status, string(body))
	}

	resp := chatResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no content")
	}
	return &resp, err
}

type chatRequest struct {
	Model       string        `json:"model,omitempty"`
	Messages    []llm.Message `json:"messages,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        int           `json:"top_p,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	SafePrompt  bool          `json:"safe_prompt,omitempty"`
	RandomSeed  int           `json:"random_seed,omitempty"`
}

type chatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      llm.Message `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
