package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

var generators = map[string]TestsGenerator{
	toCheckersPath("most-of", "bool_compare_test.go"):        BoolCompareCasesGenerator{},
	toCheckersPath("compares", "compares_test.go"):           ComparesCasesGenerator{},
	toCheckersPath("most-of", "empty_test.go"):               EmptyCasesGenerator{},
	toCheckersPath("most-of", "error_test.go"):               ErrorCasesGenerator{},
	toCheckersPath("most-of", "error_is_test.go"):            ErrorIsCasesGenerator{},
	toCheckersPath("most-of", "expected_actual_test.go"):     ExpectedActualCasesGenerator{},
	toCheckersPath("most-of", "float_compare_test.go"):       FloatCompareCasesGenerator{},
	toCheckersPath("most-of", "len_test.go"):                 LenCasesGenerator{},
	toCheckersPath("require-error", "require_error_test.go"): RequireErrorCasesGenerator{},
	toCheckersPath("suite-no-extra-assert-call",
		"suite_no_extra_assert_call_test.go"): SuiteNoExtraAssertCallCasesGenerator{},
}

func toCheckersPath(dirsFile ...string) string {
	return filepath.Join(append([]string{"pkg", "analyzer", "testdata", "src", "checkers"}, dirsFile...)...)
}

func main() {
	for output, g := range generators {
		if !strings.HasSuffix(output, "_test.go") {
			panic(output + " is not test file!")
		}

		if err := genGoFileFromTmpl(output, g.ErroredTemplate(), g.Data()); err != nil {
			log.Panic(err)
		}

		if goldenTmpl := g.GoldenTemplate(); goldenTmpl != nil {
			if err := genGoFileFromTmpl(output+".golden", goldenTmpl, g.Data()); err != nil {
				log.Panic(err)
			}
		}
	}
}

func genGoFileFromTmpl(output string, tmpl *template.Template, data any) error {
	b := bytes.NewBuffer(nil)
	if err := tmpl.Execute(b, data); err != nil {
		return fmt.Errorf("execute cases tmpl: %v", err)
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		_ = ioutil.WriteFile(output, b.Bytes(), 0o644) // For debug.
		return fmt.Errorf("format %s: %v", output, err)
	}

	return ioutil.WriteFile(output, formatted, 0o644)
}
