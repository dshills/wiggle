package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (g *Gemini) GenEmbed(_ context.Context, str string) ([]float32, error) {
	const embedEP = "/v1beta/models/%%MODEL%%:embedContent?key=%%APIKEY%%"
	ep := fmt.Sprintf("%v%v", g.baseURL, embedEP)
	ep = strings.Replace(ep, "%%MODEL%%", g.model, 1)
	ep = strings.Replace(ep, "%%APIKEY%%", g.apiKey, 1)

	req := embedRequest{
		Model: g.model,
	}
	req.Content.Parts = []part{{Text: str}}

	jsReq, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	// nolint
	resp, err := http.Post(ep, "application/json", bytes.NewReader(jsReq))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	eResp := embedResponse{}
	err = json.NewDecoder(resp.Body).Decode(&eResp)
	if err != nil {
		return nil, err
	}

	return eResp.Embedding.Values, nil
}

type embedRequest struct {
	Model   string `json:"model"`
	Content struct {
		Parts []part `json:"parts"`
	} `json:"content"`
}

type embedResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
}
