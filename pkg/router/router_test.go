package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	router := NewRouter()

	assert.NotNil(router)
}

func TestAddRoute(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	router := NewRouter()
	route := NewRoute("/events", GET, mockHandler)

	router.AddRoute(route)

	assert.Equal(1, len(router.Routes))
	assert.Equal(route.Signature, router.Routes[route.Signature].Signature)

	assert.Panics(func() {
		router.AddRoute(route)
	})
}
