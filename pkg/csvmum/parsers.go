package csvmum

import (
	"fmt"
	"reflect"
	"strconv"
)

type parser func(reflect.Value, string) (reflect.Value, error)
type parsers map[reflect.Kind]parser

var defaultParsers parsers

func parseString(_ reflect.Value, value string) (reflect.Value, error) {
	return reflect.ValueOf(value), nil
}

func parseInt(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseInt(value, 10, 64)
	return reflect.ValueOf(int(i)), err
}

func parseInt8(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseInt(value, 10, 8)
	return reflect.ValueOf(int8(i)), err
}

func parseInt16(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseInt(value, 10, 16)
	return reflect.ValueOf(int16(i)), err
}

func parseInt32(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseInt(value, 10, 32)
	return reflect.ValueOf(int32(i)), err
}

func parseInt64(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseInt(value, 10, 64)
	return reflect.ValueOf(i), err
}

func parseUint(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseUint(value, 10, 64)
	return reflect.ValueOf(uint(i)), err
}

func parseUint8(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseUint(value, 10, 8)
	return reflect.ValueOf(uint8(i)), err
}

func parseUint16(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseUint(value, 10, 16)
	return reflect.ValueOf(uint16(i)), err
}

func parseUint32(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseUint(value, 10, 32)
	return reflect.ValueOf(uint32(i)), err
}

func parseUint64(_ reflect.Value, value string) (reflect.Value, error) {
	i, err := strconv.ParseUint(value, 10, 64)
	return reflect.ValueOf(i), err
}

func parseFloat32(_ reflect.Value, value string) (reflect.Value, error) {
	f, err := strconv.ParseFloat(value, 32)
	return reflect.ValueOf(float32(f)), err
}

func parseFloat64(_ reflect.Value, value string) (reflect.Value, error) {
	f, err := strconv.ParseFloat(value, 64)
	return reflect.ValueOf(f), err
}

func parseBool(_ reflect.Value, value string) (reflect.Value, error) {
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

	newValue := reflect.New(v.Type()).Elem()
	newValue.Set(v)
	return newValue.Addr(), nil
}

func parseValue(field reflect.Value, value string) (reflect.Value, error) {
	if p, ok := defaultParsers[field.Kind()]; ok {
		return p(field, value)
	}
	return reflect.Value{}, fmt.Errorf("unsupported type %s", field.Kind())
}

func init() {
	defaultParsers = parsers{
		reflect.String:  parseString,
		reflect.Int:     parseInt,
		reflect.Int8:    parseInt8,
		reflect.Int16:   parseInt16,
		reflect.Int32:   parseInt32,
		reflect.Int64:   parseInt64,
		reflect.Uint:    parseUint,
		reflect.Uint8:   parseUint8,
		reflect.Uint16:  parseUint16,
		reflect.Uint32:  parseUint32,
		reflect.Uint64:  parseUint64,
		reflect.Float32: parseFloat32,
		reflect.Float64: parseFloat64,
		reflect.Bool:    parseBool,
		reflect.Pointer: parsePointer,
	}
}
