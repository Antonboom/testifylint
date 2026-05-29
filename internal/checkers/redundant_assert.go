package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// RedundantAssert detects redundant assert.Error calls when a stronger error assertion
// on the same error value is present in the same block.
type RedundantAssert struct{}

// NewRedundantAssert constructs RedundantAssert checker.
func NewRedundantAssert() RedundantAssert { return RedundantAssert{} }
func (RedundantAssert) Name() string      { return "redundant-assert" }

func (checker RedundantAssert) Check(pass *analysis.Pass, insp *inspector.Inspector) []analysis.Diagnostic {
	callsByBlock := make(map[*ast.BlockStmt][]redundantAssertCallMeta)

	insp.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return false
		}

		callExpr := node.(*ast.CallExpr)
		testifyCall := NewCallMeta(pass, callExpr)
		if testifyCall == nil {
			return true
		}

		parentBlock := findNearestNode[*ast.BlockStmt](stack)
		if parentBlock == nil {
			return true
		}

		parentStmt := findNearestNode[*ast.ExprStmt](stack)
		if parentStmt == nil || parentStmt.X != callExpr {
			return true
		}

		callsByBlock[parentBlock] = append(callsByBlock[parentBlock], redundantAssertCallMeta{
			call:        callExpr,
			testifyCall: testifyCall,
			parentStmt:  parentStmt,
		})
		return true
	})

	diagnostics := make([]analysis.Diagnostic, 0)

	for _, calls := range callsByBlock {
		for currIdx, curr := range calls {
			if !curr.testifyCall.IsAssert || curr.testifyCall.Fn.NameFTrimmed != "Error" {
				continue
			}
			if len(curr.testifyCall.Args) < 1 {
				continue
			}

			currErr := analysisutil.NodeString(pass.Fset, curr.testifyCall.Args[0])
			if currErr == "" {
				continue
			}

			closestStrongerFn := ""
			closestDistance := math.MaxInt

			for otherIdx, other := range calls {
				if curr.call == other.call || len(other.testifyCall.Args) < 1 {
					continue
				}

				otherFn := other.testifyCall.Fn.NameFTrimmed
				switch otherFn {
				default:
					continue
				case "ErrorContains", "ErrorIs", "ErrorAs", "EqualError":
				}

				otherErr := analysisutil.NodeString(pass.Fset, other.testifyCall.Args[0])
				if currErr != otherErr {
					continue
				}
				distance := absInt(currIdx - otherIdx)
				if distance >= closestDistance {
					continue
				}
				closestDistance = distance
				closestStrongerFn = other.testifyCall.String()
			}

			if closestStrongerFn != "" {
				diagnostics = append(diagnostics, *newDiagnostic(
					checker.Name(),
					curr.testifyCall,
					fmt.Sprintf("remove assert.%s: %s already asserts that the error is not nil", curr.testifyCall.Fn.Name, closestStrongerFn),
					analysis.SuggestedFix{
						Message: "Remove redundant assert.Error assertion",
						TextEdits: []analysis.TextEdit{{
							Pos:     curr.parentStmt.Pos(),
							End:     toStmtLineEnd(pass, curr.parentStmt.End()),
							NewText: []byte(""),
						}},
					},
				))
			}
		}
	}

	return diagnostics
}

type redundantAssertCallMeta struct {
	call        *ast.CallExpr
	testifyCall *CallMeta
	parentStmt  *ast.ExprStmt
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func toStmtLineEnd(pass *analysis.Pass, pos token.Pos) token.Pos {
	file := pass.Fset.File(pos)
	if file == nil {
		return pos
	}

	line := file.Line(pos)
	if line >= file.LineCount() {
		return pos
	}
	return file.LineStart(line + 1)
}
