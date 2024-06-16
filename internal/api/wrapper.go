package api

// import (
// 	"fmt"
// 	"net/http"

// 	lambdaEvents "github.com/aws/aws-lambda-go/events"
// )

// type WrapInput struct {
// 	StatusCode int
// 	Body       string
// 	Err        error
// 	Message    string
// }

// func wrap(data WrapInput) lambdaEvents.LambdaFunctionURLResponse {
// 	if ok := checkStatusCode(data.StatusCode); !ok {
// 		panic("Invalid status code")
// 	}

// 	if data.Err == nil {
// 		fmt.Println(data.Message, data.Err.Error())
// 		return lambdaEvents.LambdaFunctionURLResponse{
// 			StatusCode: data.StatusCode,
// 			Body:       data.Message,
// 		}
// 	}

// 	if data.Body != "" {
// 		fmt.Println(data.Message)
// 		return lambdaEvents.LambdaFunctionURLResponse{
// 			StatusCode: data.StatusCode,
// 			Body:       data.Body,
// 		}
// 	}

// 	return lambdaEvents.LambdaFunctionURLResponse{
// 		StatusCode: data.StatusCode,
// 		Body:       data.Message,
// 	}
// }

// func checkStatusCode(c int) bool {
// 	switch c {
// 	case http.StatusContinue:
// 		fallthrough
// 	case http.StatusSwitchingProtocols:
// 		fallthrough
// 	case http.StatusProcessing:
// 		fallthrough
// 	case http.StatusEarlyHints:
// 		fallthrough
// 	case http.StatusOK:
// 		fallthrough

// 	case http.StatusCreated:
// 		fallthrough
// 	case http.StatusAccepted:
// 		fallthrough
// 	case http.StatusNonAuthoritativeInfo:
// 		fallthrough
// 	case http.StatusNoContent:
// 		fallthrough
// 	case http.StatusResetContent:
// 		fallthrough
// 	case http.StatusPartialContent:
// 		fallthrough
// 	case http.StatusMultiStatus:
// 		fallthrough
// 	case http.StatusAlreadyReported:
// 		fallthrough
// 	case http.StatusIMUsed:
// 		fallthrough

// 	case http.StatusMultipleChoices:
// 		fallthrough
// 	case http.StatusMovedPermanently:
// 		fallthrough
// 	case http.StatusFound:
// 		fallthrough
// 	case http.StatusSeeOther:
// 		fallthrough
// 	case http.StatusNotModified:
// 		fallthrough
// 	case http.StatusUseProxy:
// 		fallthrough
// 	case http.StatusTemporaryRedirect:
// 		fallthrough
// 	case http.StatusPermanentRedirect:
// 		fallthrough

// 	case http.StatusBadRequest:
// 		fallthrough
// 	case http.StatusUnauthorized:
// 		fallthrough
// 	case http.StatusPaymentRequired:
// 		fallthrough
// 	case http.StatusForbidden:
// 		fallthrough
// 	case http.StatusNotFound:
// 		fallthrough
// 	case http.StatusMethodNotAllowed:
// 		fallthrough
// 	case http.StatusNotAcceptable:
// 		fallthrough
// 	case http.StatusProxyAuthRequired:
// 		fallthrough
// 	case http.StatusRequestTimeout:
// 		fallthrough
// 	case http.StatusConflict:
// 		fallthrough
// 	case http.StatusGone:
// 		fallthrough
// 	case http.StatusLengthRequired:
// 		fallthrough
// 	case http.StatusPreconditionFailed:
// 		fallthrough
// 	case http.StatusRequestEntityTooLarge:
// 		fallthrough
// 	case http.StatusRequestURITooLong:
// 		fallthrough
// 	case http.StatusUnsupportedMediaType:
// 		fallthrough
// 	case http.StatusRequestedRangeNotSatisfiable:
// 		fallthrough
// 	case http.StatusExpectationFailed:
// 		fallthrough
// 	case http.StatusTeapot:
// 		fallthrough
// 	case http.StatusMisdirectedRequest:
// 		fallthrough
// 	case http.StatusUnprocessableEntity:
// 		fallthrough
// 	case http.StatusLocked:
// 		fallthrough
// 	case http.StatusFailedDependency:
// 		fallthrough
// 	case http.StatusTooEarly:
// 		fallthrough
// 	case http.StatusUpgradeRequired:
// 		fallthrough
// 	case http.StatusPreconditionRequired:
// 		fallthrough
// 	case http.StatusTooManyRequests:
// 		fallthrough
// 	case http.StatusRequestHeaderFieldsTooLarge:
// 		fallthrough
// 	case http.StatusUnavailableForLegalReasons:
// 		fallthrough

// 	case http.StatusInternalServerError:
// 		fallthrough
// 	case http.StatusNotImplemented:
// 		fallthrough
// 	case http.StatusBadGateway:
// 		fallthrough
// 	case http.StatusServiceUnavailable:
// 		fallthrough
// 	case http.StatusGatewayTimeout:
// 		fallthrough
// 	case http.StatusHTTPVersionNotSupported:
// 		fallthrough
// 	case http.StatusVariantAlsoNegotiates:
// 		fallthrough
// 	case http.StatusInsufficientStorage:
// 		fallthrough
// 	case http.StatusLoopDetected:
// 		fallthrough
// 	case http.StatusNotExtended:
// 		fallthrough
// 	case http.StatusNetworkAuthenticationRequired:
// 		return true
// 	default:
// 		return false
// 	}
// }
