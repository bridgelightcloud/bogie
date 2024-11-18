package csvmum

import (
	"fmt"
	"reflect"
	"strconv"
)

func GetStructHeaders(v any) ([]string, error) {
	headers := []string{}
	if reflect.ValueOf(v).Kind() == reflect.Struct {
		for i := 0; i < reflect.ValueOf(v).NumField(); i++ {
			if reflect.TypeOf(v).Field(i).IsExported() {
				headers = append(headers, fmt.Sprintf("%v", reflect.TypeOf(v).Field(i).Name))
			}
		}
		return headers, nil
	} else {
		return headers, fmt.Errorf("GetStructHeaders: not a struct")
	}
}

func Marshal(v any) ([][]string, error) {
	out := [][]string{}

	if val := reflect.ValueOf(v); val.Kind() == reflect.Slice {
		len := val.Len()
		if len == 0 {
			return out, nil
		}

		v := val.Index(0).Interface()
		headers, err := GetStructHeaders(v)
		if err != nil {
			return out, err
		}
		out = append(out, headers)

		for i := 0; i < len; i++ {
			item := val.Index(i)
			record := []string{}
			for i := 0; i < item.NumField(); i++ {
				field := item.Field(i)
				if !item.Type().Field(i).IsExported() {
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
	} else {
		return out, fmt.Errorf("Marshal: not a slice")
	}

	return out, nil
}
