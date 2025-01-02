package csvmum

import (
	"fmt"
	"reflect"
)

type UnsupportedTypeError struct {
	Kind reflect.Kind
}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported type: %s", e.Kind)
}

type UnmarshalValueError struct {
	Err error
}

func (e *UnmarshalValueError) Error() string {
	return e.Err.Error()
}

func (e *UnmarshalValueError) Unwrap() error {
	return e.Err
}
