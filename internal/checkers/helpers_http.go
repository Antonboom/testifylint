package checkers

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"net/http"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

var httpMethod = map[string]string{
	http.MethodGet:     "MethodGet",
	http.MethodHead:    "MethodHead",
	http.MethodPost:    "MethodPost",
	http.MethodPut:     "MethodPut",
	http.MethodPatch:   "MethodPatch",
	http.MethodDelete:  "MethodDelete",
	http.MethodConnect: "MethodConnect",
	http.MethodOptions: "MethodOptions",
	http.MethodTrace:   "MethodTrace",
}

var httpStatusCode = map[int]string{
	http.StatusContinue:           "StatusContinue",
	http.StatusSwitchingProtocols: "StatusSwitchingProtocols",
	http.StatusProcessing:         "StatusProcessing",
	http.StatusEarlyHints:         "StatusEarlyHints",

	http.StatusOK:                   "StatusOK",
	http.StatusCreated:              "StatusCreated",
	http.StatusAccepted:             "StatusAccepted",
	http.StatusNonAuthoritativeInfo: "StatusNonAuthoritativeInfo",
	http.StatusNoContent:            "StatusNoContent",
	http.StatusResetContent:         "StatusResetContent",
	http.StatusPartialContent:       "StatusPartialContent",
	http.StatusMultiStatus:          "StatusMultiStatus",
	http.StatusAlreadyReported:      "StatusAlreadyReported",
	http.StatusIMUsed:               "StatusIMUsed",

	http.StatusMultipleChoices:   "StatusMultipleChoices",
	http.StatusMovedPermanently:  "StatusMovedPermanently",
	http.StatusFound:             "StatusFound",
	http.StatusSeeOther:          "StatusSeeOther",
	http.StatusNotModified:       "StatusNotModified",
	http.StatusUseProxy:          "StatusUseProxy",
	http.StatusTemporaryRedirect: "StatusTemporaryRedirect",
	http.StatusPermanentRedirect: "StatusPermanentRedirect",

	http.StatusBadRequest:                   "StatusBadRequest",
	http.StatusUnauthorized:                 "StatusUnauthorized",
	http.StatusPaymentRequired:              "StatusPaymentRequired",
	http.StatusForbidden:                    "StatusForbidden",
	http.StatusNotFound:                     "StatusNotFound",
	http.StatusMethodNotAllowed:             "StatusMethodNotAllowed",
	http.StatusNotAcceptable:                "StatusNotAcceptable",
	http.StatusProxyAuthRequired:            "StatusProxyAuthRequired",
	http.StatusRequestTimeout:               "StatusRequestTimeout",
	http.StatusConflict:                     "StatusConflict",
	http.StatusGone:                         "StatusGone",
	http.StatusLengthRequired:               "StatusLengthRequired",
	http.StatusPreconditionFailed:           "StatusPreconditionFailed",
	http.StatusRequestEntityTooLarge:        "StatusRequestEntityTooLarge",
	http.StatusRequestURITooLong:            "StatusRequestURITooLong",
	http.StatusUnsupportedMediaType:         "StatusUnsupportedMediaType",
	http.StatusRequestedRangeNotSatisfiable: "StatusRequestedRangeNotSatisfiable",
	http.StatusExpectationFailed:            "StatusExpectationFailed",
	http.StatusTeapot:                       "StatusTeapot",
	http.StatusMisdirectedRequest:           "StatusMisdirectedRequest",
	http.StatusUnprocessableEntity:          "StatusUnprocessableEntity",
	http.StatusLocked:                       "StatusLocked",
	http.StatusFailedDependency:             "StatusFailedDependency",
	http.StatusTooEarly:                     "StatusTooEarly",
	http.StatusUpgradeRequired:              "StatusUpgradeRequired",
	http.StatusPreconditionRequired:         "StatusPreconditionRequired",
	http.StatusTooManyRequests:              "StatusTooManyRequests",
	http.StatusRequestHeaderFieldsTooLarge:  "StatusRequestHeaderFieldsTooLarge",
	http.StatusUnavailableForLegalReasons:   "StatusUnavailableForLegalReasons",

	http.StatusInternalServerError:           "StatusInternalServerError",
	http.StatusNotImplemented:                "StatusNotImplemented",
	http.StatusBadGateway:                    "StatusBadGateway",
	http.StatusServiceUnavailable:            "StatusServiceUnavailable",
	http.StatusGatewayTimeout:                "StatusGatewayTimeout",
	http.StatusHTTPVersionNotSupported:       "StatusHTTPVersionNotSupported",
	http.StatusVariantAlsoNegotiates:         "StatusVariantAlsoNegotiates",
	http.StatusInsufficientStorage:           "StatusInsufficientStorage",
	http.StatusLoopDetected:                  "StatusLoopDetected",
	http.StatusNotExtended:                   "StatusNotExtended",
	http.StatusNetworkAuthenticationRequired: "StatusNetworkAuthenticationRequired",
}

