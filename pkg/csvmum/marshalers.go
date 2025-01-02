package csvmum

import (
	"encoding"
	"reflect"
	"strconv"
)

func marshalString(v reflect.Value) (string, error) {
	return v.String(), nil
}

func marshalInt(v reflect.Value) (string, error) {
	return strconv.FormatInt(v.Int(), 10), nil
}

func marshalUint(v reflect.Value) (string, error) {
	return strconv.FormatUint(v.Uint(), 10), nil
}

func marshalFloat32(v reflect.Value) (string, error) {
	return strconv.FormatFloat(v.Float(), 'f', -1, 32), nil
}

func marshalFloat64(v reflect.Value) (string, error) {
	return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
}

func marshalBool(v reflect.Value) (string, error) {
	return strconv.FormatBool(v.Bool()), nil
}

func marshalPointer(v reflect.Value) (string, error) {
	if v.IsNil() {
		return "", nil
	}

	return marshalValue(v.Elem())
}

func marshalTextMarshaler(tm encoding.TextMarshaler) (string, error) {
	b, err := tm.MarshalText()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func marshalValue(v reflect.Value) (string, error) {
	if tm, ok := v.Interface().(encoding.TextMarshaler); ok {
		return marshalTextMarshaler(tm)
	}

	if marshaler, ok := defaultMarshalers[v.Kind()]; ok {
		return marshaler(v)
	}

	return "", &UnsupportedTypeError{v.Kind()}
}

type marshaler func(reflect.Value) (string, error)

var defaultMarshalers map[reflect.Kind]marshaler

func init() {
	defaultMarshalers = map[reflect.Kind]marshaler{
		reflect.String:  marshalString,
		reflect.Int:     marshalInt,
		reflect.Int8:    marshalInt,
		reflect.Int16:   marshalInt,
		reflect.Int32:   marshalInt,
		reflect.Int64:   marshalInt,
		reflect.Uint:    marshalUint,
		reflect.Uint8:   marshalUint,
		reflect.Uint16:  marshalUint,
		reflect.Uint32:  marshalUint,
		reflect.Uint64:  marshalUint,
		reflect.Float32: marshalFloat32,
		reflect.Float64: marshalFloat64,
		reflect.Bool:    marshalBool,
		reflect.Pointer: marshalPointer,
	}
}
