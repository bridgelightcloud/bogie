package csvmum

import (
	"fmt"
	"reflect"
	"strconv"
)

type parser func(reflect.Value, string) (reflect.Value, error)
type parsers map[reflect.Kind]parser

var defaultParsers parsers

func parseString(field reflect.Value, value string) (reflect.Value, error) {
	return reflect.ValueOf(value), nil
}

func parseInt(field reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseInt(value, 10, 64)
	return reflect.ValueOf(int(i)), err
}

func parseFloat64(field reflect.Value, value string) (reflect.Value, error) {
	f, err := strconv.ParseFloat(value, 64)
	return reflect.ValueOf(f), err
}

func parseBool(field reflect.Value, value string) (reflect.Value, error) {
	b, err := strconv.ParseBool(value)
	return reflect.ValueOf(b), err
}

func parsePointer(field reflect.Value, value string) (reflect.Value, error) {
	if value == "" {
		return reflect.New(field.Type()).Elem(), nil
	}

	if field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
	}

	v, err := parseValue(field.Elem(), value)
	if err != nil {
		return reflect.Value{}, err
	}

	if v.CanAddr() {
		return v.Addr(), nil
	}
	newValue := reflect.New(v.Type()).Elem()
	newValue.Set(v)
	return newValue.Addr(), nil
}

func parseValue(field reflect.Value, s string) (reflect.Value, error) {
	if f, ok := defaultParsers[field.Kind()]; ok {
		return f(field, s)
	}
	return reflect.Value{}, fmt.Errorf("unsupported type %s", field.Kind())
}

func init() {
	defaultParsers = parsers{
		reflect.String:  parseString,
		reflect.Int:     parseInt,
		reflect.Float64: parseFloat64,
		reflect.Bool:    parseBool,
		reflect.Pointer: parsePointer,
	}
}
