package checkers

import (
	"go/ast"
	"go/types"
	"net/http"
	"strconv"

	"golang.org/x/tools/go/analysis"

	"github.com/Antonboom/testifylint/internal/analysisutil"
)

var httpMethod = map[string]string{
	http.MethodGet:     "http.MethodGet",
	http.MethodHead:    "http.MethodHead",
	http.MethodPost:    "http.MethodPost",
	http.MethodPut:     "http.MethodPut",
	http.MethodPatch:   "http.MethodPatch",
	http.MethodDelete:  "http.MethodDelete",
	http.MethodConnect: "http.MethodConnect",
	http.MethodOptions: "http.MethodOptions",
	http.MethodTrace:   "http.MethodTrace",
}

var httpStatusCode = map[string]string{
	strconv.Itoa(http.StatusContinue):           "http.StatusContinue",
	strconv.Itoa(http.StatusSwitchingProtocols): "http.StatusSwitchingProtocols",
	strconv.Itoa(http.StatusProcessing):         "http.StatusProcessing",
	strconv.Itoa(http.StatusEarlyHints):         "http.StatusEarlyHints",

	strconv.Itoa(http.StatusOK):                   "http.StatusOK",
	strconv.Itoa(http.StatusCreated):              "http.StatusCreated",
	strconv.Itoa(http.StatusAccepted):             "http.StatusAccepted",
	strconv.Itoa(http.StatusNonAuthoritativeInfo): "http.StatusNonAuthoritativeInfo",
	strconv.Itoa(http.StatusNoContent):            "http.StatusNoContent",
	strconv.Itoa(http.StatusResetContent):         "http.StatusResetContent",
	strconv.Itoa(http.StatusPartialContent):       "http.StatusPartialContent",
	strconv.Itoa(http.StatusMultiStatus):          "http.StatusMultiStatus",
	strconv.Itoa(http.StatusAlreadyReported):      "http.StatusAlreadyReported",
	strconv.Itoa(http.StatusIMUsed):               "http.StatusIMUsed",

	strconv.Itoa(http.StatusMultipleChoices):   "http.StatusMultipleChoices",
	strconv.Itoa(http.StatusMovedPermanently):  "http.StatusMovedPermanently",
	strconv.Itoa(http.StatusFound):             "http.StatusFound",
	strconv.Itoa(http.StatusSeeOther):          "http.StatusSeeOther",
	strconv.Itoa(http.StatusNotModified):       "http.StatusNotModified",
	strconv.Itoa(http.StatusUseProxy):          "http.StatusUseProxy",
	strconv.Itoa(http.StatusTemporaryRedirect): "http.StatusTemporaryRedirect",
	strconv.Itoa(http.StatusPermanentRedirect): "http.StatusPermanentRedirect",

	strconv.Itoa(http.StatusBadRequest):                   "http.StatusBadRequest",
	strconv.Itoa(http.StatusUnauthorized):                 "http.StatusUnauthorized",
	strconv.Itoa(http.StatusPaymentRequired):              "http.StatusPaymentRequired",
	strconv.Itoa(http.StatusForbidden):                    "http.StatusForbidden",
	strconv.Itoa(http.StatusNotFound):                     "http.StatusNotFound",
	strconv.Itoa(http.StatusMethodNotAllowed):             "http.StatusMethodNotAllowed",
	strconv.Itoa(http.StatusNotAcceptable):                "http.StatusNotAcceptable",
	strconv.Itoa(http.StatusProxyAuthRequired):            "http.StatusProxyAuthRequired",
	strconv.Itoa(http.StatusRequestTimeout):               "http.StatusRequestTimeout",
	strconv.Itoa(http.StatusConflict):                     "http.StatusConflict",
	strconv.Itoa(http.StatusGone):                         "http.StatusGone",
	strconv.Itoa(http.StatusLengthRequired):               "http.StatusLengthRequired",
	strconv.Itoa(http.StatusPreconditionFailed):           "http.StatusPreconditionFailed",
	strconv.Itoa(http.StatusRequestEntityTooLarge):        "http.StatusRequestEntityTooLarge",
	strconv.Itoa(http.StatusRequestURITooLong):            "http.StatusRequestURITooLong",
	strconv.Itoa(http.StatusUnsupportedMediaType):         "http.StatusUnsupportedMediaType",
	strconv.Itoa(http.StatusRequestedRangeNotSatisfiable): "http.StatusRequestedRangeNotSatisfiable",
	strconv.Itoa(http.StatusExpectationFailed):            "http.StatusExpectationFailed",
	strconv.Itoa(http.StatusTeapot):                       "http.StatusTeapot",
	strconv.Itoa(http.StatusMisdirectedRequest):           "http.StatusMisdirectedRequest",
	strconv.Itoa(http.StatusUnprocessableEntity):          "http.StatusUnprocessableEntity",
	strconv.Itoa(http.StatusLocked):                       "http.StatusLocked",
	strconv.Itoa(http.StatusFailedDependency):             "http.StatusFailedDependency",
	strconv.Itoa(http.StatusTooEarly):                     "http.StatusTooEarly",
	strconv.Itoa(http.StatusUpgradeRequired):              "http.StatusUpgradeRequired",
	strconv.Itoa(http.StatusPreconditionRequired):         "http.StatusPreconditionRequired",
	strconv.Itoa(http.StatusTooManyRequests):              "http.StatusTooManyRequests",
	strconv.Itoa(http.StatusRequestHeaderFieldsTooLarge):  "http.StatusRequestHeaderFieldsTooLarge",
	strconv.Itoa(http.StatusUnavailableForLegalReasons):   "http.StatusUnavailableForLegalReasons",

	strconv.Itoa(http.StatusInternalServerError):           "http.StatusInternalServerError",
	strconv.Itoa(http.StatusNotImplemented):                "http.StatusNotImplemented",
	strconv.Itoa(http.StatusBadGateway):                    "http.StatusBadGateway",
	strconv.Itoa(http.StatusServiceUnavailable):            "http.StatusServiceUnavailable",
	strconv.Itoa(http.StatusGatewayTimeout):                "http.StatusGatewayTimeout",
	strconv.Itoa(http.StatusHTTPVersionNotSupported):       "http.StatusHTTPVersionNotSupported",
	strconv.Itoa(http.StatusVariantAlsoNegotiates):         "http.StatusVariantAlsoNegotiates",
	strconv.Itoa(http.StatusInsufficientStorage):           "http.StatusInsufficientStorage",
	strconv.Itoa(http.StatusLoopDetected):                  "http.StatusLoopDetected",
	strconv.Itoa(http.StatusNotExtended):                   "http.StatusNotExtended",
	strconv.Itoa(http.StatusNetworkAuthenticationRequired): "http.StatusNetworkAuthenticationRequired",
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

	for i := 0; i < sig.Params().Len(); i++ {
		lhs := sig.Params().At(i).Type()
		rhs := pass.TypesInfo.TypeOf(fType.Params.List[i].Type)
		if !types.Identical(lhs, rhs) {
			return false
		}
	}
	return true
}
