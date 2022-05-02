package main

import "text/template"

type TestsGenerator interface {
	Template() *template.Template
	Data() any
}

type Check struct {
	Fn                       string // "Equal"
	ArgsTmpl                 string // "t, "len(%s), "0"
	ReportedMsgf, ProposedFn string // "use %s.%s" "Empty"
	ReportedMsg              string // "need to simplify the check"
}

// TODO: валидация ArgsTmpl, ".*, .*, .*"
// TODO: валидация или reportMsgF+ProposedFn или ReportedMsg
// TODO: ArgsTmpl -> Argsf

func (c Check) Validate() error {
	return nil
}
