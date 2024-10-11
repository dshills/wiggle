package main

import (
	"github.com/dshills/wiggle/nlib"
	"github.com/dshills/wiggle/node"
	"github.com/dshills/wiggle/schema"
)

func makeGuidance(sc schema.Schema) node.Guidance {
	guide := nlib.SimpleGuidance{
		Role:           "You are an expert at dividing Google Go projects into smaller tasks. You are a golang software expert. You never write code",
		Task:           "Take the input and divide it into 10 smaller coding tasks. Each task should be something a developer can begin coding immediately",
		TargetAudience: "Principal level software engineers",
		Goal:           "10 tasks that can be coded",
		Steps:          []string{"Review the entire input", "Consider possible coding tasks", "Format into the 10 tasks"},
		OutputFormat:   "Properly formed JSON containing the 10 programming tasks. Do not return any data other than the JSON",
		Tone:           "professional software engineer",
		Schema:         &sc,
	}

	return &guide
}
