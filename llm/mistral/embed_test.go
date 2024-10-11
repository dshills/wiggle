//go:build mistral
// +build mistral

package mistral_test

import (
	"context"
	"os"
	"testing"

	"github.com/dshills/wiggle/llm/mistral"
)

func TestEmbed(t *testing.T) {
	baseURL := os.Getenv("MISTRAL_API_URL")
	apiKey := os.Getenv("MISTRAL_API_KEY")
	model := "mistral-embed"

	mist := mistral.New(baseURL, model, apiKey, nil)

	ctx := context.TODO()
	embedings, err := mist.GenEmbed(ctx, "Why is the sky blue")
	if err != nil {
		t.Fatal(err)
	}

	if len(embedings) == 0 {
		t.Error("No embeddgings returned")
	}
}
