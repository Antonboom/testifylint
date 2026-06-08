package main

import "strings"

// CheckerName is a simple helper type useful for checker name transformations.
type CheckerName string

// AsPkgName transforms "suite-extra-assert-call" into "suiteextraassertcall".
func (n CheckerName) AsPkgName() string {
	return strings.ReplaceAll(string(n), "-", "")
}

// AsTestName transforms "suite-extra-assert-call" into "TestSuiteExtraAssertCallChecker".
func (n CheckerName) AsTestName() string {
	return "Test" + n.toCamelCase() + "Checker"
}

// AsSuiteName transforms "suite-extra-assert-call" into "SuiteExtraAssertCallCheckerSuite".
func (n CheckerName) AsSuiteName() string {
	return n.toCamelCase() + "CheckerSuite"
}

func (n CheckerName) toCamelCase() string {
	var result string
	for _, word := range strings.Split(string(n), "-") {
		result += title(word)
	}
	return result
}

// title is more simple analogue of deprecated strings.Title.
func title(word string) string {
	if ch := word[0]; ch >= 'a' && ch <= 'z' {
		return string(ch-('a'-'A')) + word[1:]
	}
	return word
}
