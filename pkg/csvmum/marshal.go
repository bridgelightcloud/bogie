package csvmum

import (
	"encoding"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type CSVMarshaler[T any] struct {
	writer    csv.Writer
	fieldList []int
}

func NewMarshaler[T any](w io.Writer) (*CSVMarshaler[T], error) {
	c := csv.NewWriter(w)

	return NewCSVMarshaler[T](*c)
}

func NewCSVMarshaler[T any](w csv.Writer) (*CSVMarshaler[T], error) {
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
	m.writer.Flush()

	m.fieldList = fl

	if err = m.writer.Error(); err != nil {
		return m, fmt.Errorf("cannot marshal: %w", err)
	}

	return m, nil
}

func (m *CSVMarshaler[T]) Marshal(record T) error {
	v := reflect.ValueOf(record)
	row := make([]string, len(m.fieldList))

	for ci, fi := range m.fieldList {
		f := v.Field(fi)
		if cm, ok := f.Interface().(encoding.TextMarshaler); ok {
			b, err := cm.MarshalText()
			if err != nil {
				return fmt.Errorf("cannot marshal: %w", err)
			}
			row[ci] = string(b)
			continue
		}

		switch f.Kind() {
		case reflect.String:
			row[ci] = fmt.Sprintf("%s", f.String())
		case reflect.Int:
			row[ci] = fmt.Sprintf("%d", f.Int())
		case reflect.Bool:
			row[ci] = fmt.Sprintf("%t", f.Bool())
		case reflect.Float64:
			row[ci] = fmt.Sprintf("%s", strconv.FormatFloat(f.Float(), 'f', -1, 64))
		}
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
