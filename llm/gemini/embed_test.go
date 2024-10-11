//go:build gemini
// +build gemini

package gemini_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/dshills/wiggle/llm/gemini"
)

func TestChat(t *testing.T) {
	baseURL := os.Getenv("GEMINI_API_URL")
	apiKey := os.Getenv("GEMINI_API_KEY")
	model := "text-embedding-004"
	gem := gemini.New(baseURL, model, apiKey, nil)

	ctx := context.TODO()
	embedings, err := gem.GenEmbed(ctx, "Why is the sky blue")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(embedings)
}
