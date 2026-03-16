package main

import (
	"strings"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type GracefulTeardownTestsGenerator struct{}

func (GracefulTeardownTestsGenerator) Checker() checkers.Checker {
	return checkers.NewGracefulTeardown()
}

func (g GracefulTeardownTestsGenerator) TemplateData() any {
	var (
		name   = g.Checker().Name()
		report = QuoteReport(name + ": do not use require in cleanup code")
	)

	return struct {
		CheckerName CheckerName
		Report      string
	}{
		CheckerName: CheckerName(name),
		Report:      report,
	}
}

func (GracefulTeardownTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("GracefulTeardownTestsGenerator.ErroredTemplate").
		Parse(gracefulTeardownTestTmpl))
}

func (GracefulTeardownTestsGenerator) GoldenTemplate() Executor {
	golden := strings.ReplaceAll(gracefulTeardownTestTmpl,
		"\t\trequire.NoError(t, err) // want {{ $.Report }}",
		"\t\tassert.NoError(t, err) // want {{ $.Report }}")
	golden = strings.ReplaceAll(golden,
		"\t\trequire.Nil(t, err)     // want {{ $.Report }}",
		"\t\tassert.Nil(t, err)     // want {{ $.Report }}")
	golden = strings.ReplaceAll(golden,
		"\ts.Require().NoError(nil) // want {{ $.Report }}",
		"\ts.Assert().NoError(nil) // want {{ $.Report }}")
	return template.Must(template.New("GracefulTeardownTestsGenerator.GoldenTemplate").Parse(golden))
}

const gracefulTeardownTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	var err error

	// OK: require in test body.
	require.NoError(t, err)

	// OK: assert in cleanup.
	t.Cleanup(func() {
		assert.NoError(t, err)
	})

	// Bad: require in cleanup.
	t.Cleanup(func() {
		require.NoError(t, err) // want {{ $.Report }}
	})

	// Bad: multiple requires in cleanup.
	t.Cleanup(func() {
		require.NoError(t, err) // want {{ $.Report }}
		require.Nil(t, err)     // want {{ $.Report }}
	})
}

{{ $suiteName := .CheckerName.AsSuiteName }}

type {{ $suiteName }} struct {
	suite.Suite
}

func Test{{ $suiteName }}(t *testing.T) {
	suite.Run(t, new({{ $suiteName }}))
}

func (s *{{ $suiteName }}) SetupTest() {
	// OK: require in setup.
	s.Require().NoError(nil)
}

func (s *{{ $suiteName }}) TestSomething() {
	// OK: require in test method.
	s.Require().NoError(nil)
}

func (s *{{ $suiteName }}) TearDownTest() {
	// OK: assert in teardown.
	s.Assert().NoError(nil)
	s.NoError(nil)

	// Bad: require in teardown.
	s.Require().NoError(nil) // want {{ $.Report }}
}

func (s *{{ $suiteName }}) TearDownSuite() {
	// Bad: require in teardown suite.
	s.Require().NoError(nil) // want {{ $.Report }}
}

func (s *{{ $suiteName }}) AfterTest(_, _ string) {
	// Bad: require in AfterTest.
	s.Require().NoError(nil) // want {{ $.Report }}
}
`
