package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

var generators = map[string]TestsGenerator{
	"pkg/analyzer/testdata/src/basic/bool_compare_test.go": BoolCompareCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/comparisons_generated.go":     ComparisonsCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/empty_test.go": EmptyCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/error_generated.go":           ErrorCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/error_is_test.go":        ErrorIsCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/expected_actual_test.go": ExpectedActualCasesGenerator{},
	"pkg/analyzer/testdata/src/basic/float_compare_test.go": FloatCompareCasesGenerator{},
	"pkg/analyzer/testdata/src/basic/len_test.go":           LenCasesGenerator{},
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
		_ = ioutil.WriteFile(output, b.Bytes(), 0644) // For debug.
		return fmt.Errorf("format %s: %v", output, err)
	}

	return ioutil.WriteFile(output, formatted, 0644)
}
