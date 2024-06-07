package main

import (
	"regexp"
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type SuiteBrokenParallelTestsGenerator struct{}

func (SuiteBrokenParallelTestsGenerator) Checker() checkers.Checker {
	return checkers.NewSuiteBrokenParallel()
}

func (g SuiteBrokenParallelTestsGenerator) TemplateData() any {
	var (
		checker = g.Checker().Name()
		report  = QuoteReport(checker + ": testify v1 does not support suite's parallel tests and subtests")
	)

	return struct {
		CheckerName CheckerName
		Report      string
	}{
		CheckerName: CheckerName(checker),
		Report:      report,
	}
}

func (SuiteBrokenParallelTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("SuiteBrokenParallelTestsGenerator.ErroredTemplate").
		Parse(suiteBrokenParallelTestTmpl))
}

func (SuiteBrokenParallelTestsGenerator) GoldenTemplate() Executor {
	parallelRe := regexp.MustCompile(`\n.*(s\.T\(\)|t)\.Parallel\(\) // want \{\{ \$\.Report }}`)
	return template.Must(template.New("SuiteBrokenParallelTestsGenerator.GoldenTemplate").
		Parse(parallelRe.ReplaceAllString(suiteBrokenParallelTestTmpl, "")))
}

const suiteBrokenParallelTestTmpl = header + `

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
	t.Parallel()
	suite.Run(t, new({{ $suiteName }}))
}

func (s *SuiteBrokenParallelCheckerSuite) BeforeTest(_, _ string) {
	s.T().Parallel() // want {{ $.Report }}
}

func (s *{{ $suiteName }}) SetupTest() {
	s.T().Parallel() // want {{ $.Report }}
}

func (s *{{ $suiteName }}) TestAll() {
	s.T().Parallel() // want {{ $.Report }}

	s.Run("", func() {
		s.T().Parallel() // want {{ $.Report }}
		s.Equal(1, 2)
	})

	s.T().Run("", func(t *testing.T) {
		t.Parallel() // want {{ $.Report }}
		assert.Equal(t, 1, 2)
	})

	s.Run("", func() {
		s.Equal(1, 2)

		s.Run("", func() {
			s.Equal(1, 2)
			
			s.Run("", func() {
				s.T().Parallel() // want {{ $.Report }}
				s.Equal(1, 2)
			})
		
			s.T().Run("", func(t *testing.T) {
				t.Parallel() // want {{ $.Report }}
				assert.Equal(t, 1, 2)
			})
		})
	})
}

func (s *{{ $suiteName }}) TestTable() {
	cases := []struct{ Name string }{}

	for _, tt := range cases {
		tt := tt
		s.T().Run(tt.Name, func(t *testing.T) {
			t.Parallel() // want {{ $.Report }}
			s.Equal(1, 2)
		})
	}

	for _, tt := range cases {
		tt := tt
		s.Run(tt.Name, func() {
			s.T().Parallel() // want {{ $.Report }}
			s.Equal(1, 2)
		})
	}
}

func TestSimpleTable(t *testing.T) {
	t.Parallel()

	cases := []struct{ Name string }{}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, 1, 2)
		})
	}
}
`
