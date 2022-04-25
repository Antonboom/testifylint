package main

import (
	"fmt"
	"strings"
)

func Product(a, b []any) [][]any {
	result := make([][]any, 0, len(a))
	for _, v := range a {
		for _, vv := range b {
			result = append(result, []any{v, vv})
		}
	}
	return result
}

func ExpandCheck(check Check, pkg string, argValues []any) (string, error) { // TODO: refactoring
	if err := check.Validate(); err != nil {
		return "", err
	}

	args := check.ArgsTmpl
	if len(argValues) > 0 {
		args = fmt.Sprintf(args, argValues...)
	}

	report := func(withF bool) string {
		switch {
		case check.ReportedMsg != "":
			return check.ReportedMsg
		case check.ReportedMsgf != "":
			fn := check.ProposedFn
			if withF {
				fn = withSuffixF(fn)
			}
			return fmt.Sprintf(check.ReportedMsgf, pkg, fn)
		}
		return ""
	}

	return strings.Join([]string{
		buildCheck(pkg, check.Fn, args, report(false)),
		buildCheck(pkg, check.Fn, args+`, "msg"`, report(false)),
		buildCheck(pkg, check.Fn, args+`, "msg with arg %d", 42`, report(false)),
		buildCheck(pkg, withSuffixF(check.Fn), args+`, "msg"`, report(true)),
		buildCheck(pkg, withSuffixF(check.Fn), args+`, "msg with arg %d", 42`, report(true)),
	}, "\n"), nil
}

func ExpandCheckWithoutF(check Check, pkg string, argValues []any) (string, error) {
	if err := check.Validate(); err != nil {
		return "", err
	}

	args := check.ArgsTmpl
	if len(argValues) > 0 {
		args = fmt.Sprintf(args, argValues...)
	}

	report := func(withF bool) string {
		switch {
		case check.ReportedMsg != "":
			return check.ReportedMsg
		case check.ReportedMsgf != "":
			fn := check.ProposedFn
			if withF {
				fn = withSuffixF(fn)
			}
			return fmt.Sprintf(check.ReportedMsgf, pkg, fn)
		}
		return ""
	}

	return strings.Join([]string{
		buildCheck(pkg, check.Fn, args, report(false)),
		buildCheck(pkg, check.Fn, args+`, "msg"`, report(false)),
		buildCheck(pkg, check.Fn, args+`, "msg with arg %d", 42`, report(false)),
	}, "\n"), nil
}

func withSuffixF(s string) string {
	return s + "f"
}

func buildCheck(pkg, fn, args, reportedMsg string) string {
	s := fmt.Sprintf("%s.%s(%s)", pkg, fn, args)
	if reportedMsg != "" {
		s += fmt.Sprintf(" // want %q", reportedMsg)
	}
	return s
}
