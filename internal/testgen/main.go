package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Antonboom/testifylint/internal/checkers"
)

var analyserTestdataPath = filepath.Join("analyzer", "testdata", "src")

var freeTestsGenerators = map[string]TestsGenerator{ // by subdirectory.
	"base-test":                       BaseTestsGenerator{},
	"error-as-target":                 ErrorAsTargetTestsGenerator{},
	"expected-var-custom-pattern":     ExpectedVarCustomPatternTestsGenerator{},
	"formatter-not-defaults":          FormatterNotDefaultsTestsGenerator{},
	"require-error-fn-pattern":        RequireErrorFnPatternTestsGenerator{},
	"suite-require-extra-assert-call": SuiteRequireExtraAssertCallTestsGenerator{},
}

var checkerTestsGenerators = []CheckerTestsGenerator{
	BlankImportTestsGenerator{},
	BoolCompareTestsGenerator{},
	ComparesTestsGenerator{},
	EmptyTestsGenerator{},
	ErrorNilTestsGenerator{},
	ErrorIsAsTestsGenerator{},
	ExpectedActualTestsGenerator{},
	FloatCompareTestsGenerator{},
	FormatterTestsGenerator{},
	GoRequireTestsGenerator{},
	LenTestsGenerator{},
	NegativePositiveTestsGenerator{},
	NilCompareTestsGenerator{},
	RequireErrorTestsGenerator{},
	SuiteBrokenParallelTestsGenerator{},
	SuiteDontUsePkgTestsGenerator{},
	SuiteExtraAssertCallTestsGenerator{},
	SuiteSubtestRunTestsGenerator{},
	SuiteTHelperTestsGenerator{},
	UselessAssertTestsGenerator{},
}

func init() {
	genForChecker := make(map[string]struct{}, len(checkerTestsGenerators))
	for _, g := range checkerTestsGenerators {
		name := g.Checker().Name()
		if _, ok := genForChecker[name]; ok {
			panic(fmt.Sprintf("multiple test generators for checker %q", name))
		}
		genForChecker[name] = struct{}{}
	}

	for _, ch := range checkers.All() {
		if _, ok := genForChecker[ch]; !ok {
			log.Printf("[WARN] No generated tests for %q checker\n", ch)
		}
	}
}

func main() {
	if err := generateFreeTests(); err != nil {
		log.Panic(err)
	}

	if err := generateCheckersDefaultTests(); err != nil {
		log.Panic(err)
	}
}

func generateFreeTests() error {
	for dir, g := range freeTestsGenerators {
		testFile := strings.ReplaceAll(dir, "-", "_") + "_test.go"
		output := filepath.Join(analyserTestdataPath, dir, testFile)
		if err := genTestFilesPair(g, output); err != nil {
			return fmt.Errorf("%s: %v", dir, err)
		}
	}
	return nil
}

func generateCheckersDefaultTests() error {
	for _, g := range checkerTestsGenerators {
		checker := g.Checker().Name()
		testFile := strings.ReplaceAll(checker, "-", "_") + "_test.go"
		output := filepath.Join(analyserTestdataPath, "checkers-default", checker, testFile)
		if err := genTestFilesPair(g, output); err != nil {
			return fmt.Errorf("%s: %v", checker, err)
		}
	}
	return nil
}

func genTestFilesPair(g TestsGenerator, path string) error {
	dir := filepath.Dir(path)
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("rm tests dir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir tests dir: %v", err)
	}

	tmplData := g.TemplateData()

	if err := genGoFileFromTmpl(path, g.ErroredTemplate(), tmplData); err != nil {
		return fmt.Errorf("generate test file: %v", err)
	}

	if goldenTmpl := g.GoldenTemplate(); goldenTmpl != nil {
		if err := genGoFileFromTmpl(path+".golden", goldenTmpl, tmplData); err != nil {
			return fmt.Errorf("generate golden file: %v", err)
		}
	} else {
		log.Printf("[WARN] No golden file in %q\n", dir)
	}
	return nil
}

func genGoFileFromTmpl(output string, exec Executor, data any) error {
	b := bytes.NewBuffer(nil)
	if err := exec.Execute(b, data); err != nil {
		return fmt.Errorf("execute cases tmpl: %v", err)
	}

	formatted, err := format.Source(b.Bytes())
	if err != nil {
		_ = os.WriteFile(output, b.Bytes(), 0o644) // For debug.
		return fmt.Errorf("format %s: %v", output, err)
	}

	return os.WriteFile(output, formatted, 0o644)
}
