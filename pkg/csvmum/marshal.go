package csvmum

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

type CSVMarshaler[T any] struct {
	writer    *csv.Writer
	fieldList []int
}

func NewMarshaler[T any](w io.Writer) (*CSVMarshaler[T], error) {
	c := csv.NewWriter(w)

	return NewCSVMarshaler[T](c)
}

func NewCSVMarshaler[T any](w *csv.Writer) (*CSVMarshaler[T], error) {
	m := &CSVMarshaler[T]{writer: w}

	var t T
	hm, err := buildFieldMap(reflect.TypeOf(t))
	if err != nil {
		return m, fmt.Errorf("cannot marshal: %w", err)
	}

	hh, fl := getOrderedHeaders(hm)
	if err = m.writer.Write(hh); err != nil {
		return m, fmt.Errorf("cannot marshal: %w", err)
	}

	m.fieldList = fl

	return m, nil
}

func (m *CSVMarshaler[T]) Marshal(record T) error {
	v := reflect.ValueOf(record)
	row := []string{}

	for _, fi := range m.fieldList {
		field := v.Field(fi)

		f, err := marshalValue(field)
		if err != nil {
			return fmt.Errorf("cannot marshal field %d: %w", fi, err)
		}
		row = append(row, f)
	}

	if err := m.writer.Write(row); err != nil {
		return fmt.Errorf("cannot marshal: %w", err)
	}
	return nil
}

func (m *CSVMarshaler[T]) Flush() error {
	m.writer.Flush()

	if err := m.writer.Error(); err != nil {
		return fmt.Errorf("cannot marshal: %w", err)
	}

	return nil
}
