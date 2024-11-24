package csvmum

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

type fieldData struct {
	name     string
	exported bool
	idx      int

	timeLayout string
}

var builtin = regexp.MustCompile(`^(\w{1,2}):(.+)$`)

var (
	timeLayout = "tl"
)

func getHeaderData(t reflect.Type) (map[string]fieldData, error) {
	headers := map[string]fieldData{}

	if t.Kind() != reflect.Struct {
		return headers, fmt.Errorf("cannot get headers: not a struct")
	}

	for i := range t.NumField() {
		f := t.Field(i)
		fd := getFieldData(f, i)
		if !fd.exported {
			continue
		}
		headers[fd.name] = fd
	}

	return headers, nil
}

func getFieldData(f reflect.StructField, idx int) fieldData {
	fd := fieldData{
		name:       f.Name,
		idx:        idx,
		timeLayout: time.RFC3339,
	}

	if f.IsExported() {
		fd.exported = true

		if tag, ok := f.Tag.Lookup("csv"); ok {
			tags := strings.Split(tag, ",")
			for i, tag := range tags {
				if i == 0 {
					// name is always first
					switch tag {
					case "-":
						// ignore the field
						fd.exported = false
					case "":
						// use the field name as the name
					default:
						// use the tag as the name
						fd.name = tag
					}
				} else {
					// all other tags can be in any order
					b := builtin.FindStringSubmatch(tag)
					if len(b) == 3 {
						switch b[1] {
						case timeLayout:
							fd.timeLayout = b[2]
						default:
						}
					}
				}
			}
		}
	}
	return fd
}

func getOrderedHeaders(hm map[string]fieldData) []string {
	hh := make([]string, 0, len(hm))
	for n := range hm {
		hh = append(hh, n)
	}

	sort.SliceStable(hh, func(i, j int) bool {
		return hm[hh[i]].idx < hm[hh[j]].idx
	})

	return hh
}
