package csvmum

import (
	"math/bits"
	"reflect"
	"strconv"
)

func parseString(value string, field reflect.Value) error {
	field.SetString(reflect.ValueOf(value).String())

	return nil
}

func parseIntFactory(size int) parser {
	return func(value string, field reflect.Value) error {
		i, err := strconv.ParseInt(value, 10, size)

		if err != nil {
			return err
		}

		field.SetInt(i)

		return nil
	}
}

func parseUintFactory(size int) parser {
	return func(value string, field reflect.Value) error {
		i, err := strconv.ParseUint(value, 10, size)

		if err != nil {
			return err
		}

		field.SetUint(i)

		return nil
	}
}

func parseFloatFactory(size int) parser {
	return func(value string, field reflect.Value) error {
		f, err := strconv.ParseFloat(value, size)

		if err != nil {
			return err
		}

		field.SetFloat(f)

		return nil
	}
}

func parseBool(value string, field reflect.Value) error {
	b, err := strconv.ParseBool(value)

	if err == nil {
		field.SetBool(b)
	}

	return err
}

func parsePointer(value string, field reflect.Value) error {
	if value == "" {
		return nil
	}

	if field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
	}

	err := parseValue(value, field.Elem())

	return err
}

func parseValue(value string, field reflect.Value) error {
	kind := field.Kind()
	if p, ok := defaultParsers[kind]; ok {
		err := p(value, field)
		if err != nil {
			return &ParseError{err}
		}
		return nil
	}

	return &UnsupportedTypeError{kind}
}

type parser func(string, reflect.Value) error

var defaultParsers map[reflect.Kind]parser

func init() {
	defaultParsers = map[reflect.Kind]parser{
		reflect.String:  parseString,
		reflect.Int:     parseIntFactory(bits.UintSize),
		reflect.Int8:    parseIntFactory(8),
		reflect.Int16:   parseIntFactory(16),
		reflect.Int32:   parseIntFactory(32),
		reflect.Int64:   parseIntFactory(64),
		reflect.Uint:    parseUintFactory(bits.UintSize),
		reflect.Uint8:   parseUintFactory(8),
		reflect.Uint16:  parseUintFactory(16),
		reflect.Uint32:  parseUintFactory(32),
		reflect.Uint64:  parseUintFactory(64),
		reflect.Float32: parseFloatFactory(32),
		reflect.Float64: parseFloatFactory(64),
		reflect.Bool:    parseBool,
		reflect.Pointer: parsePointer,
	}
}
