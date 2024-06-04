package router

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
	ANY     Method = "*"
)

type PathParams map[string]string

type HandlerFunc func(context.Context, PathParams, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

type Route struct {
	ParamCount        int
	Path              string
	PathElements      []string
	PathLength        int
	Method            Method
	Handler           HandlerFunc
	Signature         string
	SignatureElements []string
}

func NewRoute(path string, method Method, handler HandlerFunc) Route {
	if !strings.HasPrefix(path, "/") {
		panic("Path must start with /")
	}

	paramCount := 0
	pathElements := strings.Split(path, "/")
	signatureElements := []string{}
	for _, element := range pathElements {
		if strings.HasPrefix(element, ":") {
			signatureElements = append(signatureElements, "*")
			paramCount++
		} else {
			signatureElements = append(signatureElements, element)
		}
	}

	return Route{
		ParamCount:        paramCount,
		Path:              path,
		PathElements:      pathElements,
		PathLength:        len(pathElements),
		Method:            method,
		Handler:           handler,
		Signature:         string(method) + strings.Join(signatureElements, "/"),
		SignatureElements: signatureElements,
	}
}

func (r *Route) match(path string, method Method) (HandlerFunc, bool) {
	if r.Method != method {
		return nil, false
	}

	pathElements := strings.Split(path, "/")

	if len(pathElements) != len(r.PathElements) {
		return nil, false
	}

	for i, element := range r.PathElements {
		if element != pathElements[i] && !strings.HasPrefix(element, ":") {
			return nil, false
		}
	}

	return r.Handler, true
}

func (r *Route) params(path string) PathParams {
	pathElements := strings.Split(path, "/")
	params := PathParams{}

	for i, element := range r.PathElements {
		if strings.HasPrefix(element, ":") {
			params[element[1:]] = pathElements[i]
		}
	}

	return params
}
