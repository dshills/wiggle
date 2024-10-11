package nlib

import (
	"bytes"
	"testing"
)

func TestBasicTemplate(t *testing.T) {
	tmpl := basicTmpl

	data := BasicTemplateData{
		Role:           "",
		Task:           "",
		TargetAudience: "",
		Goal:           "",
		Steps:          []string{"banana", "grape", "apple"},
		OutputFormat:   "",
		Tone:           "",
		Context:        "",
		Input:          "",
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Error(err)
	}
}
