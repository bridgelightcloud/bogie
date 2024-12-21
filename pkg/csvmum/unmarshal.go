package csvmum

import (
	"encoding"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
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
	fm, err := buildFieldMap(reflect.TypeOf(t))
	if err != nil {
		return um, fmt.Errorf("cannot unmarshal: %w", err)
	}

	hh, err := um.reader.Read()
	if err == io.EOF {
		return um, err
	}
	if err != nil {
		return um, fmt.Errorf("cannot unmarshal: %w", err)
	}

	um.fieldList = make([]int, len(hh))
	for i, h := range hh {
		if j, ok := fm[h]; ok {
			um.fieldList[i] = j
		} else {
			um.fieldList[i] = -1
		}
	}

	return um, nil
}

func (um *CSVUnmarshaler[T]) Unmarshal(record *T) error {
	r, err := um.reader.Read()
	if err == io.EOF {
		return err
	}
	if err != nil {
		return fmt.Errorf("cannot unmarshal: %w", err)
	}

	typ := reflect.TypeOf(*record)
	n := reflect.New(typ).Elem()

	for i, j := range um.fieldList {
		if j == -1 {
			continue
		}

		f := n.Field(j)

		if m, ok := f.Addr().Interface().(encoding.TextUnmarshaler); ok {
			if err := m.UnmarshalText([]byte(r[i])); err != nil {
				return fmt.Errorf("cannot unmarshal column %d, field %d: %w", i, j, err)
			}
			continue
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(r[i])
		case reflect.Int:
			i, err := strconv.ParseInt(r[i], 10, 64)
			if err != nil {
				return fmt.Errorf("cannot unmarshal column %d, field %d: error parsing int: %w", i, j, err)
			}
			f.SetInt(i)
		case reflect.Float64:
			f64, err := strconv.ParseFloat(r[i], 64)
			if err != nil {
				return fmt.Errorf("cannot unmarshal column %d, field %d: error parsing float64: %w", i, j, err)
			}
			f.SetFloat(f64)
		case reflect.Bool:
			b, err := strconv.ParseBool(r[i])
			if err != nil {
				return fmt.Errorf("cannot unmarshal column %d, field %d: error parsing bool: %w", i, j, err)
			}
			f.SetBool(b)
		case reflect.Pointer:
			switch f.Type().Elem().Kind() {
			case reflect.Int:
				if r[i] != "" {
					i, err := strconv.Atoi(r[i])
					if err != nil {
						return fmt.Errorf("cannot unmarshal column %d, field %d: error parsing int pointer: %w", i, j, err)
					}
					f.Set(reflect.ValueOf(&i))
				}
			case reflect.Float64:
				if r[i] != "" {
					f64, err := strconv.ParseFloat(r[i], 64)
					if err != nil {
						return fmt.Errorf("cannot unmarshal column %d, field %d: error parsing float64 pointer: %w", i, j, err)
					}
					f.Set(reflect.ValueOf(&f64))
				}
			}
		}
	}

	reflect.ValueOf(record).Elem().Set(n)

	return nil
}
