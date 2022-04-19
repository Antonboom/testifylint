package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Antonboom/testifylint/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}

// Assumptions:
// - не работает, если алиас для функции сделали
