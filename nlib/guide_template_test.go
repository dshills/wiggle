package nlib

import (
	"bytes"
	"testing"
	"text/template"
)

func TestBasicTemplate(t *testing.T) {
	var err error
	tmpl := template.New("basic").Funcs(template.FuncMap{"add": AddFn})
	tmpl, err = tmpl.Parse(BasicTemplate)
	if err != nil {
		t.Fatal(err)
	}

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
