package main

import "strings"

type CheckerName string

func (n CheckerName) AsPkgName() string {
	return strings.ReplaceAll(string(n), "-", "")
}

func (n CheckerName) AsTestName() string {
	var result string
	for _, word := range strings.Split(string(n), "-") {
		result += strings.Title(word)
	}
	return "Test" + result + "Checker"
}
