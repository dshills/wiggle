package nlib

import "text/template"

// BasicTemplateData is the data to populate in the BasicTemplate
type BasicTemplateData struct {
	Role           string
	Task           string
	TargetAudience string
	Goal           string
	Steps          []string
	OutputFormat   string
	Tone           string
	Context        string
	Input          string
}

// AddFn is used by the BasicTemplate to make a numbered list
var AddFn = func(a, b int) int { return a + b }

// BasicTemaplte is a template for building Guidance
const BasicTemplate = `
{{if .Role}}
<role>
	{{ .Role }}
</role>
{{end}}

{{if .Task}}
<task>
	{{ .Task }}
</task>
{{end}}

{{if .TargetAudience}}
<target audience>
	{{ .TargetAudience }}
</target audience>
{{end}}

{{if .Goal}}
<goal>
	{{ .Goal }}
</goal>
{{end}}

{{if .Steps}}
<steps>
{{range $index, $element := .Steps}}
	{{$indexPlusOne := add $index 1}}
	{{$indexPlusOne}}. {{$element}}
{{end}}
</steps>
{{end}}

{{if .OutputFormat}}
<output format>
	{{ .OutputFormat }}
</output format>
{{end}}

{{if .Tone}}
<tone>
	{{ .Tone }}
</tone>
{{end}}

	{{if .Context}}
<context>
	{{ .Context }}
</context>
{{end}}

{{if .Input}}
<input>
	{{ .Input }}
</input>
{{end}}
`

// ParseBasicTempl will return the BasicTemplate ready for use
func ParseBasicTempl() *template.Template {
	return template.Must(template.New("basic").Funcs(template.FuncMap{"add": AddFn}).Parse(BasicTemplate))
}
