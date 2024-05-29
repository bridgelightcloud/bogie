package router

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func mockHandler(ctx context.Context, params map[string]string, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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
