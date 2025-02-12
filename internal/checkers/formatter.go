package checkers

import (
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
	"github.com/Antonboom/testifylint/internal/checkers/printf"
	"github.com/Antonboom/testifylint/internal/testify"
)

// Formatter detects situations like
//
//	assert.ElementsMatch(t, certConfig.Org, csr.Subject.Org, "organizations not equal")
//	assert.Error(t, err, fmt.Sprintf("Profile %s should not be valid", test.profile))
//	assert.Errorf(t, err, fmt.Sprintf("test %s", test.testName))
//	assert.Truef(t, targetTs.Equal(ts), "the timestamp should be as expected (%s) but was %s", targetTs)
//	...
//
// and requires
//
//	assert.ElementsMatchf(t, certConfig.Org, csr.Subject.Org, "organizations not equal")
//	assert.Errorf(t, err, "Profile %s should not be valid", test.profile)
//	assert.Errorf(t, err, "test %s", test.testName)
//	assert.Truef(t, targetTs.Equal(ts), "the timestamp should be as expected (%s) but was %s", targetTs, ts)
//
// It also checks that there are no arguments in `msgAndArgs` if the message is not a string,
// and additionally checks that the first argument of `msgAndArgs` is a string.
type Formatter struct {
	checkFormatString bool
	requireFFuncs     bool
	allowNonStringMsg bool
}

// NewFormatter constructs Formatter checker.
func NewFormatter() *Formatter {
	return &Formatter{
		checkFormatString: true,
		requireFFuncs:     false,
		allowNonStringMsg: false,
	}
}

func (Formatter) Name() string { return "formatter" }

func (checker *Formatter) SetCheckFormatString(v bool) *Formatter {
	checker.checkFormatString = v
	return checker
}

func (checker *Formatter) SetRequireFFuncs(v bool) *Formatter {
	checker.requireFFuncs = v
	return checker
}

func (checker *Formatter) SetAllowNonStringMsg(v bool) *Formatter {
	checker.allowNonStringMsg = v
	return checker
}

func (checker Formatter) Check(pass *analysis.Pass, call *CallMeta) (result *analysis.Diagnostic) {
	if call.Fn.IsFmt {
		return checker.checkFmtAssertion(pass, call)
	}
	return checker.checkNotFmtAssertion(pass, call)
}

func (checker Formatter) checkNotFmtAssertion(pass *analysis.Pass, call *CallMeta) *analysis.Diagnostic {
	msgAndArgsPos, ok := isPrintfLikeCall(pass, call)
	if !ok {
		return nil
	}

	lastArgPos := len(call.ArgsRaw) - 1
	isSingleMsgAndArgElem := msgAndArgsPos == lastArgPos
	msgAndArgs := call.ArgsRaw[msgAndArgsPos]

	if args, ok := isFmtSprintfCall(pass, msgAndArgs); ok && isSingleMsgAndArgElem {
		if checker.requireFFuncs {
			return newRemoveFnAndUseDiagnostic(pass, checker.Name(), call, call.Fn.Name+"f",
				"fmt.Sprintf", msgAndArgs, args...)
		}
		return newRemoveSprintfDiagnostic(pass, checker.Name(), call, msgAndArgs, args)
	}

	if hasStringType(pass, msgAndArgs) { //nolint:nestif // This is the best option of code organization :(
		if checker.requireFFuncs {
			return newUseFunctionDiagnostic(checker.Name(), call, call.Fn.Name+"f")
		}
	} else {
		if isSingleMsgAndArgElem { //nolint:revive // Better without early-return.
			if !checker.allowNonStringMsg {
				return newDiagnostic(checker.Name(), call, "do not use non-string value as first element of msgAndArgs",
					analysis.SuggestedFix{
						Message: `Introduce "%+v" as the message`,
						TextEdits: []analysis.TextEdit{
							{
								Pos:     msgAndArgs.Pos(),
								End:     msgAndArgs.End(),
								NewText: []byte(`"%+v", ` + analysisutil.NodeString(pass.Fset, msgAndArgs)),
							},
						},
					})
			}
		} else {
			return newDiagnostic(checker.Name(), call,
				"using arguments with non-string value as first element of msgAndArgs causes panic",
				analysis.SuggestedFix{
					Message: `Remove format arguments`,
					TextEdits: []analysis.TextEdit{
						{
							Pos:     msgAndArgs.Pos(),
							End:     call.ArgsRaw[lastArgPos].End(),
							NewText: analysisutil.NodeBytes(pass.Fset, msgAndArgs),
						},
					},
				})
		}
	}
	return nil
}

func (checker Formatter) checkFmtAssertion(pass *analysis.Pass, call *CallMeta) (result *analysis.Diagnostic) {
	formatPos := getMsgPosition(call.Fn.Signature)
	if formatPos < 0 {
		return nil
	}

	lastArgPos := len(call.ArgsRaw) - 1
	msg := call.ArgsRaw[formatPos]

	if formatPos == lastArgPos {
		if args, ok := isFmtSprintfCall(pass, msg); ok {
			return newRemoveSprintfDiagnostic(pass, checker.Name(), call, msg, args)
		}
	}

	if checker.checkFormatString {
		report := pass.Report
		defer func() { pass.Report = report }()

		pass.Report = func(d analysis.Diagnostic) {
			result = newDiagnostic(checker.Name(), call, d.Message)
		}

		format, err := strconv.Unquote(analysisutil.NodeString(pass.Fset, msg))
		if err != nil {
			return nil
		}
		printf.CheckPrintf(pass, call.Call, call.String(), format, formatPos)
	}
	return result
}

func isPrintfLikeCall(pass *analysis.Pass, call *CallMeta) (int, bool) {
	msgAndArgsPos := getMsgAndArgsPosition(call.Fn.Signature)
	if msgAndArgsPos <= 0 {
		return -1, false
	}

	if !(len(call.ArgsRaw) > msgAndArgsPos) {
		return -1, false
	}

	if !assertHasFormattedAnalogue(pass, call) {
		return -1, false
	}

	return msgAndArgsPos, true
}

func assertHasFormattedAnalogue(pass *analysis.Pass, call *CallMeta) bool {
	if fn := analysisutil.ObjectOf(pass.Pkg, testify.AssertPkgPath, call.Fn.Name+"f"); fn != nil {
		return true
	}

	if fn := analysisutil.ObjectOf(pass.Pkg, testify.RequirePkgPath, call.Fn.Name+"f"); fn != nil {
		return true
	}

	recv := call.Fn.Signature.Recv()
	if recv == nil {
		return false
	}

	recvT := recv.Type()
	if ptr, ok := recv.Type().(*types.Pointer); ok {
		recvT = ptr.Elem()
	}

	suite, ok := recvT.(*types.Named)
	if !ok {
		return false
	}
	for i := 0; i < suite.NumMethods(); i++ {
		if suite.Method(i).Name() == call.Fn.Name+"f" {
			return true
		}
	}

	return false
}

func getMsgAndArgsPosition(sig *types.Signature) int {
	params := sig.Params()
	if params.Len() < 1 {
		return -1
	}

	lastIdx := params.Len() - 1
	lastParam := params.At(lastIdx)

	_, isSlice := lastParam.Type().(*types.Slice)
	if lastParam.Name() == "msgAndArgs" && isSlice {
		return lastIdx
	}
	return -1
}

func getMsgPosition(sig *types.Signature) int {
	for i := 0; i < sig.Params().Len(); i++ {
		param := sig.Params().At(i)

		if b, ok := param.Type().(*types.Basic); ok && b.Kind() == types.String && param.Name() == "msg" {
			return i
		}
	}
	return -1
}
