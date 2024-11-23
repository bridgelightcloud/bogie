package csvmum

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHeaderNamesToIndices(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    any
		expected map[string]int
		err      error
	}{{
		name:     "empty",
		input:    struct{}{},
		expected: map[string]int{},
		err:      nil,
	}, {
		name:     "not a struct",
		input:    []int{1, 2, 3},
		expected: map[string]int{},
		err:      fmt.Errorf("cannot get headers: not a struct"),
	}, {
		name:     "simple",
		input:    struct{ One string }{},
		expected: map[string]int{"One": 0},
		err:      nil,
	}, {
		name: "complex",
		input: struct {
			One   string
			Two   int
			Three bool
		}{},
		expected: map[string]int{"One": 0, "Two": 1, "Three": 2},
		err:      nil,
	}, {
		name: "unexported",
		input: struct {
			One   string
			two   int
			Three bool
		}{},
		expected: map[string]int{"One": 0, "Three": 2},
		err:      nil,
	}, {
		name: "tagged",
		input: struct {
			One string `csv:"uno"`
			Two int    `csv:"dos"`
		}{},
		expected: map[string]int{"uno": 0, "dos": 1},
		err:      nil,
	}, {
		name: "tagged but not exported",
		input: struct {
			One string `csv:"uno"`
			two int    `csv:"dos"`
		}{},
		expected: map[string]int{"uno": 0},
		err:      nil,
	}, {
		name: "tagged with hyphen -",
		input: struct {
			One string `csv:"-"`
			Two int    `csv:"dos"`
		}{},
		expected: map[string]int{"dos": 1},
		err:      nil,
	}}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			headers, err := getHeaderNamesToIndices(reflect.TypeOf(tc.input))
			assert.Equal(tc.expected, headers)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetExportedName(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    any
		expected string
	}{{
		name:     "unexported",
		input:    struct{ one string }{},
		expected: "-",
	}, {
		name:     "exported",
		input:    struct{ One string }{},
		expected: "One",
	}, {
		name: "tagged",
		input: struct {
			One string `csv:"uno"`
		}{},
		expected: "uno",
	}, {
		name: "tagged multiple",
		input: struct {
			One string `csv:"uno,dos"`
		}{},
		expected: "uno",
	}, {
		name: "tagged empty",
		input: struct {
			One string `csv:""`
		}{},
		expected: "One",
	}, {
		name: "tagged -",
		input: struct {
			One string `csv:"-"`
		}{},
		expected: "-",
	}}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			f := reflect.TypeOf(tc.input).Field(0)

			name := getExportedName(f)
			assert.Equal(tc.expected, name)
		})
	}
}

func TestGetOrderedHeaders(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    map[string]int
		expected []string
	}{{
		name:     "empty",
		input:    map[string]int{},
		expected: []string{},
	}, {
		name:     "simple",
		input:    map[string]int{"One": 0},
		expected: []string{"One"},
	}, {
		name:     "complex",
		input:    map[string]int{"One": 0, "Two": 1, "Three": 2},
		expected: []string{"One", "Two", "Three"},
	}, {
		name:     "unordered",
		input:    map[string]int{"Two": 2, "One": 0, "Three": 7},
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
