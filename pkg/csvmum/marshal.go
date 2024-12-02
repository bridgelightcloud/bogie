package csvmum

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

func Marshal(v any) ([][]string, error) {
	out := [][]string{}

	s := reflect.ValueOf(v)
	if s.Kind() != reflect.Slice {
		return out, fmt.Errorf("cannot marshal: not a slice")
	}

	if s.Len() == 0 {
		return out, nil
	}

	typ := reflect.ValueOf(s.Index(0).Interface()).Type()
	hm, err := getHeaderNamesToIndices(typ)
	if err != nil {
		return out, err
	}

	hs := getOrderedHeaders(hm)

	out = append(out, hs)

	for i := range s.Len() {
		item := s.Index(i)
		record := []string{}
		for _, n := range hs {
			field := item.Field(hm[n])

			if m, ok := field.Interface().(encoding.TextMarshaler); ok {
				b, err := m.MarshalText()
				if err != nil {
					return out, err
				}
				record = append(record, string(b))
				continue
			}

			switch field.Kind() {
			case reflect.String:
				record = append(record, fmt.Sprintf("%s", field.String()))
			case reflect.Int:
				record = append(record, fmt.Sprintf("%d", field.Int()))
			case reflect.Bool:
				record = append(record, fmt.Sprintf("%t", field.Bool()))
			case reflect.Float64:
				record = append(record, fmt.Sprintf("%s", strconv.FormatFloat(field.Float(), 'f', -1, 64)))
			}
		}
		out = append(out, record)
	}

	return out, nil
}
