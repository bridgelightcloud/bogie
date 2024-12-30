package csvmum

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsers(t *testing.T) {
	t.Parallel()

	t.Run("parseString", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc string = "string"
		field := reflect.Value{}
		value := "string"
		v, err := parseString(field, value)

		assert.Nil(err, "%w", err)
		assert.Equal(tc, v.String())
	})

	t.Run("parseInt", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc int = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseInt(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseInt8", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc int8 = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseInt8(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseInt16", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc int16 = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseInt16(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("praseInt32", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc int32 = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseInt32(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseInt64", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc int64 = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseInt64(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseUint", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc uint = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseUint(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseUint8", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc uint8 = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseUint8(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseUint16", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc uint16 = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseUint16(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseUint32", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc uint32 = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseUint32(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseUint64", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc uint64 = 1
		field := reflect.Value{}
		value := "1"
		v, err := parseUint64(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseFloat32", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc float32 = 1.1
		field := reflect.Value{}
		value := "1.1"
		v, err := parseFloat32(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseFloat64", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		var tc float64 = 1.1
		field := reflect.Value{}
		value := "1.1"
		v, err := parseFloat64(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Interface())
	})

	t.Run("parseBool", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		field := reflect.Value{}
		value := "true"
		v, err := parseBool(field, value)

		assert.Nil(err)
		assert.True(v.Bool())
	})

	t.Run("parsePointer", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		tc := 1
		field := reflect.ValueOf(ptr(tc))
		value := "1"
		v, err := parsePointer(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Elem().Interface())
	})

	t.Run("parseInvalidPointer", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		tc := 1
		field := reflect.ValueOf(ptr(tc))
		value := "seven"
		_, err := parsePointer(field, value)

		assert.Equal(err.Error(), "strconv.ParseInt: parsing \"seven\": invalid syntax")
	})

	t.Run("invalidType", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		field := reflect.Value{}
		value := "1"
		_, err := parseValue(field, value)

		assert.Equal(err.Error(), "unsupported type invalid")
	})

	t.Run("parsePointerToPointer", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		tc := 1
		field := reflect.ValueOf(ptr(ptr(tc)))
		value := "1"
		v, err := parsePointer(field, value)

		assert.Nil(err)
		assert.Equal(tc, v.Elem().Elem().Interface())
	})
}
