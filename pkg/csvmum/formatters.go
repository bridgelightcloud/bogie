package csvmum

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

func formatString(v reflect.Value) (string, error) {
	return v.String(), nil
}

func formatInt(v reflect.Value) (string, error) {
	return strconv.FormatInt(v.Int(), 10), nil
}

func formatUint(v reflect.Value) (string, error) {
	return strconv.FormatUint(v.Uint(), 10), nil
}

func formatFloat32(v reflect.Value) (string, error) {
	return strconv.FormatFloat(v.Float(), 'f', -1, 32), nil
}

func formatFloat64(v reflect.Value) (string, error) {
	return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
}

func formatBool(v reflect.Value) (string, error) {
	return strconv.FormatBool(v.Bool()), nil
}

func formatPointer(v reflect.Value) (string, error) {
	if v.IsNil() {
		return "", nil
	}

	return formatValue(v.Elem())
}

func formatTextMarshaler(tm encoding.TextMarshaler) (string, error) {
	b, err := tm.MarshalText()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func formatValue(v reflect.Value) (string, error) {
	if tm, ok := v.Interface().(encoding.TextMarshaler); ok {
		return formatTextMarshaler(tm)
	}

	if formatter, ok := defaultFormatters[v.Kind()]; ok {
		return formatter(v)
	}

	return "", fmt.Errorf("unsupported type: %s", v.Kind())
}

type valueFormatter func(reflect.Value) (string, error)

var defaultFormatters map[reflect.Kind]valueFormatter

func init() {
	defaultFormatters = map[reflect.Kind]valueFormatter{
		reflect.String:  formatString,
		reflect.Int:     formatInt,
		reflect.Int8:    formatInt,
		reflect.Int16:   formatInt,
		reflect.Int32:   formatInt,
		reflect.Int64:   formatInt,
		reflect.Uint:    formatUint,
		reflect.Uint8:   formatUint,
		reflect.Uint16:  formatUint,
		reflect.Uint32:  formatUint,
		reflect.Uint64:  formatUint,
		reflect.Float32: formatFloat32,
		reflect.Float64: formatFloat64,
		reflect.Bool:    formatBool,
		reflect.Pointer: formatPointer,
	}
}
