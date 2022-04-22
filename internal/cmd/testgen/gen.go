package main

import "text/template"

type TestsGenerator interface {
	Template() *template.Template
	Data() any
}

type Check struct {
	Fn          string // "Equal"
	ArgsTmpl    string // "t, "len(%s), "0"
	ReportedMsg string // "use %s.Empty"
	MsgAsIs     bool
}

// TODO: валидация ArgsTmpl, ".*, .*, .*"
