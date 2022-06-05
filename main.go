package main

import (
	"flag"
	"log"
	"os"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Antonboom/testifylint/pkg/analyzer"
	"github.com/Antonboom/testifylint/pkg/config"
)

var (
	configPath = flag.String("config", "", "path to config file (yml)")
	dumpCfg    = flag.Bool("dump-config", false, "dump config example (yml) in stdout")
)

func main() {
	flag.Parse()

	if *dumpCfg {
		mustNil(config.Dump(config.Default, os.Stdout))
		return
	}

	var cfg config.Config
	if *configPath != "" {
		var err error
		cfg, err = config.ParseFromFile(*configPath)
		mustNil(err)
	}

	singlechecker.Main(analyzer.New(cfg))
}

func mustNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Assumptions:
// - не работает, если алиас для функции сделали
// - Empty проверяет только для сравнений len() == 0, не трогая zero value
// - что делать, если функции ещё нет в testify (а линтер её просит)
// - не поддерживаем кастомные типы, удовлетворяющие TestingSuite
// - not support exotic like eq := assert.Equal eq(1, 2)

// Open for contribution:
// - zero (чекает zero value)
// - negative-positive
// - http-error (HTTPSuccess + HTTPError)
// - http-code-const
// - float compare for structs with float comparisons
// - suite-test-name
// - equal-values
// - no-error-contains
// - no-equal-error
// - elements-match
//
// - старайтесь, чтобы тестовые файлы не превышали 3к строк
// - реализуйте интерфейс Disabled для выключения
// - сначала пишите генератор тестов
// (я осознаю, что местами тесты избыточны. но считаю, что тестов много не бывает)
// - добавьте тест анализатора
// - потом реализуйте чекер и укажите его в списке

// TODO:
// - п
// - TODO: кинуть issue во floatcompare о недостающих проверках
// TODO:
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
// дампить конфиг
// вывод приоритета линтеров, какие включены, какие выключены

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
линтер, пример, имеет ли автофикс, enabled by default



роверить что при go get линтера не ставится лишнего
https://stackoverflow.com/questions/64071364/best-way-to-use-test-dependencies-in-go-but-prevent-export-them

https://grep.app/search?q=make.%2A&regexp=true&filter[lang][0]=Go
https://sourcegraph.com/search?q=context:global+t+testing.TB+count:1000000&patternType=literal

*/
