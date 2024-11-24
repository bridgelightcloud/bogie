package csvmum

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHeaderData(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    any
		expected map[string]fieldData
		err      error
	}{{
		name:     "empty",
		input:    struct{}{},
		expected: map[string]fieldData{},
		err:      nil,
	}, {
		name:     "not a struct",
		input:    []int{1, 2, 3},
		expected: map[string]fieldData{},
		err:      fmt.Errorf("cannot get headers: not a struct"),
	}, {
		name:  "simple",
		input: struct{ One string }{},
		expected: map[string]fieldData{
			"One": {name: "One", exported: true, idx: 0, timeLayout: time.RFC3339},
		},
		err: nil,
	}, {
		name: "complex",
		input: struct {
			One   string
			Two   int
			Three bool
		}{},
		expected: map[string]fieldData{
			"One":   {name: "One", exported: true, idx: 0, timeLayout: time.RFC3339},
			"Two":   {name: "Two", exported: true, idx: 1, timeLayout: time.RFC3339},
			"Three": {name: "Three", exported: true, idx: 2, timeLayout: time.RFC3339}},
		err: nil,
	}, {
		name: "unexported",
		input: struct {
			One   string
			two   int
			Three bool
		}{},
		expected: map[string]fieldData{
			"One":   {name: "One", exported: true, idx: 0, timeLayout: time.RFC3339},
			"Three": {name: "Three", exported: true, idx: 2, timeLayout: time.RFC3339},
		},
		err: nil,
	}, {
		name: "tagged",
		input: struct {
			One string `csv:"uno"`
			Two int    `csv:"dos"`
		}{},
		expected: map[string]fieldData{
			"uno": {name: "uno", exported: true, idx: 0, timeLayout: time.RFC3339},
			"dos": {name: "dos", exported: true, idx: 1, timeLayout: time.RFC3339},
		},
		err: nil,
	}, {
		name: "tagged but not exported",
		input: struct {
			One string `csv:"uno"`
			two int    `csv:"dos"`
		}{},
		expected: map[string]fieldData{"uno": {name: "uno", exported: true, idx: 0, timeLayout: time.RFC3339}},
		err:      nil,
	}, {
		name: "tagged with hyphen -",
		input: struct {
			One string `csv:"-"`
			Two int    `csv:"dos"`
		}{},
		expected: map[string]fieldData{"dos": {name: "dos", exported: true, idx: 1, timeLayout: time.RFC3339}},
		err:      nil,
	}}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			headers, err := getHeaderData(reflect.TypeOf(tc.input))
			assert.Equal(tc.expected, headers)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetFieldData(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    any
		expected fieldData
	}{{
		name:     "unexported",
		input:    struct{ one string }{},
		expected: fieldData{name: "one", exported: false, idx: 0, timeLayout: time.RFC3339},
	}, {
		name:     "exported",
		input:    struct{ One string }{},
		expected: fieldData{name: "One", exported: true, idx: 0, timeLayout: time.RFC3339},
	}, {
		name: "tagged",
		input: struct {
			One string `csv:"uno"`
		}{},
		expected: fieldData{name: "uno", exported: true, idx: 0, timeLayout: time.RFC3339},
	}, {
		name: "tagged with date format",
		input: struct {
			One string `csv:"uno,tl:20060102"`
		}{},
		expected: fieldData{name: "uno", exported: true, idx: 0, timeLayout: "20060102"},
	}, {
		name: "tagged empty",
		input: struct {
			One string `csv:""`
		}{},
		expected: fieldData{name: "One", exported: true, idx: 0, timeLayout: time.RFC3339},
	}, {
		name: "tagged -",
		input: struct {
			One string `csv:"-"`
		}{},
		expected: fieldData{name: "One", exported: false, idx: 0, timeLayout: time.RFC3339},
	}}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			f := reflect.TypeOf(tc.input).Field(0)

			fd := getFieldData(f, 0)
			assert.Equal(tc.expected, fd)
		})
	}
}

func TestGetOrderedHeaders(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    map[string]fieldData
		expected []string
	}{{
		name:     "empty",
		input:    map[string]fieldData{},
		expected: []string{},
	}, {
		name:     "simple",
		input:    map[string]fieldData{"One": {name: "One"}},
		expected: []string{"One"},
	}, {
		name: "complex",
		input: map[string]fieldData{
			"One":   {idx: 0},
			"Two":   {idx: 1},
			"Three": {idx: 2},
		},
		expected: []string{"One", "Two", "Three"},
	}, {
		name: "unordered",
		input: map[string]fieldData{
			"Two":   {idx: 2},
			"One":   {idx: 0},
			"Three": {idx: 7},
		},
		expected: []string{"One", "Two", "Three"},
	}}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			headers := getOrderedHeaders(tc.input)
			assert.Equal(tc.expected, headers)
		})
	}
}