// httpNetPkgName returns the local name for "net/http" in the file containing pos,
// and an optional TextEdit to add the import when absent (if no conflict exists).
// Returns ("", nil, false) if net/http is blank-imported or all candidate names are taken.
func httpNetPkgName(pass *analysis.Pass, pos token.Pos) (string, *analysis.TextEdit, bool) {
	return AddImportFix(pass.Files, pos, "net/http")
}

func newHTTPMethodTextEdit(pass *analysis.Pass, e ast.Expr) (valueEdit *analysis.TextEdit, importEdit *analysis.TextEdit) {
	bt, ok := typeSafeBasicLit(e, token.STRING)
	if !ok {
		return nil, nil
	}
	currentVal, ok := unquoteBasicLitValue(bt)
	if !ok {
		return nil, nil
	}
	constName, ok := httpMethod[strings.ToUpper(currentVal)]
	if !ok {
		return nil, nil
	}
	newVal, importEdit, ok := httpQualifiedName(pass, bt.Pos(), constName)
	if !ok {
		return nil, nil
	}
	return &analysis.TextEdit{
		Pos:     bt.Pos(),
		End:     bt.End(),
		NewText: []byte(newVal),
	}, importEdit
}

func newHTTPStatusCodeTextEdit(pass *analysis.Pass, e ast.Expr) (valueEdit *analysis.TextEdit, importEdit *analysis.TextEdit) {
	bt, ok := typeSafeBasicLit(e, token.INT)
	if !ok {
		return nil, nil
	}
	// Use go/constant to parse the literal, correctly handling all Go integer literal forms
	// (decimal, hex 0xC8, octal 0o310, binary 0b11001000, underscore separators 2_00, etc.).
	v := constant.MakeFromLiteral(bt.Value, token.INT, 0)
	if v.Kind() != constant.Int {
		return nil, nil
	}
	intVal, exact := constant.Int64Val(v)
	if !exact {
		return nil, nil
	}
	// Guard against overflow when converting int64 to int on 32-bit platforms.
	if int64(int(intVal)) != intVal {
		return nil, nil
	}
	constName, ok := httpStatusCode[int(intVal)]
	if !ok {
		return nil, nil
	}
	newVal, importEdit, ok := httpQualifiedName(pass, bt.Pos(), constName)
	if !ok {
		return nil, nil
	}
	return &analysis.TextEdit{
		Pos:     bt.Pos(),
		End:     bt.End(),
		NewText: []byte(newVal),
	}, importEdit
}

// httpQualifiedName returns ("qualifier.constName", importEdit, true) using the local import name of net/http,
// or ("constName", nil, true) if net/http is dot-imported, or ("", nil, false) if net/http
// cannot be provided (blank-imported or all candidate names are taken).
// importEdit is non-nil when net/http is absent and needs to be added.
func httpQualifiedName(pass *analysis.Pass, pos token.Pos, constName string) (string, *analysis.TextEdit, bool) {
	qualifier, importEdit, ok := httpNetPkgName(pass, pos)
	if !ok {
		return "", nil, false
	}
	if qualifier == "" {
		return constName, nil, true // dot-import: no qualifier needed
	}
	return qualifier + "." + constName, importEdit, true
}

func mimicHTTPHandler(pass *analysis.Pass, fType *ast.FuncType) bool {
	httpHandlerFuncObj := analysisutil.ObjectOf(pass.Pkg, "net/http", "HandlerFunc")
	if httpHandlerFuncObj == nil {
		return false
	}

	sig, ok := httpHandlerFuncObj.Type().Underlying().(*types.Signature)
	if !ok {
		return false
	}

	if len(fType.Params.List) != sig.Params().Len() {
		return false
	}

	for i := range sig.Params().Len() {
		lhs := sig.Params().At(i).Type()
		rhs := pass.TypesInfo.TypeOf(fType.Params.List[i].Type)
		if !types.Identical(lhs, rhs) {
			return false
		}
	}
	return true
}
