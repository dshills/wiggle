//go:build mistral
// +build mistral

package mistral_test

import (
	"context"
	"os"
	"testing"

	"github.com/dshills/wiggle/llm"
	"github.com/dshills/wiggle/llm/mistral"
)

func TestChat(t *testing.T) {
	baseURL := os.Getenv("MISTRAL_API_URL")
	apiKey := os.Getenv("MISTRAL_API_KEY")
	model := "mistral-small-latest"

	mist := mistral.New(baseURL, model, apiKey, nil)

	ctx := context.TODO()
	msgs := llm.MessageList{
		llm.Message{Role: llm.RoleUser, Content: "Why is the sky blue?"},
	}
	respMsg, err := mist.Chat(ctx, msgs)
	if err != nil {
		t.Fatal(err)
	}

	if respMsg.Content == "" {
		t.Errorf("Expected a response got none")
	}
}
