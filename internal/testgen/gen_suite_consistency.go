package main

import (
	"regexp"
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteConsistencyTestsGenerator struct{}

func (SuiteConsistencyTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteConsistency()
}

func (g SuiteConsistencyTestsGenerator) TemplateData() any {
	name := g.Checker().Name()

	return struct {
		CheckerName        CheckerName
		RunnerReport       string
		ReceiverReport     string
		UnnamedRcvReport   string
		MethodReport       string
		DefaultNamePattern string
	}{
		CheckerName:      CheckerName(name),
		RunnerReport:     QuoteReport(name + ": suite test function name %s does not match suite name (expected %s or %s_<Variant>)"),
		ReceiverReport:   QuoteReport(name + ": suite receiver name %s does not match configured name %s"),
		UnnamedRcvReport: QuoteReport(name + ": suite receiver should be named %s"),
		MethodReport:     QuoteReport(name + ": suite test method name %s does not match pattern %s"),
		DefaultNamePattern: strings.ReplaceAll(
			regexp.QuoteMeta(checkers.DefaultSuiteConsistencyTestNamePattern.String()),
			`\`,
			`\\`,
		),
	}
}

func (SuiteConsistencyTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteConsistencyTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(suiteConsistencyTestTmpl))
}

func (SuiteConsistencyTestsGenerator) GoldenTemplate() Executor {
	golden := strings.NewReplacer(
		`func TestBadName_Run(t *testing.T) { // want {{ printf $.RunnerReport "TestBadName_Run" "TestBadNameSuite" "TestBadNameSuite" }}`,
		`func TestBadNameSuite(t *testing.T) { // want {{ printf $.RunnerReport "TestBadName_Run" "TestBadNameSuite" "TestBadNameSuite" }}`,
		`var badRunner = TestBadName_Run`,
		`var badRunner = TestBadNameSuite`,
		`func TestCompletelyDifferent(t *testing.T) { // want {{ printf $.RunnerReport "TestCompletelyDifferent" "TestBadNameSuite2" "TestBadNameSuite2" }}`,
		`func TestBadNameSuite2(t *testing.T) { // want {{ printf $.RunnerReport "TestCompletelyDifferent" "TestBadNameSuite2" "TestBadNameSuite2" }}`,
		`func (suite *ReceiverSuite) TestUsecase() { // want {{ printf $.ReceiverReport "suite" "s" }}`,
		`func (s *ReceiverSuite) TestUsecase() { // want {{ printf $.ReceiverReport "suite" "s" }}`,
		`suite.Equal(1, 1)`,
		`s.Equal(1, 1)`,
		`func (suite *ReceiverSuite) Test_Usecase_Success() { // want {{ printf $.ReceiverReport "suite" "s" }} {{ printf $.MethodReport "Test_Usecase_Success" $.DefaultNamePattern }}`,
		`func (s *ReceiverSuite) TestUsecase_Success() { // want {{ printf $.ReceiverReport "suite" "s" }} {{ printf $.MethodReport "Test_Usecase_Success" $.DefaultNamePattern }}`,
		`suite.NotZero(1)`,
		`s.NotZero(1)`,
		`var badMethod = (*ReceiverSuite).Test_Usecase_Success`,
		`var badMethod = (*ReceiverSuite).TestUsecase_Success`,
		`func (suite *ReceiverSuite) helper() { // want {{ printf $.ReceiverReport "suite" "s" }}`,
		`func (s *ReceiverSuite) helper() { // want {{ printf $.ReceiverReport "suite" "s" }}`,
		`suite.True(true)`,
		`s.True(true)`,
	).Replace(suiteConsistencyTestTmpl)
	return template.Must(template.New("SuiteConsistencyTestsGenerator.GoldenTemplate").
		Funcs(fm).
		Parse(golden))
}

const suiteConsistencyTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Bad: runner name does not match suite name.
type BadNameSuite struct {
	suite.Suite
}

func TestBadName_Run(t *testing.T) { // want {{ printf $.RunnerReport "TestBadName_Run" "TestBadNameSuite" "TestBadNameSuite" }}
	suite.Run(t, new(BadNameSuite))
}

var badRunner = TestBadName_Run

// Bad: runner name is completely different.
type BadNameSuite2 struct {
	suite.Suite
}

func TestCompletelyDifferent(t *testing.T) { // want {{ printf $.RunnerReport "TestCompletelyDifferent" "TestBadNameSuite2" "TestBadNameSuite2" }}
	suite.Run(t, &BadNameSuite2{})
}

// Bad: runner rename would conflict.
type ConflictRunnerSuite struct {
	suite.Suite
}

func TestConflictRunnerSuite(t *testing.T) {}

func TestConflictRunner_Run(t *testing.T) { // want {{ printf $.RunnerReport "TestConflictRunner_Run" "TestConflictRunnerSuite" "TestConflictRunnerSuite" }}
	suite.Run(t, new(ConflictRunnerSuite))
}

// Bad: runner suffix is malformed, so autofix should not drop the intended variant.
type BadVariantSuite struct {
	suite.Suite
}

func TestBadVariantSuite_noSignature(t *testing.T) { // want {{ printf $.RunnerReport "TestBadVariantSuite_noSignature" "TestBadVariantSuite" "TestBadVariantSuite" }}
	suite.Run(t, new(BadVariantSuite))
}

// Good: runner name matches suite name.
type GoodNameSuite struct {
	suite.Suite
}

func TestGoodNameSuite(t *testing.T) {
	suite.Run(t, new(GoodNameSuite))
}

// Good: runner variants are allowed.
type VariantSuite struct {
	suite.Suite
}

func TestVariantSuite_NoSignature(t *testing.T) {
	suite.Run(t, new(VariantSuite))
}

func TestVariantSuite_SignedVerdicts(t *testing.T) {
	suite.Run(t, new(VariantSuite))
}

func TestVariantSuite_SignedVerdicts_WithoutCache(t *testing.T) {
	suite.Run(t, new(VariantSuite))
}

// Good: suite instance is stored in a variable.
type VariableSuite struct {
	suite.Suite
}

func TestVariableSuite(t *testing.T) {
	s := new(VariableSuite)
	suite.Run(t, s)
}

// Good: generic suite instance.
type GenericSuite[T any] struct {
	suite.Suite
}

func TestGenericSuite(t *testing.T) {
	suite.Run(t, new(GenericSuite[int]))
}

// Good: suite.Run in nested callback is ignored.
type CallbackSuite struct {
	suite.Suite
}

func TestCallbackWrapper(t *testing.T) {
	callback := func() {
		suite.Run(t, new(CallbackSuite))
	}
	callback()
}

// Good: non-test helper function that wraps suite.Run is ignored.
type HelperSuite struct {
	suite.Suite
}

func runHelperSuite(t *testing.T) {
	suite.Run(t, new(HelperSuite))
}

// Bad: receiver name is inconsistent.
type ReceiverSuite struct {
	suite.Suite
}

func (suite *ReceiverSuite) TestUsecase() { // want {{ printf $.ReceiverReport "suite" "s" }}
	suite.Equal(1, 1)
}

func (s *ReceiverSuite) TestUsecaseSuccess() {}

func (s *ReceiverSuite) TestUsecase_SubTest() {}

func (s *ReceiverSuite) TestUsecase_SubTest_AnotherCase() {}

func (suite *ReceiverSuite) Test_Usecase_Success() { // want {{ printf $.ReceiverReport "suite" "s" }} {{ printf $.MethodReport "Test_Usecase_Success" $.DefaultNamePattern }}
	suite.NotZero(1)
}

var badMethod = (*ReceiverSuite).Test_Usecase_Success

func (s *ReceiverSuite) Testusecase() { // want {{ printf $.MethodReport "Testusecase" $.DefaultNamePattern }}
}

func (s *ReceiverSuite) TestUsecase__SubTest() { // want {{ printf $.MethodReport "TestUsecase__SubTest" $.DefaultNamePattern }}
}

func (s *ReceiverSuite) TestUsecase_() { // want {{ printf $.MethodReport "TestUsecase_" $.DefaultNamePattern }}
}

func (s *ReceiverSuite) TestAction() {}

func (s *ReceiverSuite) Test_Action() { // want {{ printf $.MethodReport "Test_Action" $.DefaultNamePattern }}
}

func (*ReceiverSuite) TestUnnamedReceiver() { // want {{ printf $.UnnamedRcvReport "s" }}
}

func (suite *ReceiverSuite) helperWithConflict(s string) { // want {{ printf $.ReceiverReport "suite" "s" }}
	suite.Equal(s, s)
}

func (suite *ReceiverSuite) helper() { // want {{ printf $.ReceiverReport "suite" "s" }}
	suite.True(true)
}

// Good: non-suite receiver is ignored.
type NonSuite struct{}

func (suite *NonSuite) TestSomething() {}

// Good: indirect suite embedding is supported.
type EmbeddedBaseSuite struct {
	suite.Suite
}

type IndirectSuite struct {
	EmbeddedBaseSuite
}

func (s *IndirectSuite) TestIndirectEmbedding() {}
`
