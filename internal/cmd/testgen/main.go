package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
)

var generators = map[string]TestsGenerator{
	"pkg/analyzer/testdata/src/basic/bool_compare_generated.go":  BoolCompareCasesGenerator{},
	"pkg/analyzer/testdata/src/basic/empty_generated.go":         EmptyCasesGenerator{},
	"pkg/analyzer/testdata/src/basic/float_compare_generated.go": FloatCompareCasesGenerator{},
}

func main() {
	for output, gen := range generators {
		if err := genTests(output, gen); err != nil {
			log.Panic(err)
		}
	}
}

func genTests(output string, g TestsGenerator) error {
	b := bytes.NewBuffer(nil)
	if err := g.Template().Execute(b, g.Data()); err != nil {
		return fmt.Errorf("execute gen tmpl: %v", err)
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		_ = ioutil.WriteFile(output, b.Bytes(), 0644) // For debug.
		return fmt.Errorf("fmt result: %v", err)
	}

	return ioutil.WriteFile(output, formatted, 0644)
}
