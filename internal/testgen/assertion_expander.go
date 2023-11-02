package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type expandMode int

const (
	expandModeFull expandMode = iota
	expandModeNotFmtSet
	expandModeFmtSet
	expandModeNotFmtSingle
	expandModeFmtSingle
	expandModeExtreme
)

type AssertionExpander struct {
	mode     expandMode
	asGolden bool
}

func NewAssertionExpander() *AssertionExpander {
	return &AssertionExpander{
		asGolden: false,
		mode:     expandModeExtreme,
	}
}

// AsGolden turns AssertionExpander into "golden file" mode.
// This means that assertions from Expand will be constructed based on Assertion.Proposed* fields.
func (e *AssertionExpander) AsGolden() *AssertionExpander {
	e.asGolden = true
	return e
}

func (e *AssertionExpander) FullMode() *AssertionExpander {
	e.mode = expandModeFull
	return e
}

func (e *AssertionExpander) FmtSingleMode() *AssertionExpander {
	e.mode = expandModeFmtSingle
	return e
}

func (e *AssertionExpander) NotFmtSetMode() *AssertionExpander {
	e.mode = expandModeNotFmtSet
	return e
}

func (e *AssertionExpander) NotFmtSingleMode() *AssertionExpander {
	e.mode = expandModeNotFmtSingle
	return e
}

// Expand expands incoming Assertion into several assertion strings according to the rules:
//
//   - build arguments string: testingTParam + Assertion.Argsf + argValues, e.g.
//     "t" + "%s, 0" + "arr" = "t, arr, 0"
//
//   - build assertion string: selector + Assertion.Fn + arguments, e.g.
//     "assert" + "Len" + "t, arr, 0" = "assert.Len(t, arr, 0)"
//
//   - if Assertion.ReportMsgf defined then append analysistest.Run diagnostic comment, e.g.
//     "assert.Len(t, arr, 0)" + "empty: use %s.%s" = `assert.Len(t, arr, 0) // want "empty: use assert\\.Empty"`
//
// For diagnostic comment formatting Expand uses Assertion.ProposedSelector and Assertion.ProposedFn
// or selector and Assertion.Fn if one of them is not defined. If you do not need formatting use `%.s` (or `%.0s`).
//
// The final string is copied in several variants, depending on the current AssertionExpander mode.
func (e *AssertionExpander) Expand(assrn Assertion, selector, testingTParam string, argValues []any) string {
	fn, args := assrn.Fn, assrn.Argsf

	if e.asGolden {
		selector = or(assrn.ProposedSelector, selector)
		fn = or(assrn.ProposedFn, fn)
		args = or(assrn.ProposedArgsf, args)
	}

	if testingTParam != "" {
		args = testingTParam + ", " + args
	}
	if len(argValues) > 0 {
		args = fmt.Sprintf(args, argValues...)
	}

	as := map[expandMode][]string{ // as is for "assertionSets".
		expandModeNotFmtSet: {
			buildAssertion(selector, fn, args, assrn.ReportMsgf, assrn.ProposedSelector, assrn.ProposedFn),
		},
	}
	for _, newArgs := range []string{
		`, "msg"`,
		`, "msg with arg %d", 42`,
		`, "msg with args %d %s", 42, "42"`,
	} {
		as[expandModeNotFmtSet] = append(as[expandModeNotFmtSet],
			buildAssertion(selector, fn, args+newArgs, assrn.ReportMsgf, assrn.ProposedSelector, assrn.ProposedFn))

		as[expandModeFmtSet] = append(as[expandModeFmtSet],
			buildFmtAssertion(selector, fn, args+newArgs, assrn.ReportMsgf, assrn.ProposedSelector, assrn.ProposedFn))
	}
	as[expandModeFull] = append(slices.Clone(as[expandModeNotFmtSet]), as[expandModeFmtSet]...)
	as[expandModeNotFmtSingle] = []string{as[expandModeNotFmtSet][0]}
	as[expandModeFmtSingle] = []string{as[expandModeFmtSet][len(as[expandModeFmtSet])-1]}
	as[expandModeExtreme] = append(slices.Clone(as[expandModeNotFmtSingle]), as[expandModeFmtSingle]...)

	assertions, ok := as[e.mode]
	if !ok {
		panic(fmt.Sprintf("unsupported expand mode: %v", e.mode))
	}

	return strings.Join(assertions, "\n")
}

func buildAssertion(selector, fn, args, reportedMsgf, proposedSel, proposedFn string) string {
	s := fmt.Sprintf("%s.%s(%s)", selector, fn, args)

	if reportedMsgf != "" {
		if or(proposedSel, proposedFn) != "" {
			reportedMsgf = fmt.Sprintf(reportedMsgf, or(proposedSel, selector), or(proposedFn, fn))
		}
		s += " // want " + QuoteReport(reportedMsgf)
	}
	return s
}

func buildFmtAssertion(selector, fn, args, reportedMsgf, proposedSel, proposedFn string) string {
	return buildAssertion(selector, withSuffixF(fn), args, reportedMsgf, proposedSel, withSuffixF(proposedFn))
}

func or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func withSuffixF(s string) string {
	if s == "" {
		return s
	}
	return s + "f"
}

func QuoteReport(msg string) string {
	return fmt.Sprintf("%q", regexp.QuoteMeta(msg))
}
