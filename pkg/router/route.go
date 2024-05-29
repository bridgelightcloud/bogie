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
)

type Params map[string]string

type HandlerFunc func(context.Context, Params, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

type Route struct {
	Path         string
	PathElements []string
	Method       Method
	Handler      HandlerFunc
	Signature    string
}

func NewRoute(path string, method Method, handler HandlerFunc) Route {
	return Route{
		Path:         path,
		PathElements: strings.Split(path, "/"),
		Method:       method,
		Handler:      handler,
		Signature:    string(method) + path,
	}
}

func (r *Route) Match(path string, method Method) (HandlerFunc, bool) {
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

func (r *Route) Params(path string) Params {
	pathElements := strings.Split(path, "/")
	params := Params{}

	for i, element := range r.PathElements {
		if strings.HasPrefix(element, ":") {
			params[element[1:]] = pathElements[i]
		}
	}

	return params
}
