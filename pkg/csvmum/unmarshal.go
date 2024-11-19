package csvmum

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(data [][]string, v any) error {
	if len(data) == 0 {
		return nil
	}

	p := reflect.ValueOf(v)
	if p.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	e := p.Elem()
	t := e.Type()

	ftoi := map[string]int{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.IsExported() {
			name := f.Name
			if tag, ok := f.Tag.Lookup("csv"); ok {
				tags := strings.Split(tag, ",")
				if len(tags) > 0 {
					if t := tags[0]; t != "" {
						name = tags[0]
					}
				}
				if len(tags) > 1 {
					switch tags[1] {
					}
				}
			}
			ftoi[name] = i
		}
	}

	fmt.Printf("ftoi: %v\n", ftoi)

	headers := data[0]
	if len(headers) == 0 {
		return fmt.Errorf("no headers")
	}

	hm := map[int]int{}
	for i, h := range headers {
		if j, ok := ftoi[h]; ok {
			hm[i] = j
		}
	}

	fmt.Printf("hm: %v\n", hm)

	for _, record := range data[1:] {
		if len(record) == 0 {
			continue
		}
		if len(record) != len(headers) {
			return fmt.Errorf("record length mismatch")
		}
		for i, j := range hm {
			f := e.Field(j)
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
		fmt.Printf("v: %v\n", v)
	}

	return nil
}
