package csvmum

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsupportedTypeError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	err := &UnsupportedTypeError{}
	assert.Equal("unsupported type: invalid", err.Error())

	err = &UnsupportedTypeError{reflect.Float32}
	assert.Equal("unsupported type: float32", err.Error())
}

func TestParseError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	err := &ParseError{fmt.Errorf("some error")}
	assert.Equal("some error", err.Error())
	assert.Equal(fmt.Errorf("some error"), err.Unwrap())
}
