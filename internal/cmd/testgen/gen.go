package main

import "text/template"

type TestsGenerator interface {
	Data() any
	ErroredTemplate() *template.Template
	GoldenTemplate() *template.Template
}

type Check struct {
	Fn    string // "Equal"
	Argsf string // "%s, %s"

	ReportMsgf string //  "use %s.%s"

	ProposedFn    string // "InDelta"
	ProposedArgsf string // %s, %s, 0.0001"
}
