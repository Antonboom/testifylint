package main

import (
	"text/template"

	"github.com/Antonboom/testifylint/internal/checkers"
)

type GoRequireTestsGenerator struct{}

func (GoRequireTestsGenerator) Checker() checkers.Checker {
	return checkers.NewGoRequire()
}

func (g GoRequireTestsGenerator) TemplateData() any {
	var (
		name                = g.Checker().Name()
		requireReport       = name + ": require must only be used in the goroutine running the test function%.s%.s"
		assertFailNowReport = name + ": %s.%s must only be used in the goroutine running the test function"
		fnReport            = name + ": %s contains assertions that must only be used in the goroutine running the test function"
	)

	return struct {
		CheckerName CheckerName
		FnReport    string
		Assertions  []Assertion
		Requires    []Assertion
	}{
		CheckerName: CheckerName(name),
		FnReport:    fnReport,
		Assertions: []Assertion{
			{Fn: "Fail", Argsf: `"boom!"`},
			{Fn: "FailNow", Argsf: `"boom!"`, ReportMsgf: assertFailNowReport, ProposedFn: "FailNow"},
			{Fn: "NoError", Argsf: "err"},
			{Fn: "True", Argsf: "b"},
		},
		Requires: []Assertion{
			{Fn: "Fail", Argsf: `"boom!"`, ReportMsgf: requireReport, ProposedFn: "Fail"},
			{Fn: "FailNow", Argsf: `"boom!"`, ReportMsgf: requireReport, ProposedFn: "FailNow"},
			{Fn: "NoError", Argsf: "err", ReportMsgf: requireReport, ProposedFn: "NoError"},
			{Fn: "True", Argsf: "b", ReportMsgf: requireReport, ProposedFn: "True"},
		},
	}
}

func (GoRequireTestsGenerator) ErroredTemplate() Executor {
	return template.Must(template.New("GoRequireTestsGenerator.ErroredTemplate").
		Funcs(fm).
		Parse(goRequireTestTmpl))
}

func (GoRequireTestsGenerator) GoldenTemplate() Executor {
	// NOTE(a.telyshev): Usually this warning leads to full refactoring of test architecture.
	return nil
}

