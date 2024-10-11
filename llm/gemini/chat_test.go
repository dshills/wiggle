//go:build gemini
// +build gemini

package gemini_test

import (
	"context"
	"os"
	"testing"

	"github.com/dshills/wiggle/llm"
	"github.com/dshills/wiggle/llm/gemini"
)

func TestChat(t *testing.T) {
	baseURL := os.Getenv("GEMINI_API_URL")
	apiKey := os.Getenv("GEMINI_API_KEY")
	model := "gemini-1.5-flash"
	gem := gemini.New(baseURL, model, apiKey, nil)

	ctx := context.TODO()
	msgs := llm.MessageList{
		llm.Message{Role: llm.RoleUser, Content: "Why is the sky blue?"},
	}
	respMsg, err := gem.Chat(ctx, msgs)
	if err != nil {
		t.Fatal(err)
	}

	if respMsg.Content == "" {
		t.Errorf("Expected a response got none")
	}
}
