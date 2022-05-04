package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"text/template"
)

var generators = map[string]TestsGenerator{
	//"pkg/analyzer/testdata/src/basic/bool_compare_generated.go":    BoolCompareCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/comparisons_generated.go":     ComparisonsCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/empty_generated.go":           EmptyCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/error_generated.go":           ErrorCasesGenerator{},
	"pkg/analyzer/testdata/src/basic/error_is.go": ErrorIsCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/expected_actual_generated.go": ExpectedActualCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/float_compare_generated.go":   FloatCompareCasesGenerator{},
	//"pkg/analyzer/testdata/src/basic/len_generated.go":             LenCasesGenerator{},
}

func main() {
	for output, g := range generators {
		if err := genGoFileFromTmpl(output, g.ErroredTemplate(), g.Data()); err != nil {
			log.Panic(err)
		}

		if err := genGoFileFromTmpl(output+".golden", g.GoldenTemplate(), g.Data()); err != nil {
			log.Panic(err)
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
		return fmt.Errorf("fmt result: %v", err)
	}

	return ioutil.WriteFile(output, formatted, 0644)
}
