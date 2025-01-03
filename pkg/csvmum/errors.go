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

type UnmarshalError struct {
	Err         error
	FieldIndex  *int
	ColumnIndex *int
}

func (e UnmarshalError) indices() string {
	var out string
	if e.FieldIndex != nil && e.ColumnIndex != nil {
		out = fmt.Sprintf(" (field index: %d, column index: %d)", *e.FieldIndex, *e.ColumnIndex)
	}
	return out
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("csv unmarshal: %s%s", e.Err.Error(), e.indices())
}

func (e *UnmarshalError) Unwrap() error {
	return e.Err
}

type MarshalError struct {
	Err error
}

func (e *MarshalError) Error() string {
	return e.Err.Error()
}

func (e *MarshalError) Unwrap() error {
	return e.Err
}
