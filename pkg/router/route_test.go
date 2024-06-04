package router

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func mockHandler(ctx context.Context, params PathParams, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{}, nil
}

func TestNewRoute(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	path := "/events"
	method := GET
	route := NewRoute(path, method, mockHandler)

	assert.Equal(path, route.Path)
	assert.Equal(2, len(route.PathElements))
	assert.Equal(method, route.Method)
	assert.NotNil(route.Handler)
	assert.Equal(string(method)+path, route.Signature)
}

func TestRouteMustStartWithSlash(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	assert.Panics(func() {
		NewRoute("events", GET, mockHandler)
	})
}

func TestMatch(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name          string
		routePath     string
		routeMethod   Method
		requestPath   string
		requestMethod Method
		expected      bool
	}{
		{
			name:          "Match",
			routePath:     "/events",
			routeMethod:   GET,
			requestPath:   "/events",
			requestMethod: GET,
			expected:      true,
		},
		{
			name:          "MatchWithParams",
			routePath:     "/events/:id",
			routeMethod:   GET,
			requestPath:   "/events/1",
			requestMethod: GET,
			expected:      true,
		},
		{
			name:          "NoMatchWrongMethod",
			routePath:     "/events",
			routeMethod:   GET,
			requestPath:   "/events",
			requestMethod: POST,
			expected:      false,
		},
		{
			name:          "NoMatchDifferentPath",
			routePath:     "/events",
			routeMethod:   GET,
			requestPath:   "/events/1",
			requestMethod: GET,
			expected:      false,
		},
		{
			name:          "NoMatchDifferentPath",
			routePath:     "/events/abc",
			routeMethod:   GET,
			requestPath:   "/events/1",
			requestMethod: GET,
			expected:      false,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			route := NewRoute(tc.routePath, tc.routeMethod, mockHandler)

			handler, match := route.match(tc.requestPath, tc.requestMethod)
			assert.Equal(tc.expected, match)
			if tc.expected {
				assert.NotNil(handler)
			} else {
				assert.Nil(handler)
			}
		})
	}
}

func TestParams(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		routePath   string
		requestPath string
		expected    PathParams
	}{
		{
			name:        "NoParams",
			routePath:   "/events",
			requestPath: "/events",
			expected:    PathParams{},
		},
		{
			name:        "SingleParam",
			routePath:   "/events/:id",
			requestPath: "/events/1",
			expected: PathParams{
				"id": "1",
			},
		},
		{
			name:        "MultipleParams",
			routePath:   "/events/:id/:name",
			requestPath: "/events/1/john",
			expected: PathParams{
				"id":   "1",
				"name": "john",
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			route := NewRoute(tc.routePath, GET, mockHandler)

			params := route.params(tc.requestPath)
			assert.Equal(tc.expected, params)
		})
	}
}