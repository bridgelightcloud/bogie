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
			One   string `csv:"uno"`
			Two   int    `csv:"dos"`
			Three bool   `csv:""`
			Four  string
		}{},
		expected: map[string]int{"uno": 0, "dos": 1, "Four": 3},
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
