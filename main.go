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
// - не поддерживаем кастомные типы, удовлетворяющие TestingSuite
// - not support exotic like eq := assert.Equal eq(1, 2)

// Open for contribution:
// - Empty with zerovalue
// - Zero
// - Negative
// - HTTPCodeConstant + HTTPSuccess + HTTPError
// - нейминг для теста, запускающего suite, а также его местоположение
// - float compare for structs with float comparisons
// - старайтесь, чтобы тестовые файлы не превышали 3к строк

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

	if !strings.HasSuffix(pass.Fset.Position(expr.Range()).Filename, "float_compare_generated.go") {
		return false
	}
*/
// - сам линтер зависимостей не имеет (или по минимуму, например, нет testify)
// кинуть issue Бакину о пресете test
// TODO: https://github.com/ghetzel/hydra/blob/master/gen_test.go
// TODO: checker msg in constant?
// TODO: issue, что floatcompare можно побороть с помощью generics
// подебажить, какие Range лучше
// suggested fixes: https://github.com/golang/tools/blob/master/go/analysis/doc/suggested_fixes.md
// suite checkers only if suite imported

// почему validateCheck не check.Validate? потому что чек не знает, что с ним будут делать

// Как создавался этот линтер (ссылка на курс в ридми)

// readme: имя чекера, пример, предлагаемый фикс, автоматический ли он
// todo: бакину рассказать про неуникальность token.Pos

// todo save idea runs

// todo отказаться от ExprString

// todo нашёл issue, пока реализовывал checker len

// todo: все тесты на XOR вынести в отдельный testdata/src/strange, обратить пристальное внимание на ошибки в выражениях

// тред о паниках: https://github.com/orgs/golangci/teams/team/discussions/33?from_comment=5
// todo: https://habr.com/ru/company/joom/blog/666440/

// todo: дока к каждому чекеру

/*
❌	require.Nil(t, err)
✅	require.NoError(t, err)

❌	assert.Equal(t, 300.0, float64(price.Amount))
✅	assert.EqualValues(t, 300.0, price.Amount)

❌	assert.Equal(t, 0, len(result.Errors))
✅	assert.Empty(t, result.Errors)

❌	require.Equal(t, len(expected), len(result)
	sort.Slice(expected, ...)
	sort.Slice(result, ...)
	for i := range result {
		assert.Equal(t, expected[i], result[i])
	}
✅	assert.ElementsMatch(t, expected, result)




Также стоит быть осторожнее при использовании горутин в тестах. require-проверки производятся через runtime.goexit(),
так что они сработают ожидаемым образом только в основной горутине.
https://github.com/golang/go/issues/20940

не использовать equalError и ErrorContains
(покрывается линтером forbidigo)


таблика:
линтер, пример, имеет ли автофикс
*/
