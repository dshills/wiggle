package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dshills/wiggle/llm"
)

func (o *Ollama) GenerateResponse(info string, instruct string) (string, error) {
	msgList := llm.MessageList{llm.UserMsg(fmt.Sprintf("%s %s", info, instruct))}
	resp, err := o.Chat(context.TODO(), msgList)
	return resp.Content, err
}

func (o *Ollama) Chat(ctx context.Context, conv llm.MessageList) (llm.Message, error) {
	oreq := chatRequest{
		Stream:   false,
		Messages: conv,
		Options:  o.options,
		Model:    o.model,
	}
	js, err := json.Marshal(&oreq)
	if err != nil {
		return llm.Message{}, err
	}
	reader := bytes.NewReader(js)
	resp, err := o.send(ctx, o.baseURL, reader)
	if err != nil {
		return llm.Message{}, err
	}

	return resp.Message, nil
}

func (o *Ollama) send(ctx context.Context, baseURL string, reader io.Reader) (*chatResponse, error) {
	const ollamaChatEP = "api/chat"

	ep, err := url.JoinPath(baseURL, ollamaChatEP)
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, err
	}

	chatResp := chatResponse{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, err
	}
	if len(chatResp.Message.Content) == 0 {
		return nil, fmt.Errorf("no content")
	}

	return &chatResp, nil
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []llm.Message `json:"messages"`
	Stream   bool          `json:"stream"`
	Options  Options       `json:"options"`
}

type chatResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done               bool  `json:"done"`
	TotalDuration      int64 `json:"total_duration"`
	LoadDuration       int   `json:"load_duration"`
	PromptEvalCount    int   `json:"prompt_eval_count"`
	PromptEvalDuration int   `json:"prompt_eval_duration"`
	EvalCount          int   `json:"eval_count"`
	EvalDuration       int64 `json:"eval_duration"`
}
