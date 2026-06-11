package checkers

import (
	"go/ast"
	"regexp"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

var DefaultTimeEqualitySuppressCallsPattern = regexp.MustCompile(`Add|AddDate|Date|In|Local|Round|Truncate|UTC`)

// TimeCompare detects situations like
//
//	assert.Equal(t, expTime, actualTime)
//
// and requires
//
// TODO.
type TimeCompare struct {
	warnOnTimeEquality               bool
	timeEqualitySuppressCallsPattern *regexp.Regexp
}

// NewTimeCompare constructs TimeCompare checker.
func NewTimeCompare() *TimeCompare {
	return &TimeCompare{
		warnOnTimeEquality:               false,
		timeEqualitySuppressCallsPattern: DefaultTimeEqualitySuppressCallsPattern,
	}
}

func (TimeCompare) Name() string { return "time-compare" }

func (checker *TimeCompare) SetWarnOnTimeEquality(v bool) *TimeCompare {
	checker.warnOnTimeEquality = v
	return checker
}

func (checker *TimeCompare) SetTimeEqualitySuppressCallsPattern(v *regexp.Regexp) *TimeCompare {
	if v != nil {
		checker.timeEqualitySuppressCallsPattern = v
	}
	return checker
}

func (checker TimeCompare) Check(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if d := checker.checkSimplification(pass, call); d != nil {
		return d
	}

	return checker.checkTimeEquality(pass, call)
}

func (checker TimeCompare) checkSimplification(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	return nil
}

func (checker TimeCompare) checkTimeEquality(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	if !checker.warnOnTimeEquality {
		return nil
	}

	switch call.Fn.NameFTrimmed {
	case "Equal", "EqualValues", "Exactly", "NotEqual", "NotEqualValues", "NotExactly":
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
	return checker.timeEqualitySuppressCallsPattern.Match(analysisutil.NodeBytes(pass.Fset, e))
}

func isTimeInstance(pass *analysis.Pass, e ast.Expr) bool {
	return isNamedType(pass.TypesInfo.TypeOf(e), "time", "Time")
}
