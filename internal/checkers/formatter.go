package checkers

import (
	"fmt"
	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/ast/inspector"
	"maps"
	"slices"
	"strings"
)

// Formatter detects situations like
// TODO

// TODO: флаг format-string-checks:

// https://cs.opensource.google/go/x/tools/+/master:go/analysis/passes/printf/doc.go
type Formatter struct{}

// NewFormatter constructs Formatter checker.
func NewFormatter() Formatter  { return Formatter{} }
func (Formatter) Name() string { return "formatter" }

func (checker Formatter) Check(pass *analysis.Pass, inspector *inspector.Inspector) ([]analysis.Diagnostic, error) {
	var result []analysis.Diagnostic

	assertFuncsByName := make(map[rangeImpl]*CallMeta) // {"stretchr/testify/assert.Implementsf": "assert.Implementsf"}

	// Stage 1. Collect assert functions used in package (f-like only).
	// Along the way we swear at functions that behave like format ones, but do not end in "f".

	inspector.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node) {
		ce := node.(*ast.CallExpr)
		call := NewCallMeta(pass, ce)
		if call == nil {
			return
		}

		if call.Fn.IsFmt {
			k := rangeImpl{
				pos: call.Pos(),
				end: call.End(),
			}
			assertFuncsByName[k] = call
		} else {
			if isPrintfLikeCall(pass, call, call.Fn.Obj) {
				result = append(result, *newUseFunctionDiagnostic(
					checker.Name(),
					call,
					call.Fn.Name+"f",
					newSuggestedFuncReplacement(call, call.Fn.Name+"f"),
				))
			}
		}
	})
	if len(assertFuncsByName) == 0 {
		return nil, nil
	}

	assertFuncNames := make([]string, 0, len(assertFuncsByName))
	for _, v := range assertFuncsByName {
		assertFuncNames = append(assertFuncNames, v.Fn.Obj.FullName())
	}

	// Stage 2. Configure go vet's printf with collected testify functions and run it.

	type diagKey struct {
		pos, end token.Pos
		msg      string
	}
	defaultCfgDiags := make(map[diagKey]struct{})
	addToCache := func(d analysis.Diagnostic) {
		defaultCfgDiags[diagKey{
			pos: d.Pos,
			end: d.End,
			msg: d.Message,
		}] = struct{}{}
	}
	if err := runPrintfAnalyzer(pass, inspector, "", addToCache); err != nil {
		return nil, fmt.Errorf("govet's printf dry run: %v", err)
	}

	addToResult := func(d analysis.Diagnostic) {
		//fmt.Println("PRINTF", pass.Fset.Position(d.Pos).String())
		//fmt.Println("PRINTF", pass.Fset.Position(d.End).String())
		if _, ok := defaultCfgDiags[diagKey{
			pos: d.Pos,
			end: d.End,
			msg: d.Message,
		}]; !ok {
			result = append(result, *newDiagnostic(
				checker.Name(),
				rangeImpl{d.Pos, d.End},
				shortFuncNames(d.Pos, d.Message, assertFuncsByName),
				nil,
			))
		}
	}
	if err := runPrintfAnalyzer(pass, inspector, strings.Join(assertFuncNames, ","), addToResult); err != nil {
		return nil, fmt.Errorf("govet's printf run for testify funcs: %v", err)
	}

	return result, nil
}

func runPrintfAnalyzer(
	pass *analysis.Pass,
	inspector *inspector.Inspector,
	funcs string,
	handler func(d analysis.Diagnostic),
) error {
	originalReport := pass.Report
	pass.Report = handler
	defer func() { pass.Report = originalReport }()

	anlzr := printf.Analyzer
	if funcs != "" {
		if err := anlzr.Flags.Set("funcs", funcs); err != nil {
			return fmt.Errorf("set `printf.funcs` flag: %v", err)
		}
	}
	// TODO (full copy)
	if _, err := anlzr.Run(&analysis.Pass{
		Analyzer:     anlzr,
		Fset:         &(*pass.Fset),
		Files:        slices.Clone(pass.Files),
		OtherFiles:   slices.Clone(pass.OtherFiles),
		IgnoredFiles: slices.Clone(pass.IgnoredFiles),
		Pkg:          &(*pass.Pkg),
		TypesInfo: &types.Info{
			Types:        maps.Clone(pass.TypesInfo.Types),
			Instances:    maps.Clone(pass.TypesInfo.Instances),
			Defs:         maps.Clone(pass.TypesInfo.Defs),
			Uses:         maps.Clone(pass.TypesInfo.Uses),
			Implicits:    maps.Clone(pass.TypesInfo.Implicits),
			Selections:   maps.Clone(pass.TypesInfo.Selections),
			Scopes:       maps.Clone(pass.TypesInfo.Scopes),
			InitOrder:    pass.TypesInfo.InitOrder,
			FileVersions: maps.Clone(pass.TypesInfo.FileVersions),
		},
		TypesSizes:        pass.TypesSizes,
		TypeErrors:        slices.Clone(pass.TypeErrors),
		Report:            handler,
		ResultOf:          map[*analysis.Analyzer]any{inspect.Analyzer: inspector},
		ImportObjectFact:  func(obj types.Object, fact analysis.Fact) bool { return false },
		ImportPackageFact: func(pkg *types.Package, fact analysis.Fact) bool { return false },
		ExportObjectFact:  func(obj types.Object, fact analysis.Fact) {},
		ExportPackageFact: func(fact analysis.Fact) {},
		AllPackageFacts:   func() []analysis.PackageFact { return slices.Clone(pass.AllPackageFacts()) },
		AllObjectFacts:    func() []analysis.ObjectFact { return slices.Clone(pass.AllObjectFacts()) },
	}); err != nil {
		return fmt.Errorf("analyzer run: %v", err)
	}
	return nil
}

func isPrintfLikeCall(pass *analysis.Pass, call *CallMeta, fn *types.Func) bool {
	msgAndArgsPos := getMsgAndArgsPosition(fn)
	if msgAndArgsPos == 0 {
		return false
	}

	fmtFn := analysisutil.ObjectOf(pass.Pkg, testify.AssertPkgPath, call.Fn.Name+"f")
	if fmtFn == nil {
		// NOTE(a.telyshev): No formatted analogue.
		return false
	}

	return len(call.ArgsRaw) > msgAndArgsPos
}

func shortFuncNames(pos token.Pos, msg string, assertFuncsByName map[rangeImpl]*CallMeta) string {
	for rng, call := range assertFuncsByName {
		if rng.Includes(pos) {
			return strings.ReplaceAll(msg, call.Fn.Obj.FullName(), call.String())
		}
	}
	return msg
}

type rangeImpl struct{ pos, end token.Pos }

func (r rangeImpl) Pos() token.Pos            { return r.pos }
func (r rangeImpl) End() token.Pos            { return r.end }
func (r rangeImpl) Includes(p token.Pos) bool { return r.pos <= p && p <= r.end }

func getMsgAndArgsPosition(fn *types.Func) int {
	signature, ok := fn.Type().(*types.Signature)
	if !ok {
		return 0
	}

	params := signature.Params()
	if params.Len() < 1 {
		return 0
	}
	if last := params.Len() - 1; params.At(last).Name() == "msgAndArgs" {
		return last
	}

	return 0
}

// {pos, end}: {fnName, shortName}
// for printf diags
//		if diag.range in {pos, end}
//			try replace fnName -> shortName
//  chatgpt про эффективный поиск

// comment про segment tree
