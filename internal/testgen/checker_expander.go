package main

import (
	"fmt"
	"regexp"
	"strings"
)

// TODO -> comments to types and methods

type CheckerExpander struct {
	toMultiple bool
	withFFuncs bool
	asGolden   bool
}

func NewCheckerExpander() *CheckerExpander {
	return &CheckerExpander{
		toMultiple: false,
		withFFuncs: true,
		asGolden:   false,
	}
}

func (e *CheckerExpander) WithoutFFuncs() *CheckerExpander {
	e.withFFuncs = false
	return e
}

func (e *CheckerExpander) AsGolden() *CheckerExpander {
	e.asGolden = true
	return e
}

func (e *CheckerExpander) ToMultiple() *CheckerExpander {
	e.toMultiple = true
	return e
}

func (e *CheckerExpander) Expand(check Check, selector, testingTParam string, argValues []string) string {
	fn, args := check.Fn, check.Argsf

	if e.asGolden {
		if check.ProposedSelector != "" {
			selector = check.ProposedSelector
		}
		if check.ProposedFn != "" {
			fn = check.ProposedFn
		}
		if check.ProposedArgsf != "" {
			args = check.ProposedArgsf
		}
	}

	if testingTParam != "" {
		args = testingTParam + ", " + args
	}
	if len(argValues) > 0 {
		args = fmt.Sprintf(args, toSliceOfAny(argValues)...)
	}

	checks := []string{
		buildCheck(selector, fn, args, check.ReportMsgf, check.ProposedFn, false),
	}

	if e.toMultiple {
		checks = append(checks,
			buildCheck(selector, fn, args+`, "msg"`, check.ReportMsgf, check.ProposedFn, false),
			buildCheck(selector, fn, args+`, "msg with arg %d", 42`, check.ReportMsgf, check.ProposedFn, false),
		)

		if e.withFFuncs {
			checks = append(checks, []string{
				buildCheck(selector, fn, args+`, "msg"`, check.ReportMsgf, check.ProposedFn, true),
				buildCheck(selector, fn, args+`, "msg with arg %d", 42`, check.ReportMsgf, check.ProposedFn, true),
			}...)
		}
	}
	return strings.Join(checks, "\n")
}

func buildCheck(selector, fn, args, reportedMsgf, proposedFn string, withFSuffix bool) string {
	if withFSuffix {
		fn = withSuffixF(fn)
		proposedFn = withSuffixF(proposedFn)
	}

	s := fmt.Sprintf("%s.%s(%s)", selector, fn, args)

	if reportedMsgf != "" {
		if proposedFn != "" {
			reportedMsgf = fmt.Sprintf(reportedMsgf, selector, proposedFn)
		}
		s += fmt.Sprintf(" // want %q", regexp.QuoteMeta(reportedMsgf))
	}
	return s
}

func withSuffixF(s string) string {
	if s == "" {
		return s
	}
	return s + "f"
}

func toSliceOfAny(in []string) []any {
	result := make([]any, len(in))
	for i := range in {
		result[i] = in[i]
	}
	return result
}
