package checkers

import (
	"go/ast"
	"regexp"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

// DefaultTimeCompareSuppressCallsPattern contains functions that usually make time equality safe.
var DefaultTimeCompareSuppressCallsPattern = regexp.MustCompile(`Add|AddDate|Date|In|Local|Round|Truncate|UTC`)

// TimeCompare detects flaky time assertions like
//
//	assert.Equal(t, expTime, actualTime)
//	assert.EqualValues(t, expTime, actualTime)
//	assert.Exactly(t, expTime, actualTime)
//	assert.NotEqual(t, expTime, actualTime)
//	assert.NotEqualValues(t, expTime, actualTime)
type TimeCompare struct {
	suppressCallsPattern *regexp.Regexp
}

// NewTimeCompare constructs TimeCompare checker.
func NewTimeCompare() *TimeCompare {
	return &TimeCompare{
		suppressCallsPattern: DefaultTimeCompareSuppressCallsPattern,
	}
}

func (TimeCompare) Name() string { return "time-compare" }

func (checker *TimeCompare) SetSuppressCallsPattern(v *regexp.Regexp) *TimeCompare {
	if v != nil {
		checker.suppressCallsPattern = v
	}
	return checker
}

func (checker TimeCompare) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	fn := call.Fn.NameFTrimmed
	switch fn {
	case "Equal", "EqualValues", "Exactly", "NotEqual", "NotEqualValues":
	default:
		return nil
	}

	if len(call.Args) < 2 {
		return nil
	}

	lhs, rhs := call.Args[0], call.Args[1]

	if isTimeInstance(pass, lhs) || isTimeInstance(pass, rhs) {
		if checker.needSuppressCall(pass, lhs) || checker.needSuppressCall(pass, rhs) {
			return nil
		}

		const report = "equality-based assertion on time.Time can be flaky"
		return newDiagnostic(checker.Name(), call, report)
	}
	return nil
}

func (checker TimeCompare) needSuppressCall(pass *analysis.Pass, e ast.Expr) bool {
	return checker.suppressCallsPattern.Match(analysisutil.NodeBytes(pass.Fset, e))
}

func isTimeInstance(pass *analysis.Pass, e ast.Expr) bool {
	return isNamedType(pass.TypesInfo.TypeOf(e), "time", "Time")
}
