package checkers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
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
			parentBlock: parentBlock,
		})
		return true
	})

	diagnostics := make([]analysis.Diagnostic, 0)

	for block, calls := range callsByBlock {
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

				fromStmt, toStmt := curr.parentStmt, other.parentStmt
				if fromStmt.Pos() > toStmt.Pos() {
					fromStmt, toStmt = toStmt, fromStmt
				}
				if isReassignedBetween(block, curr.testifyCall.Args[0], fromStmt, toStmt, pass.Fset, pass.TypesInfo) {
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
	parentBlock *ast.BlockStmt
}

func isReassignedBetween(
	block *ast.BlockStmt,
	errExpr ast.Expr,
	fromStmt *ast.ExprStmt,
	toStmt *ast.ExprStmt,
	fset *token.FileSet,
	info *types.Info,
) bool {
	if block == nil || errExpr == nil || fromStmt == nil || toStmt == nil {
		return false
	}

	var (
		errObj = referencedObject(errExpr, info)
		errStr string
	)
	if errObj == nil {
		errStr = analysisutil.NodeString(fset, errExpr)
		if errStr == "" {
			return false
		}
	}

	fromPos, toPos := fromStmt.Pos(), toStmt.Pos()
	if fromPos > toPos {
		fromPos, toPos = toPos, fromPos
	}

	for _, stmt := range block.List {
		if stmt.Pos() <= fromPos || stmt.Pos() >= toPos {
			continue
		}

		reassigned := false
		ast.Inspect(stmt, func(node ast.Node) bool {
			assignStmt, ok := node.(*ast.AssignStmt)
			if !ok {
				return true
			}

			for _, lhs := range assignStmt.Lhs {
				if sameErrTarget(lhs, errObj, errStr, fset, info) {
					reassigned = true
					return false
				}
			}
			return true
		})
		if reassigned {
			return true
		}
	}

	return false
}

func sameErrTarget(lhs ast.Expr, errObj types.Object, errStr string, fset *token.FileSet, info *types.Info) bool {
	if errObj != nil {
		if lhsObj := referencedObject(lhs, info); lhsObj != nil && lhsObj == errObj {
			return true
		}
	}

	return errStr != "" && analysisutil.NodeString(fset, lhs) == errStr
}

func referencedObject(expr ast.Expr, info *types.Info) types.Object {
	if info == nil {
		return nil
	}

	switch expr := expr.(type) {
	case *ast.Ident:
		return info.ObjectOf(expr)
	case *ast.SelectorExpr:
		if sel, ok := info.Selections[expr]; ok && sel != nil {
			return sel.Obj()
		}
		return info.ObjectOf(expr.Sel)
	default:
		return nil
	}
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
