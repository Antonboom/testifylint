package checkers

import (
	"fmt"
	"go/types"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"golang.org/x/tools/go/analysis"
)

// TimeEqual detects situations like
//
//	assert.Equal(t, timeA, timeB)
//	assert.NotEqual(t, timeA, timeB)
//
// and requires
//
//	assert.True(t, timeA.Equal(timeB))
//	assert.False(t, timeA.Equal(timeB))

type TimeEqual struct{}

func NewTimeEqual() TimeEqual {
	return TimeEqual{}
}

func (TimeEqual) Name() string {
	return "time-compare"
}

func (checker TimeEqual) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	assrn := call.Fn.NameFTrimmed

	var isEq bool
	switch assrn {
	default:
		return nil
	case "Equal":
		isEq = true
	case "NotEqual":
		isEq = false
	}

	if len(call.Args) < 2 {
		return nil
	}

	first, second := call.Args[0], call.Args[1]

	ft, st := pass.TypesInfo.TypeOf(first), pass.TypesInfo.TypeOf(second)

	if !isTimeType(ft) || !isTimeType(st) {
		return nil
	}

	if isEq {
		return newDiagnostic(
			checker.Name(),
			call,
			fmt.Sprintf("replace %s.Equal with %s.True passing time.Equal", call.SelectorXStr, call.SelectorXStr),
			newSuggestedFuncReplacement(
				call,
				"True",
				analysis.TextEdit{
					Pos:     first.Pos(),
					End:     second.End(),
					NewText: []byte(analysisutil.NodeString(pass.Fset, first) + ".Equal(" + analysisutil.NodeString(pass.Fset, second) + ")"),
				},
			),
		)
	} else {
		return newDiagnostic(
			checker.Name(),
			call,
			fmt.Sprintf("replace %s.NotEqual with %s.False passing time.Equal", call.SelectorXStr, call.SelectorXStr),
			newSuggestedFuncReplacement(
				call,
				"False",
				analysis.TextEdit{
					Pos:     first.Pos(),
					End:     second.End(),
					NewText: []byte(analysisutil.NodeString(pass.Fset, first) + ".Equal(" + analysisutil.NodeString(pass.Fset, second) + ")"),
				},
			),
		)
	}
}

func isTimeType(typ types.Type) bool {
	n, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	typeName := n.Obj()
	return typeName != nil && typeName.Pkg() != nil && typeName.Pkg().Path() == "time" && typeName.Name() == "Time"
}
