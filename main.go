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

// TODO:
// - поддержка алиасов
// - тест на переоределённый builtint
// - проверка наличия импортов
// - проверка тестов в соответствии с каждым методом API
// - fuzzy testing?
// - написать генератор тестов
