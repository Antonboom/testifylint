package main

import (
	"fmt"
	"strings"
)

func ExpandCheck(check check, pkg string, dynamicArgValue string) string {
	args := make([]string, len(check.Args))
	copy(args, check.Args)
	args[check.DynamicArgIndex] = fmt.Sprintf(args[check.DynamicArgIndex], dynamicArgValue)

	return strings.Join([]string{
		buildCheck(pkg, check.Fn, args, check.ReportedMsg),
		buildCheck(pkg, check.Fn, append(args, `"msg"`), check.ReportedMsg),
		buildCheck(pkg, check.Fn, append(args, `"msg with arg %d"`, "42"), check.ReportedMsg),
		buildCheck(pkg, withSuffixF(check.Fn), append(args, `"msg"`), withSuffixF(check.ReportedMsg)),
		buildCheck(pkg, withSuffixF(check.Fn), append(args, `"msg with arg %d"`, "42"), withSuffixF(check.ReportedMsg)),
	}, "\n")
}

func withSuffixF(s string) string {
	if s == "" {
		return s
	}
	return s + "f"
}

func buildCheck(pkg, fn string, args []string, reportedMsg string) string {
	s := fmt.Sprintf("%s.%s(%s)", pkg, fn, strings.Join(args, ", "))
	if reportedMsg != "" {
		s += fmt.Sprintf("// want %q", fmt.Sprintf(reportedMsg, pkg))
	}
	return s
}
