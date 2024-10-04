package anthropic

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

func (ant *Anthropic) GenerateResponse(info string, instruct string) (string, error) {
	msgList := llm.MessageList{llm.UserMsg(fmt.Sprintf("%s %s", info, instruct))}
	resp, err := ant.Chat(context.TODO(), msgList)
	return resp.Content, err
}

func (ant *Anthropic) Chat(ctx context.Context, msgs llm.MessageList) (llm.Message, error) {
	oreq := chatRequest{
		Stream:    false,
		Messages:  msgs,
		Model:     ant.model,
		MaxTokens: ant.maxTokens,
	}
	js, err := json.Marshal(&oreq)
	if err != nil {
		return llm.Message{}, err
	}
	reader := bytes.NewReader(js)
	resp, err := ant.send(ctx, ant.baseURL, reader)
	if err != nil {
		return llm.Message{}, err
	}
	if resp == nil || len(resp.Content) == 0 {
		return llm.Message{}, fmt.Errorf("nothing returned")
	}
	msg := llm.Message{
		Role:    llm.RoleAssistant,
		Content: resp.Content[0].Text,
	}
	return msg, nil
}

func (ant *Anthropic) send(ctx context.Context, baseURL string, reader io.Reader) (*chatResponse, error) {
	const chatEndpoint = "/v1//messages"

	ep, err := url.JoinPath(baseURL, chatEndpoint)
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
	req.Header.Add("x-api-key", ant.apiKey)
	req.Header.Add("anthropic-version", "2023-06-01")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var byts []byte
		if resp != nil {
			byts, _ = io.ReadAll(resp.Body)
		}
		return nil, fmt.Errorf("ERROR: %v %v %v", resp.StatusCode, resp.Status, string(byts))
	}

	chatResp := chatResponse{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, err
	}

	return &chatResp, nil
}

type chatRequest struct {
	Model         string        `json:"model,omitempty"`          // REQUIRED
	MaxTokens     int           `json:"max_tokens,omitempty"`     // The maximum number of tokens to generate before stopping.
	Messages      []llm.Message `json:"messages,omitempty"`       // REQUIRED
	MetaData      MetaData      `json:"metadata,omitempty"`       // Set a user id
	StopSequences []string      `json:"stop_sequences,omitempty"` // Set of text strings that will trigger a stop
	Stream        bool          `json:"stream,omitempty"`         // Whether to incrementally stream the response using server-sent events.
	System        string        `json:"system"`                   // System prompt
	Temperature   float32       `json:"temperature,omitempty"`    // Amount of randomness injected into the response. 0.0 - 1.0
}

type MetaData struct {
	UserID string `json:"user_id,omitempty"`
}

type chatResponse struct {
	ID      string `json:"id"`
	Content []struct {
		Text  string `json:"text,omitempty"`
		ID    string `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Input struct {
		} `json:"input,omitempty"`
	} `json:"content"`
	Model        string  `json:"model"`
	StopReason   string  `json:"stop_reason"`
	StopSequence *string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}
