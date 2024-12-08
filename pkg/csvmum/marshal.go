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
	r := make([]string, 0, len(m.fieldList))

	for _, i := range m.fieldList {
		f := v.Field(i)
		if m, ok := f.Interface().(encoding.TextMarshaler); ok {
			b, err := m.MarshalText()
			if err != nil {
				return fmt.Errorf("cannot marshal: %w", err)
			}
			r = append(r, string(b))
			continue
		}

		switch f.Kind() {
		case reflect.String:
			r = append(r, fmt.Sprintf("%s", f.String()))
		case reflect.Int:
			r = append(r, fmt.Sprintf("%d", f.Int()))
		case reflect.Bool:
			r = append(r, fmt.Sprintf("%t", f.Bool()))
		case reflect.Float64:
			r = append(r, fmt.Sprintf("%s", strconv.FormatFloat(f.Float(), 'f', -1, 64)))
		}
	}

	err := m.writer.Write(r)
	if err != nil {
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
