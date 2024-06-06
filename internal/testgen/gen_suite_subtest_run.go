package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteSubtestRunTestsGenerator struct{}

func (SuiteSubtestRunTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteSubtestRun()
}

func (g SuiteSubtestRunTestsGenerator) TemplateData() any {
	var (
		checker     = g.Checker().Name()
		sReport     = QuoteReport(checker + ": use s.Run to run subtest")
		suiteReport = QuoteReport(checker + ": use suite.Run to run subtest")
	)

	return struct {
		CheckerName CheckerName
		SReport     string
		SuiteReport string
	}{
		CheckerName: CheckerName(checker),
		SReport:     sReport,
		SuiteReport: suiteReport,
	}
}

func (SuiteSubtestRunTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteSubtestRunTestsGenerator.ErroredTemplate").
		Parse(suiteSubtestRunTestTmpl))
}

func (SuiteSubtestRunTestsGenerator) GoldenTemplate() Executor {
	return nil
}

const suiteSubtestRunTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

{{ $suiteName := .CheckerName.AsSuiteName }}

type {{ $suiteName }} struct {
	suite.Suite
}

func Test{{ $suiteName }}(t *testing.T) {
	suite.Run(t, new({{ $suiteName }}))
}

func (s *{{ $suiteName }}) BeforeTest(suiteName, testName string) {
	s.Equal(1, 2)

	s.T().Run("init 1", func(t *testing.T) { // want {{ $.SReport }}
		s.Require().Equal(1, 2)
		assert.Equal(t, 1, 2)
	})

	s.Run("init 2", func() {
		s.Require().Equal(2, 1)
	})
}

func (s *{{ $suiteName }}) TestOne() {
	s.Equal(1, 2)

	s.T().Run("init 1", func(t *testing.T) { // want {{ $.SReport }}
		s.Require().Equal(11, 22)

		s.Run("init 1.1", func() {
			s.T().Run("init 1.2", func(t *testing.T) { // want {{ $.SReport }}
				s.Equal(111, 222)
				assert.Equal(t, 111, 222)
			})
		})
	})

	s.Run("init 2", func() {
		s.Equal(2, 1)
	})
}

func (suite *{{ $suiteName }}) TestTwo() {
	suite.Equal(1, 2)

	suite.Run("init 1", func() {
		suite.Equal(1, 2)
	})

	suite.Run("init 2", func() {
		suite.Equal(2, 1)

		suite.T().Run("init 2.1", func(t *testing.T) { // want {{ $.SuiteReport }}
			suite.Run("init 2.2", func() {
				suite.Require().Equal(222, 111)
			})
			assert.Equal(t, 22, 11)
		})
	})

	cases := []struct{ Name string }{}
	for _, tt := range cases {
		suite.T().Run(tt.Name, func(t *testing.T) { // want {{ $.SuiteReport }}
			suite.Equal(1, 2)
		})
	}
}
`
