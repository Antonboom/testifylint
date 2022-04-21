package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
)

var generators = map[string]TestsGenerator{
	"pkg/analyzer/testdata/src/basic/empty_generated.go": new(EmptyCasesGenerator),
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
		fmt.Println(b.String())
		return fmt.Errorf("fmt result: %v", err)
	}

	return ioutil.WriteFile(output, formatted, 0644)
}
