package csvmum

import (
	"encoding"
	"math/bits"
	"reflect"
	"strconv"
)

func unmarshalString(value string, field reflect.Value) error {
	field.SetString(reflect.ValueOf(value).String())

	return nil
}

func unmarshalIntFactory(size int) unmarshaler {
	return func(value string, field reflect.Value) error {
		i, err := strconv.ParseInt(value, 10, size)

		if err != nil {
			return err
		}

		field.SetInt(i)

		return nil
	}
}

func unmarshalUintFactory(size int) unmarshaler {
	return func(value string, field reflect.Value) error {
		i, err := strconv.ParseUint(value, 10, size)

		if err != nil {
			return err
		}

		field.SetUint(i)

		return nil
	}
}

func unmarshalFloatFactory(size int) unmarshaler {
	return func(value string, field reflect.Value) error {
		f, err := strconv.ParseFloat(value, size)

		if err != nil {
			return err
		}

		field.SetFloat(f)

		return nil
	}
}

func unmarshalBool(value string, field reflect.Value) error {
	b, err := strconv.ParseBool(value)

	if err == nil {
		field.SetBool(b)
	}

	return err
}

func unmarshalPointer(value string, field reflect.Value) error {
	if value == "" {
		return nil
	}

	if field.IsNil() {
		field.Set(reflect.New(field.Type().Elem()))
	}

	err := unmarshalValue(value, field.Elem())

	return err
}

func unmarshalValue(value string, field reflect.Value) error {
	if m, ok := field.Addr().Interface().(encoding.TextUnmarshaler); ok {
		if err := m.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		return nil
	}

	kind := field.Kind()
	if p, ok := defaultUnmarshalers[kind]; ok {
		err := p(value, field)
		if err != nil {
			return err
		}
		return nil
	}

	return &UnsupportedTypeError{kind}
}

type unmarshaler func(string, reflect.Value) error

var defaultUnmarshalers map[reflect.Kind]unmarshaler

func init() {
	defaultUnmarshalers = map[reflect.Kind]unmarshaler{
		reflect.String:  unmarshalString,
		reflect.Int:     unmarshalIntFactory(bits.UintSize),
		reflect.Int8:    unmarshalIntFactory(8),
		reflect.Int16:   unmarshalIntFactory(16),
		reflect.Int32:   unmarshalIntFactory(32),
		reflect.Int64:   unmarshalIntFactory(64),
		reflect.Uint:    unmarshalUintFactory(bits.UintSize),
		reflect.Uint8:   unmarshalUintFactory(8),
		reflect.Uint16:  unmarshalUintFactory(16),
		reflect.Uint32:  unmarshalUintFactory(32),
		reflect.Uint64:  unmarshalUintFactory(64),
		reflect.Float32: unmarshalFloatFactory(32),
		reflect.Float64: unmarshalFloatFactory(64),
		reflect.Bool:    unmarshalBool,
		reflect.Pointer: unmarshalPointer,
	}
}
