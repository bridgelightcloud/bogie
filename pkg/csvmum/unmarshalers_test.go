package csvmum

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalValue(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		zero     any
		value    string
		expected any
		err      any
	}{
		"string":                   {"", "string", "string", nil},
		"int":                      {0, "27", 27, nil},
		"intParseErr":              {0, "twenty-seven", 0, ptr[*strconv.NumError](nil)},
		"intOverflow":              {0, "9223372036854775808", 0, ptr[*strconv.NumError](nil)},
		"int8":                     {int8(0), "27", int8(27), nil},
		"int8ParseErr":             {int8(0), "twenty-seven", int8(0), ptr[*strconv.NumError](nil)},
		"int8Overflow":             {int8(0), "128", int8(0), ptr[*strconv.NumError](nil)},
		"int16":                    {int16(0), "27", int16(27), nil},
		"int16ParseErr":            {int16(0), "twenty-seven", int16(0), ptr[*strconv.NumError](nil)},
		"int16Overflow":            {int16(0), "32768", int16(0), ptr[*strconv.NumError](nil)},
		"int32":                    {int32(0), "27", int32(27), nil},
		"int32ParseErr":            {int32(0), "twenty-seven", int32(0), ptr[*strconv.NumError](nil)},
		"int32Overflow":            {int32(0), "2147483648", int32(0), ptr[*strconv.NumError](nil)},
		"int64":                    {int64(0), "27", int64(27), nil},
		"int64ParseErr":            {int64(0), "twenty-seven", int64(0), ptr[*strconv.NumError](nil)},
		"int64Overflow":            {int64(0), "9223372036854775808", int64(0), ptr[*strconv.NumError](nil)},
		"uint":                     {uint(0), "27", uint(27), nil},
		"uintParseErr":             {uint(0), "twenty-seven", uint(0), ptr[*strconv.NumError](nil)},
		"uintOverflow":             {uint(0), "18446744073709551616", uint(0), ptr[*strconv.NumError](nil)},
		"uint8":                    {uint8(0), "27", uint8(27), nil},
		"uint8ParseErr":            {uint8(0), "twenty-seven", uint8(0), ptr[*strconv.NumError](nil)},
		"uint8Overflow":            {uint8(0), "256", uint8(0), ptr[*strconv.NumError](nil)},
		"uint16":                   {uint16(0), "27", uint16(27), nil},
		"uint16ParseErr":           {uint16(0), "twenty-seven", uint16(0), ptr[*strconv.NumError](nil)},
		"uint16Overflow":           {uint16(0), "65536", uint16(0), ptr[*strconv.NumError](nil)},
		"uint32":                   {uint32(0), "27", uint32(27), nil},
		"uint32ParseErr":           {uint32(0), "twenty-seven", uint32(0), ptr[*strconv.NumError](nil)},
		"uint32Overflow":           {uint32(0), "4294967296", uint32(0), ptr[*strconv.NumError](nil)},
		"uint64":                   {uint64(0), "27", uint64(27), nil},
		"uint64ParseErr":           {uint64(0), "twenty-seven", uint64(0), ptr[*strconv.NumError](nil)},
		"uint64Overflow":           {uint64(0), "18446744073709551616", uint64(0), ptr[*strconv.NumError](nil)},
		"float32":                  {float32(0), "27.27", float32(27.27), nil},
		"float32ParseErr":          {float32(0), "twenty-seven", float32(0), ptr[*strconv.NumError](nil)},
		"float32Overflow":          {float32(0), "1e39", float32(0), ptr[*strconv.NumError](nil)},
		"float64":                  {float64(0), "27.27", float64(27.27), nil},
		"float64ParseErr":          {float64(0), "twenty-seven", float64(0), ptr[*strconv.NumError](nil)},
		"float64Overflow":          {float64(0), "1e309", float64(0), ptr[*strconv.NumError](nil)},
		"bool":                     {false, "true", true, nil},
		"boolParseErr":             {false, "twenty-seven", false, ptr[*strconv.NumError](nil)},
		"boolEmptyErr":             {false, " ", false, ptr[*strconv.NumError](nil)},
		"pointer":                  {ptr(0), "27", ptr(27), nil},
		"pointer empty":            {ptr(0), "", nilPtr[int](), nil},
		"pointer to pointer":       {ptr(ptr(0)), "27", ptr(ptr(27)), nil},
		"customTextMarshaler":      {customMarshalAndUnmarshal{}, "~string~", customMarshalAndUnmarshal{"string"}, nil},
		"customTextMarshalerError": {customMarshalAndUnmarshal{}, "", customMarshalAndUnmarshal{}, ptr[error](nil)},
		"invalid":                  {struct{}{}, "27", struct{}{}, ptr[*UnsupportedTypeError](nil)},
	}

	for name, tc := range tt {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			field := reflect.New(reflect.TypeOf(tc.zero)).Elem()
			err := unmarshalValue(tc.value, field)

			if tc.err != nil {
				assert.ErrorAs(err, tc.err)
			} else {
				assert.NoError(err)
			}

			assert.Equal(tc.expected, field.Interface())
		})
	}
}
