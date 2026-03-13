package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteTestNameTestsGenerator struct{}

func (SuiteTestNameTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteTestName()
}

func (g SuiteTestNameTestsGenerator) TemplateData() any {
	name := g.Checker().Name()

	return struct {
		CheckerName CheckerName
		Report      string
	}{
		CheckerName: CheckerName(name),
		Report:      QuoteReport(name + ": suite test function name %s does not match suite name (expected %s)"),
	}
}

func (SuiteTestNameTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteTestNameTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(suiteTestNameTestTmpl))
}

func (SuiteTestNameTestsGenerator) GoldenTemplate() Executor {
	// The fix only renames the function identifier, so // want comments remain in the fixed output.
	golden := strings.NewReplacer(
		`func TestBadNameSuite_Run(t *testing.T) { // want {{ printf $.Report "TestBadNameSuite_Run" "TestBadNameSuite" }}`,
		`func TestBadNameSuite(t *testing.T) { // want {{ printf $.Report "TestBadNameSuite_Run" "TestBadNameSuite" }}`,
		`func TestCompletelyDifferent(t *testing.T) { // want {{ printf $.Report "TestCompletelyDifferent" "TestBadNameSuite2" }}`,
		`func TestBadNameSuite2(t *testing.T) { // want {{ printf $.Report "TestCompletelyDifferent" "TestBadNameSuite2" }}`,
	).Replace(suiteTestNameTestTmpl)
	return template.Must(template.New("SuiteTestNameTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(golden))
}

const suiteTestNameTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Bad: runner name does not match suite name.
type BadNameSuite struct {
	suite.Suite
}

func TestBadNameSuite_Run(t *testing.T) { // want {{ printf $.Report "TestBadNameSuite_Run" "TestBadNameSuite" }}
	suite.Run(t, new(BadNameSuite))
}

// Bad: runner name is completely different.
type BadNameSuite2 struct {
	suite.Suite
}

func TestCompletelyDifferent(t *testing.T) { // want {{ printf $.Report "TestCompletelyDifferent" "TestBadNameSuite2" }}
	suite.Run(t, new(BadNameSuite2))
}

// Good: runner name matches suite name (new form).
type GoodNameSuite struct {
	suite.Suite
}

func TestGoodNameSuite(t *testing.T) {
	suite.Run(t, new(GoodNameSuite))
}

// Good: runner name matches suite name (address-of form).
type GoodNameSuite2 struct {
	suite.Suite
}

func TestGoodNameSuite2(t *testing.T) {
	suite.Run(t, &GoodNameSuite2{})
}

// Good: non-test helper function that wraps suite.Run – should not be reported.
type GoodNameSuite3 struct {
	suite.Suite
}

func runGoodNameSuite3(t *testing.T) {
	suite.Run(t, new(GoodNameSuite3))
}

// Good: suite method wrapping suite.Run – should not be reported (has a receiver).
type GoodNameSuite4 struct {
	suite.Suite
}

func TestGoodNameSuite4(t *testing.T) {
	suite.Run(t, new(GoodNameSuite4))
}

func (s *GoodNameSuite4) helperRun(t *testing.T) {
	suite.Run(t, new(GoodNameSuite4))
}
`
