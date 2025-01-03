package csvmum

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsupportedTypeError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	err := &UnsupportedTypeError{}
	assert.Equal("unsupported type: invalid", err.Error())
}

func TestUnmarshalError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	err := &UnmarshalError{Err: fmt.Errorf("some error")}
	assert.Equal("csv unmarshal: some error", err.Error())
	assert.Equal(fmt.Errorf("some error"), err.Unwrap())

	err = &UnmarshalError{
		Err:         fmt.Errorf("some error"),
		FieldIndex:  ptr(1),
		ColumnIndex: ptr(2),
	}
	assert.Equal("csv unmarshal: some error (field index: 1, column index: 2)", err.Error())
}

func TestMarshalError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	err := &MarshalError{Err: fmt.Errorf("some error")}
	assert.Equal("some error", err.Error())
	assert.Equal(fmt.Errorf("some error"), err.Unwrap())
}
