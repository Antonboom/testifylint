package main

import (
	"flag"
	"log"
	"os"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Antonboom/testifylint/analyzer"
	"github.com/Antonboom/testifylint/config"
)

var (
	configPath = flag.String("config", ".testifylint.yml", "path to config file (yml)")
	dumpCfg    = flag.Bool("dump-config", false, "dump default config (yml) in stdout")
)

func main() {
	flag.Parse()

	if *dumpCfg {
		mustNil(config.DumpDefault(os.Stdout))
		return
	}

	var cfg config.Config
	if *configPath != "" {
		if _, err := os.Stat(*configPath); err != nil {
			cfg = config.Default
		} else {
			var err error
			cfg, err = config.ParseFromFile(*configPath)
			mustNil(err)
			mustNil(config.Validate(cfg))
		}
	}

	singlechecker.Main(analyzer.New(cfg))
}

func mustNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
- Flags Issue: https://github.com/golang/go/issues/53336
- expected-actual config test
- бейдж доки, в доке чекеров пример ассерта и текста
- TODO: https://github.com/ghetzel/hydra/blob/master/gen_test.go -> в тест
- todo: все тесты на XOR вынести в отдельный testdata/src/strange, обратить пристальное внимание на ошибки в выражениях
- подебажить, какие Range лучше
тест, интересует ginkgo.T()
https://github.com/kubernetes/ingress-nginx/blob/main/test/e2e/loadbalance/ewma.go
- ревью приоритетов чекеров, проверять при сборке, что приоритеты разные
- продублировать полезные флаги
- godoc к чекерам
- пристальное ревью каждого файла
- negative test cases
- поддержка pkg alias
- обновить vendor в testdata
// TODO(a.telyshev): s.T().Run( -> s.Run
// TODO: имя теста не повторяет имя сьюта TestService_Run() {}
// TODO: имя тестового метода (похоже на линтер для имён теста в cases)
-переименовать internal/checker в internal/checkers
- вынести генератор куда-нибудь?
- пересекающиеся тесты

func Get(name string) (Checker, bool) {
	ch, ok := checkersByName[name]
	return ch, ok
}
описать багу – Checker встраивался, но дальше не конвертировался в CallChecker и AdvancedChecker
от этого спасёт маркер ну или прост использовать переменную
покрыть тестом
- заполнить URL'ы диагностик ссылками на ридми/checker
- review приоритета
readme – tests is code too

- облегчить тесты, унести различные вариации на тесты фактов
- deduplication репортов скрывает баги опхода –> хуже optimization
- parallel tests? speed optimization
- снять профили
- уменьшить количество тестов?
- негативные ветки в анализаторе (импорты, файлы и тд)

- посмотреть на тесты и что мы тестируем, нельзя ли вынести общее?
- финальное ревью каждого файла и тестового файла

suite.T() vs suite.Run
- no-f-assertions no-fmt-mess
	проверить, как go vet ведёт себя на примерах
	как вариант – с форматированием только f, без форматирования только Equal
	https://github.com/stretchr/testify/issues/339
	https://github.com/stretchr/testify/issues/471
	https://go.googlesource.com/tools/+/refs/heads/release-branch.go1.12/go/analysis/passes/printf/printf.go?pli=1#

CheckerExpander -> AssertionExpander

применение ChatGPT

- про testify v1 и ожидание v2
 https://github.com/stretchr/testify/issues/1089
 https://github.com/stretchr/testify/milestone/4
или он не произойдёт?
https://github.com/stretchr/testify/pull/1109#issuecomment-1650619745

todo -> readme старайтесь поддерживать тесты маленькими, выделяя общий код
100-200 строк достаточно обычно (взять среднее по тому, что получилось)

e2e тесты на cli (config issues)?

в README к каждому чекеру добавить причину

грепнуть FIXME TODO без маски *.go


- уйти от `types.ExprString`, чекнуть, что forbidigo работает

https://floating-point-gui.de/
https://floating-point-gui.de/errors/comparison/
https://www.exploringbinary.com/why-0-point-1-does-not-exist-in-floating-point/

require-len – по аналогии с require-error
error-compare - запрет на EqualError, ErrorContains

suggested edit + golden file + more granular message for require-error

вынести contribution в отдельный md

не заменить ли конфиг полностью на флаги? или и то и то? погрумить.

readme – таблица с череками и включен/выключен по умолчанию

TODO: -race significantly decrease test speed
*/
