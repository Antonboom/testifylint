package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/testify"
)

const (
	DefaultSuiteConsistencyReceiverName = "s"
	upperCamelNameTermPattern           = `[A-Z][a-zA-Z0-9]*`
)

var (
	DefaultSuiteConsistencyTestNamePattern = regexp.MustCompile(
		"^Test" + upperCamelNameTermPattern + "(_" + upperCamelNameTermPattern + ")*$",
	)
	suiteRunnerVariantPattern = regexp.MustCompile(
		"^" + upperCamelNameTermPattern + "(_" + upperCamelNameTermPattern + ")*$",
	)
)

// SuiteConsistency detects inconsistent suite-related names:
//   - top-level suite runner names should be Test<SuiteName> or Test<SuiteName>_<Variant>;
//   - suite receiver names should be consistent;
//   - suite test method names should match the configured pattern.
//
// Indirect suite.Run callbacks are intentionally out of scope. The checker only analyzes direct
// suite.Run calls inside top-level test functions, which matches Go's own test discovery model.
type SuiteConsistency struct {
	receiverName   string
	testNameRegexp *regexp.Regexp
}

// NewSuiteConsistency constructs SuiteConsistency checker.
func NewSuiteConsistency() *SuiteConsistency {
	return &SuiteConsistency{
		receiverName:   DefaultSuiteConsistencyReceiverName,
		testNameRegexp: DefaultSuiteConsistencyTestNamePattern,
	}
}

func (*SuiteConsistency) Name() string { return "suite-consistency" }

func (checker *SuiteConsistency) SetReceiverName(name string) *SuiteConsistency {
	checker.receiverName = name
	return checker
}

func (checker *SuiteConsistency) SetTestNamePattern(pattern *regexp.Regexp) *SuiteConsistency {
	checker.testNameRegexp = pattern
	return checker
}

func (checker *SuiteConsistency) Check(pass *analysis.Pass, insp *inspector.Inspector) (diagnostics []analysis.Diagnostic) {
	suiteRunObj := analysisutil.ObjectOf(pass.Pkg, testify.SuitePkgPath, "Run")

	insp.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(node ast.Node) {
		fd := node.(*ast.FuncDecl)

		if isSuiteMethod(pass, fd) {
			if diag := checker.checkReceiver(pass, fd); diag != nil {
				diagnostics = append(diagnostics, *diag)
			}
			if diag := checker.checkTestMethodName(pass, fd); diag != nil {
				diagnostics = append(diagnostics, *diag)
			}
		}

		if suiteRunObj != nil {
			if diag := checker.checkSuiteRunFunc(pass, fd, suiteRunObj); diag != nil {
				diagnostics = append(diagnostics, *diag)
			}
		}
	})

	return diagnostics
}

func (checker *SuiteConsistency) checkReceiver(pass *analysis.Pass, fd *ast.FuncDecl) *analysis.Diagnostic {
	if checker.receiverName == "" {
		return nil
	}

	rcv := fd.Recv.List[0]
	if len(rcv.Names) != 1 || rcv.Names[0] == nil {
		msg := "suite receiver should be named " + checker.receiverName
		return newDiagnostic(checker.Name(), fd, msg)
	}

	currentName := rcv.Names[0]
	if currentName.Name == checker.receiverName {
		return nil
	}

	msg := fmt.Sprintf("suite receiver name %s does not match configured name %s", currentName.Name, checker.receiverName)
	if identConflictsWithFuncScope(pass, fd, currentName, checker.receiverName) {
		return newDiagnostic(checker.Name(), currentName, msg)
	}

	edits := []analysis.TextEdit{{
		Pos:     currentName.Pos(),
		End:     currentName.End(),
		NewText: []byte(checker.receiverName),
	}}
	if fd.Body != nil {
		ast.Inspect(fd.Body, func(node ast.Node) bool {
			id, ok := node.(*ast.Ident)
			if !ok || pass.TypesInfo.Uses[id] != pass.TypesInfo.Defs[currentName] {
				return true
			}
			edits = append(edits, analysis.TextEdit{
				Pos:     id.Pos(),
				End:     id.End(),
				NewText: []byte(checker.receiverName),
			})
			return true
		})
	}

	return newDiagnostic(checker.Name(), currentName, msg, analysis.SuggestedFix{
		Message:   "Rename receiver to " + checker.receiverName,
		TextEdits: edits,
	})
}

