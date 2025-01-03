package csvmum

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalValue(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		field    any
		expected string
		err      error
	}{
		"string":                    {"string", "string", nil},
		"int":                       {27, "27", nil},
		"int8":                      {int8(27), "27", nil},
		"int16":                     {int16(27), "27", nil},
		"int32":                     {int32(27), "27", nil},
		"int64":                     {int64(27), "27", nil},
		"uint":                      {uint(27), "27", nil},
		"uint8":                     {uint8(27), "27", nil},
		"uint16":                    {uint16(27), "27", nil},
		"uint32":                    {uint32(27), "27", nil},
		"uint64":                    {uint64(27), "27", nil},
		"float32":                   {float32(3.14), "3.14", nil},
		"float64":                   {3.14, "3.14", nil},
		"bool":                      {true, "true", nil},
		"pointer":                   {ptr(27), "27", nil},
		"pointerNil":                {nilPtr[int](), "", nil},
		"textMarshaler":             {customMarshalAndUnmarshal{One: "one"}, "~one~", nil},
		"textMarshalerInvalidValue": {customMarshalAndUnmarshal{}, "", errors.New("invalid text: ")},
		"unsupported":               {struct{}{}, "", errors.New("unsupported type: struct")},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			v, err := marshalValue(reflect.ValueOf(tc.field))

			if tc.err != nil {
				assert.EqualError(err, tc.err.Error())
			} else {
				assert.NoError(err)
			}

			assert.Equal(tc.expected, v)
		})
	}
}
