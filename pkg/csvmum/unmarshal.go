package csvmum

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

func Unmarshal(data [][]string, v any) error {
	if len(data) == 0 {
		return nil
	}

	p := reflect.ValueOf(v)
	if p.Kind() != reflect.Ptr {
		return fmt.Errorf("cannot unmarshal: not a pointer")
	}

	pe := p.Elem()
	if pe.Kind() != reflect.Slice {
		return fmt.Errorf("cannot unmarshal: not a pointer to a slice")
	}

	typ := pe.Type().Elem()
	ftoi, err := getHeaderNamesToIndices(typ)
	if err != nil {
		return fmt.Errorf("cannot unmarshal: %v", err)
	}

	headers := data[0]
	if len(headers) == 0 {
		return fmt.Errorf("cannot unmarshal: no headers")
	}

	hm := map[int]int{}
	for i, h := range headers {
		if j, ok := ftoi[h]; ok {
			hm[i] = j
		}
	}
	if len(hm) == 0 {
		return fmt.Errorf("cannot unmarshal: no headers matched")
	}

	for _, record := range data[1:] {
		if len(record) == 0 {
			continue
		}
		if len(record) != len(headers) {
			continue
		}

		n := reflect.New(typ).Elem()
		for i, j := range hm {
			f := n.Field(j)

			if m, ok := f.Addr().Interface().(encoding.TextUnmarshaler); ok {
				err := m.UnmarshalText([]byte(record[i]))
				if err != nil {
					return fmt.Errorf("cannot unmarshal: %v", err)
				}
				continue
			}

			switch f.Kind() {
			case reflect.String:
				f.SetString(record[i])
			case reflect.Int:
				i, err := strconv.ParseInt(record[i], 10, 64)
				if err != nil {
					fmt.Printf("error parsing int: %v\n", err)
				}
				f.SetInt(i)
			case reflect.Bool:
				b, err := strconv.ParseBool(record[i])
				if err != nil {
					fmt.Printf("error parsing bool: %v\n", err)
				}
				f.SetBool(b)
			case reflect.Float64:
				f64, err := strconv.ParseFloat(record[i], 64)
				if err != nil {
					fmt.Printf("error parsing float64: %v\n", err)
				}
				f.SetFloat(f64)
			}
		}
		pe.Set(reflect.Append(pe, n))
	}

	return nil
}
