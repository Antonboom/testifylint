package main

import (
	"fmt"
	"strings"
)

func ExpandCheck(check Check, pkg string, argValues []any) string {
	args := check.ArgsTmpl

	if len(args) > 0 {
		args = fmt.Sprintf(args, argValues...)
	}

	return strings.Join([]string{
		buildCheck(pkg, check.Fn, args, check.ReportedMsg),
		buildCheck(pkg, check.Fn, args+`, "msg"`, check.ReportedMsg),
		buildCheck(pkg, check.Fn, args+`, "msg with arg %d", 42`, check.ReportedMsg),
		buildCheck(pkg, withSuffixF(check.Fn), args+`, "msg"`, withSuffixF(check.ReportedMsg)),
		buildCheck(pkg, withSuffixF(check.Fn), args+`, "msg with arg %d", 42`, withSuffixF(check.ReportedMsg)),
	}, "\n")
}

func withSuffixF(s string) string {
	if s == "" {
		return s
	}
	return s + "f"
}

func buildCheck(pkg, fn string, args string, reportedMsg string) string {
	s := fmt.Sprintf("%s.%s(%s)", pkg, fn, args)
	if reportedMsg != "" {
		s += fmt.Sprintf(" // want %q", fmt.Sprintf(reportedMsg, pkg))
	}
	return s
}
