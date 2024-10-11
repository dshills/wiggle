//go:build gemini
// +build gemini

package gemini_test

import (
	"context"
	"os"
	"testing"

	"github.com/dshills/wiggle/llm/gemini"
)

func TestEmbed(t *testing.T) {
	baseURL := os.Getenv("GEMINI_API_URL")
	apiKey := os.Getenv("GEMINI_API_KEY")
	model := "text-embedding-004"
	gem := gemini.New(baseURL, model, apiKey, nil)

	ctx := context.TODO()
	_, err := gem.GenEmbed(ctx, "Why is the sky blue")
	if err != nil {
		t.Fatal(err)
	}
}
