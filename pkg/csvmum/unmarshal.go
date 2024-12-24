package csvmum

import (
	"encoding"
	"encoding/csv"
	"fmt"
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
		return um, fmt.Errorf("cannot unmarshal: %w", err)
	}

	headers, err := um.reader.Read()
	if err == io.EOF {
		return um, err
	}
	if err != nil {
		return um, fmt.Errorf("cannot unmarshal: %w", err)
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
		return fmt.Errorf("cannot unmarshal: %w", err)
	}

	nsv := newSettableValue(record)

	for i, j := range um.fieldList {
		if j == -1 {
			continue
		}

		f := nsv.Field(j)

		if m, ok := f.Addr().Interface().(encoding.TextUnmarshaler); ok {
			if err := m.UnmarshalText([]byte(line[i])); err != nil {
				return fmt.Errorf("cannot unmarshal column %d, field %d: %w", i, j, err)
			}
			continue
		}

		v, err := parseValue(f, line[i])
		if err != nil {
			return fmt.Errorf("cannot unmarshal column %d, field %d: %w", i, j, err)
		}

		f.Set(v)
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
