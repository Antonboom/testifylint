package main

import (
	"fmt"
	"strings"
)

func ExpandCheck(check Check, pkg string, argValues []any) string {
	args := check.ArgsTmpl

	if len(argValues) > 0 {
		args = fmt.Sprintf(args, argValues...)
	}

	return strings.Join([]string{
		buildCheck(pkg, check.Fn, args, check.ReportedMsg, check.MsgAsIs),
		buildCheck(pkg, check.Fn, args+`, "msg"`, check.ReportedMsg, check.MsgAsIs),
		buildCheck(pkg, check.Fn, args+`, "msg with arg %d", 42`, check.ReportedMsg, check.MsgAsIs),
		buildCheck(pkg, withSuffixF(check.Fn, check.MsgAsIs), args+`, "msg"`, withSuffixF(check.ReportedMsg, check.MsgAsIs), check.MsgAsIs),
		buildCheck(pkg, withSuffixF(check.Fn, check.MsgAsIs), args+`, "msg with arg %d", 42`, withSuffixF(check.ReportedMsg, check.MsgAsIs), check.MsgAsIs),
	}, "\n")
}

func withSuffixF(s string, asIs bool) string {
	if s == "" || asIs {
		return s
	}
	return s + "f"
}

func buildCheck(pkg, fn string, args string, reportedMsg string, asIs bool) string {
	s := fmt.Sprintf("%s.%s(%s)", pkg, fn, args)
	if reportedMsg != "" {
		if !asIs {
			reportedMsg = fmt.Sprintf(reportedMsg, pkg)
		}
		s += fmt.Sprintf(" // want %q", reportedMsg)
	}
	return s
}
