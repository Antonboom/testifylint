package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Antonboom/testifylint/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
