package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteMethodSignatureTestsGenerator struct{}

func (SuiteMethodSignatureTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteMethodSignature()
}

func (g SuiteMethodSignatureTestsGenerator) TemplateData() any {
	name := g.Checker().Name()

	return struct {
		CheckerName               CheckerName
		ReportMethodShouldBeClean string
		ReportConflictWithIface   string
	}{
		CheckerName:               CheckerName(name),
		ReportMethodShouldBeClean: QuoteReport(name + ": test method should not have any arguments or returning values"),
		ReportConflictWithIface:   QuoteReport(name + ": method conflicts with suite.%s interface"),
	}
}

func (SuiteMethodSignatureTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteMethodSignatureTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(SuiteMethodSignatureTestTmpl))
}

func (SuiteMethodSignatureTestsGenerator) GoldenTemplate() Executor {
	// NOTE(a.telyshev): Autofix may result in uncompiled code.
	return nil
}

const SuiteMethodSignatureTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidSuite struct {
	suite.Suite
}

func TestValidSuite(t *testing.T) {
	suite.Run(t, new(ValidSuite))
}

func (s *ValidSuite) SetupSuite() {}
func (s *ValidSuite) SetupTest() {}
func (s *ValidSuite) TearDownSuite() {}
func (s *ValidSuite) TearDownTest() {}
func (s *ValidSuite) BeforeTest(suiteName, testName string) {}
func (s *ValidSuite) AfterTest(suiteName, testName string) {}
func (s *ValidSuite) HandleStats(suiteName string, stats *suite.SuiteInformation) {}
func (s *ValidSuite) SetupSubTest() {}
func (s *ValidSuite) TearDownSubTest() {}
func (s *ValidSuite) TestTrue() { s.True(true) }

var (
	_ suite.SetupAllSuite     = (*ValidSuite)(nil)
	_ suite.SetupTestSuite    = (*ValidSuite)(nil)
	_ suite.TearDownAllSuite  = (*ValidSuite)(nil)
	_ suite.TearDownTestSuite = (*ValidSuite)(nil)
	_ suite.BeforeTest        = (*ValidSuite)(nil)
	_ suite.AfterTest         = (*ValidSuite)(nil)
	_ suite.WithStats         = (*ValidSuite)(nil)
	_ suite.SetupSubTest      = (*ValidSuite)(nil)
	_ suite.TearDownSubTest   = (*ValidSuite)(nil)
)

type InvalidSuite struct {
	suite.Suite
}

func TestInvalidSuite(t *testing.T) {
	suite.Run(t, new(InvalidSuite))
}

func (s *InvalidSuite) SetupSuite(_ bool) {} // want {{ printf $.ReportConflictWithIface "SetupAllSuite" }}
func (s *InvalidSuite) SetupTest() int { return 0 } // want {{ printf $.ReportConflictWithIface "SetupTestSuite" }}
func (s *InvalidSuite) TearDownSuite(_ bool, _ int) {} // want {{ printf $.ReportConflictWithIface "TearDownAllSuite" }}
func (s *InvalidSuite) TearDownTest() (string, bool) { return "", false } // want {{ printf $.ReportConflictWithIface "TearDownTestSuite" }}
func (s *InvalidSuite) BeforeTest(suiteName string, testName int) {} // want {{ printf $.ReportConflictWithIface "BeforeTest" }}
func (s *InvalidSuite) AfterTest(suiteName int, testName string) {} // want {{ printf $.ReportConflictWithIface "AfterTest" }}
func (s *InvalidSuite) HandleStats(suiteName string, stats suite.SuiteInformation) {} // want {{ printf $.ReportConflictWithIface "WithStats" }}
func (s *InvalidSuite) SetupSubTest(_ string) {} // want {{ printf $.ReportConflictWithIface "SetupSubTest" }}
func (s *InvalidSuite) TearDownSubTest(ss string) string { return ss } // want {{ printf $.ReportConflictWithIface "TearDownSubTest" }}

func (s *InvalidSuite) TestTrue() { s.True(true) }
func (s *InvalidSuite) Test1(t *testing.T) { s.True(true) } // want {{ $.ReportMethodShouldBeClean }}
func (s *InvalidSuite) Test2() bool { return s.True(true) } // want {{ $.ReportMethodShouldBeClean }}
func (s *InvalidSuite) Test3(_ *testing.T) bool { return s.True(true) } // want {{ $.ReportMethodShouldBeClean }}

type MixedSuite struct {
	suite.Suite
}

func TestMixedSuite(t *testing.T) {
	suite.Run(t, new(MixedSuite))
}

func (s *MixedSuite) SetupSuite() {}
func (s *MixedSuite) SetupTest() int { return 0 } // want {{ printf $.ReportConflictWithIface "SetupTestSuite" }}
func (s *MixedSuite) TearDownSuite() {}
func (s MixedSuite) TearDownTest() (string, bool) { return "", false } // Value receivers are not supported.
func (s MixedSuite) BeforeTest(suiteName, testName string) {}
func (s MixedSuite) AfterTest(suiteName, testName string) {}
func (s MixedSuite) HandleStats(suiteName string, stats suite.SuiteInformation) {} // Value receivers are not supported.
func (s *MixedSuite) SetupSubTest() {}
func (s *MixedSuite) TearDownSubTest() {}

func (s *MixedSuite) TestTrue() { s.True(true) }
func (s *MixedSuite) TestFalse(t *testing.T) { s.False(false) } // want {{ $.ReportMethodShouldBeClean }}
`
