package main

import "text/template"

var fm = template.FuncMap{
	"NewCheckerExpander": NewCheckerExpander,
}

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
	return "", nil
}
