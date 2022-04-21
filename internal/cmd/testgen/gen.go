package main

import "text/template"

type TestsGenerator interface {
	Template() *template.Template
	Data() any
}

type Check struct {
	Fn              string   // "Equal"
	Args            []string // ["t", "len(%s)", "0"]
	DynamicArgIndex int      // 1
	ReportedMsg     string   // "use assert.Empty"
}
