package csvmum

import (
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

	len := s.Len()
	if len == 0 {
		return out, nil
	}

	typ := reflect.ValueOf(s.Index(0).Interface()).Type()
	hm, err := getHeaderNamesToIndices(typ)
	if err != nil {
		return out, err
	}

	hs := []string{}
	for n := range hm {
		hs = append(hs, n)
	}

	out = append(out, hs)

	for i := 0; i < len; i++ {
		item := s.Index(i)
		record := []string{}
		for _, i := range hm {
			field := item.Field(i)
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
