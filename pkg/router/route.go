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

type HandlerFunc func(context.Context, map[string]string, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

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
