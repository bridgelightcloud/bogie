package csvmum

import (
	"fmt"
	"reflect"
	"strings"
)

func getHeaderNamesToIndices(t reflect.Type) (map[string]int, error) {
	headers := map[string]int{}

	if t.Kind() != reflect.Struct {
		return headers, fmt.Errorf("cannot get headers: not a struct")
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if name := getExportedName(f); name != "" {
			headers[name] = i
		}
	}
	return headers, nil
}

func getExportedName(f reflect.StructField) string {
	var name string
	if f.IsExported() {
		name = f.Name
		if tag, ok := f.Tag.Lookup("csv"); ok {
			tags := strings.Split(tag, ",")
			if len(tags) > 0 {
				name = tags[0]
			}
		}
	}
	return name
}
