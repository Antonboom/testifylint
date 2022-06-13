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
		mustNil(config.Validate(cfg))
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
- ругаться, если ни одного чекера не включено
- expected-actual config test
- бейдж доки, в доке чекеров пример ассерта и текста
- TODO: https://github.com/ghetzel/hydra/blob/master/gen_test.go -> в тест
- todo: все тесты на XOR вынести в отдельный testdata/src/strange, обратить пристальное внимание на ошибки в выражениях
- подебажить, какие Range лучше
тест, интересует ginkgo.T()
https://github.com/kubernetes/ingress-nginx/blob/main/test/e2e/loadbalance/ewma.go
- ревью приоритетов чекеров, проверять при сборке, что приоритеты разные
*/
