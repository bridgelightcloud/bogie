package csvmum

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
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
	hd, err := getHeaderData(typ)
	if err != nil {
		return out, err
	}

	hs := getOrderedHeaders(hd)

	out = append(out, hs)

	for i := range s.Len() {
		item := s.Index(i)
		record := []string{}
		for _, n := range hs {
			field := item.Field(hd[n].idx)
			switch field.Kind() {
			case reflect.String:
				record = append(record, fmt.Sprintf("%s", field.String()))
			case reflect.Int:
				record = append(record, fmt.Sprintf("%d", field.Int()))
			case reflect.Bool:
				record = append(record, fmt.Sprintf("%t", field.Bool()))
			case reflect.Float64:
				record = append(record, fmt.Sprintf("%s", strconv.FormatFloat(field.Float(), 'f', -1, 64)))
			case reflect.Struct:
				switch field.Type() {
				case timeType:
					record = append(record, fmt.Sprintf("%s", field.Interface().(time.Time).Format(hd[n].timeLayout)))
				default:
				}
			}
		}
		out = append(out, record)
	}

	return out, nil
}
