package csvmum

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func getHeaderNamesToIndices(t reflect.Type) (map[string]int, error) {
	headers := map[string]int{}

	if t.Kind() != reflect.Struct {
		return headers, fmt.Errorf("cannot get headers: not a struct")
	}
	for i := range t.NumField() {
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

func getOrderedHeaders(hm map[string]int) []string {
	hh := make([]string, 0, len(hm))
	for n := range hm {
		hh = append(hh, n)
	}

	sort.SliceStable(hh, func(i, j int) bool {
		return hm[hh[i]] < hm[hh[j]]
	})

	return hh
}