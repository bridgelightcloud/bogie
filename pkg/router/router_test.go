package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	router := NewRouter([]Route{})

	assert.NotNil(router)
}

func TestAddRoute(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	route1 := NewRoute("/events", GET, mockHandler)
	route2 := NewRoute("/events", GET, mockHandler)
	router := NewRouter([]Route{
		route1,
	})

	assert.Equal(1, len(router.Routes))
	assert.Equal(route1.Signature, router.Routes[route1.Signature].Signature)

	assert.Panics(func() {
		router.addRoute(route2)
	})
}
