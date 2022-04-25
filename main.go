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
// - Empty проверяет только для сравнений len() == 0, не трогая zero value
// - что делать, если функции ещё нет в testify (а линтер её просит)

// TODO:
// - поддержка алиасов
// - тест на переоределённый builtint
// - проверка наличия импортов
// - проверка, что мы сейчас находимся в тестах
// - проверка тестов в соответствии с каждым методом API
// - fuzzy testing?
// - написать генератор тестов
// - проверить что при go get линтера не ставится лишнего
// - описать правила контрибьютинга (чекер, генератор тестов)
// - TODO: кинуть issue во floatcompare о недостающих проверках
// TODO: я осознаю, что местами тесты избыточны. но считаю, что тестов много не бывает :)
// todo: предлагаю решать по задачке в день, чтобы не утомляться и не становилось скучно и лень
/*
как дебагать

	if !strings.HasSuffix(pass.Fset.Position(expr.Pos()).Filename, "float_compare_generated.go") {
		return false
	}
*/
// - сам линтер зависимостей не имеет (или по минимуму, например, нет testify)
