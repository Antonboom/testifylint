package main

// Assertion is a generic view of testify assertion.
// Designed to be expanded (by AssertionExpander) to multiple lines of assertions.
type Assertion struct {
	Fn    string // "Equal"
	Argsf string // "%s, %s"

	ReportMsgf string //  "use %s.%s"

	ProposedSelector string // s.Require()
	ProposedFn       string // "InDelta"
	ProposedArgsf    string // %s, %s, 0.0001"
}

// WithoutReport strips Assertion expected warning.
func (a Assertion) WithoutReport() Assertion {
	return Assertion{Fn: a.Fn, Argsf: a.Argsf}
}
