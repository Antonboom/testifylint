package main

import (
	"fmt"
	"strings"
)

type CheckerExpander struct {
	withTArg   bool
	withFFuncs bool
	asGolden   bool
}

func NewCheckerExpander() *CheckerExpander {
	return &CheckerExpander{
		withTArg:   true,
		withFFuncs: true,
		asGolden:   false,
	}
}

func (e *CheckerExpander) WithoutTArg() *CheckerExpander {
	e.withTArg = false
	return e
}

func (e *CheckerExpander) WithoutFFuncs() *CheckerExpander {
	e.withFFuncs = false
	return e
}

func (e *CheckerExpander) AsGolden() *CheckerExpander {
	e.asGolden = true
	return e
}

func (e *CheckerExpander) Expand(check Check, selector string, argValues []any) string {
	fn, args := check.Fn, check.Argsf

	if e.asGolden {
		if check.ProposedFn != "" {
			fn = check.ProposedFn
		}
		if check.ProposedArgsf != "" {
			args = check.ProposedArgsf
		}
	}

	if e.withTArg {
		args = "t, " + args
	}
	if len(argValues) > 0 {
		args = fmt.Sprintf(args, argValues...)
	}

	checks := []string{
		buildCheck(selector, fn, args, check.ReportMsgf, check.ProposedFn, false),
		buildCheck(selector, fn, args+`, "msg"`, check.ReportMsgf, check.ProposedFn, false),
		buildCheck(selector, fn, args+`, "msg with arg %d", 42`, check.ReportMsgf, check.ProposedFn, false),
	}
	if e.withFFuncs {
		checks = append(checks, []string{
			buildCheck(selector, fn, args+`, "msg"`, check.ReportMsgf, check.ProposedFn, true),
			buildCheck(selector, fn, args+`, "msg with arg %d", 42`, check.ReportMsgf, check.ProposedFn, true),
		}...)
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
		s += fmt.Sprintf(" // want %q", reportedMsgf)
	}
	return s
}

func withSuffixF(s string) string {
	if s == "" {
		return s
	}
	return s + "f"
}
