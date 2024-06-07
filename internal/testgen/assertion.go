package main

import "sort"

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

var allAssertions = []Assertion{
	{Fn: "Condition", Argsf: "nil"},
	{Fn: "Contains", Argsf: "nil, nil"},
	{Fn: "DirExists", Argsf: `""`},
	{Fn: "ElementsMatch", Argsf: "nil, nil"},
	{Fn: "Empty", Argsf: "nil"},
	{Fn: "Equal", Argsf: "nil, nil"},
	{Fn: "EqualError", Argsf: `nil, ""`},
	{Fn: "EqualExportedValues", Argsf: "nil, nil"},
	{Fn: "EqualValues", Argsf: "nil, nil"},
	{Fn: "Error", Argsf: "nil"},
	{Fn: "ErrorAs", Argsf: "nil, nil"},
	{Fn: "ErrorContains", Argsf: `nil, ""`},
	{Fn: "ErrorIs", Argsf: "nil, nil"},
	{Fn: "Eventually", Argsf: "nil, 0, 0"},
	{Fn: "EventuallyWithT", Argsf: "nil, 0, 0"},
	{Fn: "Exactly", Argsf: "nil, nil"},
	{Fn: "Fail", Argsf: `""`},
	{Fn: "FailNow", Argsf: `""`},
	{Fn: "False", Argsf: "false"},
	{Fn: "FileExists", Argsf: `""`},
	{Fn: "Implements", Argsf: "nil, nil"},
	{Fn: "InDelta", Argsf: "0., 0., 0."},
	{Fn: "InDeltaMapValues", Argsf: "nil, nil, 0."},
	{Fn: "InDeltaSlice", Argsf: "nil, nil, 0."},
	{Fn: "InEpsilon", Argsf: "nil, nil, 0."},
	{Fn: "InEpsilonSlice", Argsf: "nil, nil, 0."},
	{Fn: "IsType", Argsf: "nil, nil"},
	{Fn: "JSONEq", Argsf: `"", ""`},
	{Fn: "Len", Argsf: "nil, 0"},
	{Fn: "Never", Argsf: "nil, 0, 0"},
	{Fn: "Nil", Argsf: "nil"},
	{Fn: "NoDirExists", Argsf: `""`},
	{Fn: "NoError", Argsf: "nil"},
	{Fn: "NoFileExists", Argsf: `""`},
	{Fn: "NotContains", Argsf: "nil, nil"},
	{Fn: "NotEmpty", Argsf: "nil"},
	{Fn: "NotEqual", Argsf: "nil, nil"},
	{Fn: "NotEqualValues", Argsf: "nil, nil"},
	{Fn: "NotErrorIs", Argsf: "nil, nil"},
	{Fn: "NotNil", Argsf: "nil"},
	{Fn: "NotPanics", Argsf: "nil"},
	{Fn: "NotRegexp", Argsf: `nil, ""`},
	{Fn: "NotSame", Argsf: "nil, nil"},
	{Fn: "NotSubset", Argsf: "nil, nil"},
	{Fn: "NotZero", Argsf: "nil"},
	{Fn: "Panics", Argsf: "nil"},
	{Fn: "PanicsWithError", Argsf: `"", nil`},
	{Fn: "PanicsWithValue", Argsf: "nil, nil"},
	{Fn: "Regexp", Argsf: "nil, nil"},
	{Fn: "Same", Argsf: "nil, nil"},
	{Fn: "Subset", Argsf: "nil, nil"},
	{Fn: "True", Argsf: "true"},
	{Fn: "WithinDuration", Argsf: "time.Time{}, time.Time{}, 0"},
	{Fn: "WithinRange", Argsf: "time.Time{}, time.Time{}, time.Time{}"},
	{Fn: "YAMLEq", Argsf: `"", ""`},
	{Fn: "Zero", Argsf: "nil"},
}

func sortAssertions(assertions []Assertion) {
	sort.Slice(assertions, func(i, j int) bool {
		lhs, rhs := assertions[i], assertions[j]
		if lhs.Fn == rhs.Fn {
			return lhs.Argsf < rhs.Argsf
		}
		return lhs.Fn < rhs.Fn
	})
}
