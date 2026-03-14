package checkers

import (
	"go/ast"
	"go/token"
	"go/types"
	"net/http"

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
	return addImportFix(pass.Files, pos, "net/http")
}

// httpNetQualifiedConst returns the qualified net/http constant reference for constName
// using httpQual as the local package name.
// Returns "constName" when httpQual is "" (dot-import), otherwise "httpQual.constName".
func httpNetQualifiedConst(httpQual, constName string) string {
	if httpQual == "" {
		return constName
	}
	return httpQual + "." + constName
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
