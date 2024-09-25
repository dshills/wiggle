package openai

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

func (ai *OpenAI) GenerateResponse(info string, instruct string) (string, error) {
	msgList := llm.MessageList{llm.UserMsg(fmt.Sprintf("%s %s", info, instruct))}
	resp, err := ai.Chat(context.TODO(), msgList)
	return resp.Content, err
}

func (ai *OpenAI) Chat(ctx context.Context, msgs llm.MessageList) (llm.Message, error) {
	js, err := ai.encodeRequest(msgs)
	if err != nil {
		return llm.Message{}, err
	}
	reader := bytes.NewReader(js)
	resp, err := ai.send(ctx, ai.baseURL, reader)
	if err != nil {
		return llm.Message{}, err
	}
	if resp == nil || len(resp.Choices) == 0 {
		return llm.Message{}, fmt.Errorf("OpenAI: Chat: No data returned")
	}

	return resp.Choices[0].Message, nil
}

func (ai *OpenAI) encodeRequest(msgs llm.MessageList) ([]byte, error) {
	var js []byte
	var err error
	switch {
	case ai.options != nil && len(ai.options.Tools) > 0:
		req := ai.options.asRequest()
		req.Stream = false
		req.Messages = msgs
		req.Model = ai.model
		js, err = json.Marshal(&req)
		if err != nil {
			return nil, err
		}

	case ai.options != nil:
		req := chatRequest{
			Stream:      false,
			Messages:    msgs,
			Model:       ai.model,
			Temperature: ai.options.Temperature,
			MaxTokens:   ai.options.MaxTokens,
		}
		js, err = json.Marshal(&req)
		if err != nil {
			return nil, err
		}

	default:
		req := chatRequest{
			Stream:   false,
			Messages: msgs,
			Model:    ai.model,
		}
		js, err = json.Marshal(&req)
		if err != nil {
			return nil, err
		}
	}
	return js, nil
}

func (ai *OpenAI) send(ctx context.Context, baseURL string, reader io.Reader) (*chatResponse, error) {
	const chatEP = "/v1/chat/completions"

	ep, err := url.JoinPath(baseURL, chatEP)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, ep, reader)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ai.apiKey))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI: Chat: %v %v\n%v", resp.StatusCode, resp.Status, string(body))
	}

	chatResp := chatResponse{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, err
	}
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no data returned")
	}

	return &chatResp, nil
}

type chatResponse struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index        int         `json:"index"`
		Message      llm.Message `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
