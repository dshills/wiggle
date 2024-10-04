package anthropic

import (
	"context"
	"os"
	"testing"

	"github.com/dshills/wiggle/llm"
)

func TestChat(t *testing.T) {
	baseURL := os.Getenv("ANTHROPIC_API_URL")
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	maxTokens := 1024
	ant := New(baseURL, ModelSonnet35, apiKey, maxTokens)

	ctx := context.TODO()
	msgs := llm.MessageList{
		llm.Message{Role: llm.RoleUser, Content: "Why is the sky blue?"},
	}
	respMsg, err := ant.Chat(ctx, msgs)
	if err != nil {
		t.Fatal(err)
	}

	if respMsg.Content == "" {
		t.Errorf("Expected a response got none")
	}
}