func (checker *SuiteConsistency) checkTestMethodName(pass *analysis.Pass, fd *ast.FuncDecl) *analysis.Diagnostic {
	if fd.Name == nil || !isSuiteTestMethod(fd.Name.Name) || checker.testNameRegexp == nil {
		return nil
	}
	if checker.testNameRegexp.MatchString(fd.Name.Name) {
		return nil
	}

	msg := fmt.Sprintf("suite test method name %s does not match pattern %s", fd.Name.Name, checker.testNameRegexp.String())

	const legacyPrefix = "Test_"
	if !strings.HasPrefix(fd.Name.Name, legacyPrefix) {
		return newDiagnostic(checker.Name(), fd.Name, msg)
	}

	proposedName := "Test" + strings.TrimPrefix(fd.Name.Name, legacyPrefix)
	if !checker.testNameRegexp.MatchString(proposedName) || methodExists(pass, fd, proposedName) {
		return newDiagnostic(checker.Name(), fd.Name, msg)
	}

	return newDiagnostic(checker.Name(), fd.Name, msg, analysis.SuggestedFix{
		Message:   "Rename to " + proposedName,
		TextEdits: renameObjectEdits(pass, pass.TypesInfo.Defs[fd.Name], proposedName),
	})
}

func (checker *SuiteConsistency) checkSuiteRunFunc(pass *analysis.Pass,
	fd *ast.FuncDecl, suiteRunObj types.Object,
) *analysis.Diagnostic {
	if fd.Name == nil || fd.Recv != nil || !strings.HasPrefix(fd.Name.Name, "Test") || !hasStdTestingTParam(pass, fd.Type) {
		return nil
	}

	var suiteName string
	ast.Inspect(fd.Body, func(node ast.Node) bool {
		switch node.(type) {
		case nil:
			return false
		case *ast.FuncLit:
			return false
		}

		ce, ok := node.(*ast.CallExpr)
		if !ok || len(ce.Args) < 2 {
			return true
		}

		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok || !analysisutil.IsObj(pass.TypesInfo, se.Sel, suiteRunObj) {
			return true
		}

		suiteName = suiteTypeName(pass.TypesInfo.TypeOf(ce.Args[1]))
		return suiteName == ""
	})
	if suiteName == "" {
		return nil
	}

	expectedName := "Test" + suiteName
	if isValidSuiteRunnerName(fd.Name.Name, expectedName) {
		return nil
	}

	msg := fmt.Sprintf("suite test function name %s does not match suite name (expected %s or %s_<Variant>)",
		fd.Name.Name, expectedName, expectedName)
	if strings.HasPrefix(fd.Name.Name, expectedName+"_") || packageScopeNameExists(pass, expectedName) {
		return newDiagnostic(checker.Name(), fd.Name, msg)
	}

	return newDiagnostic(checker.Name(), fd.Name, msg, analysis.SuggestedFix{
		Message:   "Rename to " + expectedName,
		TextEdits: renameObjectEdits(pass, pass.TypesInfo.Defs[fd.Name], expectedName),
	})
}

func isValidSuiteRunnerName(name, expectedName string) bool {
	if name == expectedName {
		return true
	}
	suffix, ok := strings.CutPrefix(name, expectedName+"_")
	if !ok {
		return false
	}

	return suiteRunnerVariantPattern.MatchString(suffix)
}

func suiteTypeName(typ types.Type) string {
	if typ == nil {
		return ""
	}
	typ = types.Unalias(typ)
	for {
		ptr, ok := typ.(*types.Pointer)
		if !ok {
			break
		}
		typ = types.Unalias(ptr.Elem())
	}

	named, ok := typ.(*types.Named)
	if !ok {
		return ""
	}
	return named.Obj().Name()
}

func identConflictsWithFuncScope(pass *analysis.Pass, fd *ast.FuncDecl, original *ast.Ident, proposedName string) bool {
	originalObj := pass.TypesInfo.Defs[original]
	if originalObj == nil {
		return true
	}

	conflict := false
	ast.Inspect(fd, func(node ast.Node) bool {
		if conflict {
			return false
		}

		id, ok := node.(*ast.Ident)
		if !ok || id.Name != proposedName {
			return true
		}

		if obj := pass.TypesInfo.Defs[id]; obj != nil && obj != originalObj {
			conflict = true
			return false
		}
		return true
	})
	return conflict
}

func methodExists(pass *analysis.Pass, fd *ast.FuncDecl, methodName string) bool {
	rcvType := pass.TypesInfo.TypeOf(fd.Recv.List[0].Type)
	if rcvType == nil {
		return true
	}

	obj, _, _ := types.LookupFieldOrMethod(rcvType, true, pass.Pkg, methodName)
	return obj != nil
}

func packageScopeNameExists(pass *analysis.Pass, name string) bool {
	return pass.Pkg.Scope().Lookup(name) != nil
}

func renameObjectEdits(pass *analysis.Pass, obj types.Object, proposedName string) []analysis.TextEdit {
	if obj == nil {
		return nil
	}

	edits := []analysis.TextEdit{{
		Pos:     obj.Pos(),
		End:     obj.Pos() + token.Pos(len(obj.Name())),
		NewText: []byte(proposedName),
	}}
	for id, useObj := range pass.TypesInfo.Uses {
		if useObj != obj {
			continue
		}
		edits = append(edits, analysis.TextEdit{
			Pos:     id.Pos(),
			End:     id.End(),
			NewText: []byte(proposedName),
		})
	}
	return edits
}

func IsValidSuiteConsistencyReceiverName(name string) bool {
	return token.IsIdentifier(name) && name != "_"
}
