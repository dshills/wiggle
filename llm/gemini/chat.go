package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dshills/wiggle/llm"
)

func (g *Gemini) GenerateResponse(info string, instruct string) (string, error) {
	msgList := llm.MessageList{llm.UserMsg(fmt.Sprintf("%s %s", info, instruct))}
	resp, err := g.Chat(context.TODO(), msgList)
	return resp.Content, err
}

func (g *Gemini) Chat(ctx context.Context, conv llm.MessageList) (llm.Message, error) {
	conlist := []content{}
	for _, m := range conv {
		con := content{Role: m.Role, Parts: []part{{Text: m.Content}}}
		conlist = append(conlist, con)
	}
	req := chatRequest{Contents: conlist}
	js, err := json.Marshal(&req)
	if err != nil {
		return llm.Message{}, err
	}
	reader := bytes.NewReader(js)
	resp, err := g.send(ctx, g.baseURL, reader)
	if err != nil {
		return llm.Message{}, err
	}

	retRespone := llm.Message{
		Role:    resp.Candidates[0].Content.Role,
		Content: resp.Candidates[0].Content.Parts[0].Text,
	}
	return retRespone, nil
}

func (g *Gemini) send(ctx context.Context, baseURL string, reader io.Reader) (*chatResponse, error) {
	const geminiEP = "/v1beta/models/%%MODEL%%:generateContent?key=%%APIKEY%%"
	ep := fmt.Sprintf("%v%v", baseURL, geminiEP)
	ep = strings.Replace(ep, "%%MODEL%%", g.model, 1)
	ep = strings.Replace(ep, "%%APIKEY%%", g.apiKey, 1)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, ep, reader)
	if err != nil {
		return nil, fmt.Errorf("completion: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%v %v", resp.StatusCode, resp.Status)
	}

	chatResp := chatResponse{}
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return nil, err
	}
	if len(chatResp.Candidates) == 0 {
		return nil, fmt.Errorf("no content")
	}

	return &chatResp, nil
}

type chatRequest struct {
	Contents []content `json:"contents"`
}

type chatResponse struct {
	Candidates []candidate
}

type candidate struct {
	Content      content `json:"content"`
	FinishReason string  `json:"finish_reason"`
	TokenCount   int     `json:"token_count"`
	Index        int     `json:"index"`
}

type content struct {
	Role  string `json:"role"`
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}