const goRequireTestTmpl = header + `

package {{ .CheckerName.AsPkgName }}

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

{{ define "assertions" }}
	run()
	assertSomething(t)
	requireSomething(t) // want "{{ printf .FnReport "requireSomething" }}"

	{{ range $ai, $assrn := $.Assertions }}
		{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
	{{- end }}
	{{ range $ai, $assrn := $.Requires }}
		{{ NewAssertionExpander.Expand $assrn "require" "t" nil }}
	{{- end }}
{{- end }}

{{ define "silent-assertions" }}
	run()
	assertSomething(t)
	requireSomething(t)
	
	{{ range $ai, $assrn := $.Assertions }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "assert" "t" nil }}
	{{- end }}
	{{ range $ai, $assrn := $.Requires }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "require" "t" nil }}
	{{- end }}
{{- end }}

{{ define "assertions-short" }}
	run()
	assertSomething(t)
	requireSomething(t) // want "{{ printf .FnReport "requireSomething" }}"
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Assertions 0) "assert" "t" nil }}
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Requires 0) "require" "t" nil }}
{{- end }}

{{ define "silent-assertions-short" }}
	run()
	assertSomething(t)
	requireSomething(t)
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Assertions 0).WithoutReport "assert" "t" nil }}
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Requires 0).WithoutReport "require" "t" nil }}
{{- end }}

func {{ .CheckerName.AsTestName }}_Smoke(t *testing.T) {
	var err error
	var b bool

	{{ template "silent-assertions" . }}

	var wg sync.WaitGroup
	defer wg.Wait()

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			{{ template "assertions" . }}
			
			if assert.Error(t, err) {
				{{- range $ai, $assrn := $.Assertions }}
					{{- if eq $assrn.Fn "FailNow"}}
						{{ NewAssertionExpander.Expand $assrn "assert" "t" nil }}
					{{- end }}
				{{- end }}
			}
		}(i)
	}
}

func {{ .CheckerName.AsTestName }}(t *testing.T) {
	defer func() {
		{{- template "silent-assertions-short" . }}

		func() {
			{{- template "silent-assertions-short" . }}
		}()
	}()
	
	defer func() {
		go func() {
			{{ template "assertions-short" . }}
		}()
	}()

	defer run()
	defer assertSomething(t)
	defer requireSomething(t)

	t.Cleanup(func() {
		{{- template "silent-assertions-short" . }}
		
		func() {
			{{- template "silent-assertions-short" . }}
		}()
	
		go func() {
			{{- template "assertions-short" . }}

			func() {
				func(t *testing.T) {
					{{- template "assertions-short" . }}
				}(t)
			}()
		}()
	})

	func() {
		func() {
			func() {
				func() {
					{{- template "silent-assertions-short" . }}
				}()
				
				func(t *testing.T) {
					{{- template "silent-assertions-short" . }}

					t.Run("", func(t *testing.T) {
						{{- template "silent-assertions-short" . }}

						go func() {
							{{- template "assertions-short" . }}
						}()
					})
				}(t)

				go func() {
					{{- template "assertions-short" . }}
				}()
			}()
		}()
	}()

	if false {
		{{- template "silent-assertions-short" . }}
	}

	go func() {
		{{- template "assertions-short" . }}

		go func() {
			{{- template "assertions-short" . }}

			go func() {
				{{- template "assertions-short" . }}
				
				t.Run("", func(t *testing.T) {
					{{- template "silent-assertions-short" . }}
				})
				{{ template "assertions-short" . }}
			}()
			{{ template "assertions-short" . }}
		}()
		{{ template "assertions-short" . }}
	}()

	t.Run("", func(t *testing.T) {
		{{- template "silent-assertions-short" . }}

		t.Run("", func(t *testing.T) {
			{{- template "silent-assertions-short" . }}

			t.Run("", func(t *testing.T) {
				{{- template "silent-assertions-short" . }}
			
				go func() {
					{{- template "assertions-short" . }}
					
					go func () {
						go genericHelper[*testing.T](t) // want {{ QuoteReport (printf .FnReport "genericHelper[*testing.T]") }}
						go superGenericHelper[*testing.T, int](t) // want {{ QuoteReport (printf .FnReport "superGenericHelper[*testing.T, int]") }}
					}()
				}()
				{{ template "silent-assertions-short" . }}
			})
			{{ template "silent-assertions-short" . }}
		})
		{{ template "silent-assertions-short" . }}
	})

	go func() {
		{{- template "assertions-short" . }}

		t.Run("", func(t *testing.T) {
			{{- template "silent-assertions-short" . }}

			go func(t *testing.T) {
				{{- template "assertions-short" . }}

				t.Run("", func(t *testing.T) {
					{{- template "silent-assertions-short" . }}
				})
				
				if true {
					{{- template "assertions-short" . }}
				}
			}(t)
			{{ template "silent-assertions-short" . }}
		})
		{{ template "assertions-short" . }}
	}()

	cases := []struct{}{}
	for _, tt := range cases {
		tt := tt
		t.Run("", func(t *testing.T) {
			{{- template "silent-assertions-short" . }}

			go func() {
				{{- template "assertions-short" . }}
				_ = tt
			}()
		})
	}

	run()
	assertSomething(t)
	requireSomething(t)
	
	go run()
	go assertSomething(t)
	go requireSomething(t) // want "{{ printf .FnReport "requireSomething" }}"

	go func() {
		run()
	}()
	go func() {
		assertSomething(t)
	}()
	go func() {
		requireSomething(t) // want "{{ printf .FnReport "requireSomething" }}"
	}()

	var err error
	var b bool
	{{ range $ai, $assrn := $.Assertions }}
		go {{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "assert" "t" nil }}
	{{- end }}
	{{ range $ai, $assrn := $.Requires }}
		go {{ NewAssertionExpander.NotFmtSingleMode.Expand $assrn "require" "t" nil }}
	{{- end }}

	go requireSomethingInGo(t)
	go func() {
		requireSomethingInGo(t)
	}()

	go proxy(t) // want "{{ printf .FnReport "proxy" }}"
	go func() {
		proxy(t) // want "{{ printf .FnReport "proxy" }}"
	}()

	t.Run("", assertSomething)
	t.Run("", requireSomething)

	go t.Run("", assertSomething)
	go t.Run("", requireSomething)

	var wg sync.WaitGroup
	wg.Add(2)
	for i := 1; i <= 2; i++ {
		i := i
		go t.Run(fmt.Sprintf("uncaught-%d", i), func(t *testing.T) {
			defer wg.Done()

			{{- template "silent-assertions-short" . }}
		})
	}
	wg.Wait()
	
	genericHelper(t)
	genericHelper[*testing.T](t)
	superGenericHelper[*testing.T, int](t)

	go func() {
		genericHelper(t) // want "{{ printf .FnReport "genericHelper" }}"
		genericHelper[*testing.T](t) // want {{ QuoteReport (printf .FnReport "genericHelper[*testing.T]") }}
		superGenericHelper[*testing.T, int](t) // want {{ QuoteReport (printf .FnReport "superGenericHelper[*testing.T, int]") }}
	}()
}

{{ define "suite-assertions-short" }}
	run()
	assertSomething(s.T())
	requireSomething(s.T()) // want "{{ printf .FnReport "requireSomething" }}"
	s.suiteHelper() // want "{{ printf .FnReport "s.suiteHelper" }}"
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Assertions 0) "s" "" nil }}
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Requires 0) "s.Require()" "" nil }}
{{- end }}

{{ define "suite-silent-assertions-short" }}
	run()
	assertSomething(s.T())
	requireSomething(s.T())
	s.suiteHelper()
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Assertions 0).WithoutReport "s" "" nil }}
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Requires 0).WithoutReport "s.Require()" "" nil }}
{{- end }}

{{ $suiteName := .CheckerName.AsSuiteName }}

type {{ $suiteName }} struct {
	suite.Suite
}

func Test{{ $suiteName }}(t *testing.T) {
	suite.Run(t, new({{ $suiteName }}))
}

func (s *{{ $suiteName }}) TestAll() {
	{{- template "suite-silent-assertions-short" . }}

	go func() {
		{{- template "suite-assertions-short" . }}

		s.Run("", func() {
			{{- template "suite-silent-assertions-short" . }}

			go func(t *testing.T) {
				{{- template "assertions-short" . }}

				s.T().Run("", func(t *testing.T) {
					{{- template "silent-assertions-short" . }}
				})
				
				if true {
					{{- template "assertions-short" . }}
				}
			}(s.T())
			{{ template "suite-silent-assertions-short" . }}
		})
		{{ template "suite-assertions-short" . }}
	}()
	
	{{ template "suite-silent-assertions-short" . }}

	go run()
	go assertSomething(s.T())
	go requireSomething(s.T()) // want "{{ printf .FnReport "requireSomething" }}"
	go s.suiteHelper() // want "{{ printf .FnReport "s.suiteHelper" }}"

	s.Run("", s.suiteHelper)
	s.T().Run("", func(t *testing.T) { requireSomething (t) })

	go s.T().Run("", requireSomething)
	go s.Run("", s.suiteHelper)
}

func (s *{{ $suiteName }}) TestAsertFailNow() {
	var err error
	var b bool

	{{ range $ai, $assrn := $.Assertions }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "s" "" nil }}
		{{ NewAssertionExpander.Expand $assrn.WithoutReport "s.Assert()" "" nil }}
	{{- end }}

	go func() {
		{{- range $ai, $assrn := $.Assertions }}
			{{ NewAssertionExpander.Expand $assrn "s" "" nil }}
			{{ NewAssertionExpander.Expand $assrn "s.Assert()" "" nil }}
		{{- end }}
	}()
}

func (s *{{ $suiteName }}) suiteHelper() {
	s.T().Helper()

	{{ template "suite-silent-assertions-short" . }}

	go func() {
		{{- template "suite-assertions-short" . }}
	}()

	{{ template "suite-silent-assertions-short" . }}
}

func run() {}

func assertSomething(t *testing.T) {
	t.Helper()

	assert.NoError(t, nil)
	assert.Error(t, nil)
	assert.True(t, false)
}

func proxy(t *testing.T) {
	requireSomething(t)
}

func requireSomething(t *testing.T) {
	t.Helper()

	{{ template "silent-assertions-short" . }}
}

func requireSomethingInGo(t *testing.T) {
	go func() {
		{{- template "assertions-short" . }}
	}()
}

func helperNotUsedInGoroutine(t *testing.T) {
	t.Helper()

	{{ template "silent-assertions-short" . }}
	
	go func() {
		{{- template "assertions-short" . }}
	}()
}

type testingT interface {
	require.TestingT
}

func genericHelper[T testingT](t T) {
	run()

	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Assertions 0).WithoutReport "assert" "t" nil }}
	{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Requires 0).WithoutReport "require" "t" nil }}

	go func() {
		{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Assertions 0) "assert" "t" nil }}
		{{ NewAssertionExpander.NotFmtSingleMode.Expand (index $.Requires 0) "require" "t" nil }}
	}()
}

func superGenericHelper[T testingT, T2 any](t T) T2 {
	require.Fail(t, "boom!")
	var zero T2
	return zero
}
`
