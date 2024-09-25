package nlib

import (
	"fmt"
	"strings"

	"github.com/dshills/wiggle/node"
)

// Compile-time check
var _ node.Guidance = (*SimpleGuidance)(nil)

type SimpleGuidance struct {
	Role           string `json:"role"`
	Task           string `json:"task"`
	TargetAudience string `json:"target_audience"`
	Goal           string `json:"goal"`
	Steps          string `json:"steps"`
	OutputFormat   string `json:"output_format"`
	Tone           string `json:"tone"`
}

func NewSimpleGuidance() *SimpleGuidance {
	return &SimpleGuidance{}
}

func (g *SimpleGuidance) Generate(sig node.Signal) (node.Signal, error) {
	con, err := sig.Context.GetContext(sig.NodeID)
	if err != nil {
		prompt := g.prompt(sig.Data.String(), "")
		sig.Data = NewStringData(prompt)
		return sig, nil
	}
	prompt := g.prompt(sig.Data.String(), con.String())
	sig.Data = NewStringData(prompt)
	return sig, nil
}

type TextGuidance struct {
}

func (g *SimpleGuidance) prompt(input, context string) string {
	builder := strings.Builder{}
	if g.Role != "" {
		builder.WriteString(fmt.Sprintf("<role>%s</role>\n", g.Role))
	}
	if g.Task != "" {
		builder.WriteString(fmt.Sprintf("<task>%s</task>\n", g.Task))
	}
	if g.TargetAudience != "" {
		builder.WriteString(fmt.Sprintf("<target audience>%s</target audience>\n", g.TargetAudience))
	}
	if g.Goal != "" {
		builder.WriteString(fmt.Sprintf("<goal>%s</goal>\n", g.Goal))
	}
	if g.Steps != "" {
		builder.WriteString(fmt.Sprintf("<steps>%s</steps>\n", g.Steps))
	}
	if g.OutputFormat != "" {
		builder.WriteString(fmt.Sprintf("<format>%s</format>\n", g.OutputFormat))
	}
	if g.Tone != "" {
		builder.WriteString(fmt.Sprintf("<tone>%s</tone>\n", g.Tone))
	}
	if context != "" {
		builder.WriteString(fmt.Sprintf("<context>%s</context>", context))
	}
	if input != "" {
		builder.WriteString(fmt.Sprintf("<input>%s</input>", input))
	}

	return builder.String()
}
