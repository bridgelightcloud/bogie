package csvmum

import (
	"encoding/csv"
	"io"
	"reflect"
)

type CSVUnmarshaler[T any] struct {
	reader    *csv.Reader
	fieldList []int
}

func NewUnmarshaler[T any](r io.Reader) (*CSVUnmarshaler[T], error) {
	c := csv.NewReader(r)
	return NewCSVUnmarshaler[T](c)
}

func NewCSVUnmarshaler[T any](r *csv.Reader) (*CSVUnmarshaler[T], error) {
	um := &CSVUnmarshaler[T]{reader: r}

	var t T
	fields, err := buildFieldMap(reflect.TypeOf(t))
	if err != nil {
		return um, &UnmarshalError{Err: err}
	}

	headers, err := um.reader.Read()
	if err == io.EOF {
		return um, err
	}
	if err != nil {
		return um, &UnmarshalError{Err: err}
	}

	um.fieldList = make([]int, len(headers))
	for i, h := range headers {
		if j, ok := fields[h]; ok {
			um.fieldList[i] = j
		} else {
			um.fieldList[i] = -1
		}
	}

	return um, nil
}

func (um *CSVUnmarshaler[T]) Unmarshal(record *T) error {
	line, err := um.reader.Read()
	if err == io.EOF {
		return err
	}
	if err != nil {
		return &UnmarshalError{Err: err}
	}

	nsv := newSettableValue(record)

	for i, j := range um.fieldList {
		if j == -1 {
			continue
		}

		field := nsv.Field(j)

		if err := unmarshalValue(line[i], field); err != nil {
			return &UnmarshalError{
				Err:         err,
				FieldIndex:  &i,
				ColumnIndex: &j,
			}
		}
	}

	setRecordValue(record, nsv)

	return nil
}

func newSettableValue[T any](r *T) reflect.Value {
	typ := reflect.TypeOf(*r)
	return reflect.New(typ).Elem()
}

func setRecordValue[T any](record *T, v reflect.Value) {
	reflect.ValueOf(record).Elem().Set(v)
}
