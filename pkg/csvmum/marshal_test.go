package csvmum

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStructHeaders(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    any
		expected []string
		err      error
	}{{
		name:     "empty",
		input:    struct{}{},
		expected: []string{},
		err:      nil,
	}, {
		name:     "not a struct",
		input:    []int{1, 2, 3},
		expected: []string{},
		err:      fmt.Errorf("GetStructHeaders: not a struct"),
	}, {
		name:     "simple",
		input:    struct{ One string }{},
		expected: []string{"One"},
		err:      nil,
	}, {
		name: "complex",
		input: struct {
			One   string
			Two   int
			Three bool
		}{},
		expected: []string{"One", "Two", "Three"},
		err:      nil,
	}, {
		name: "unexported",
		input: struct {
			One   string
			two   int
			Three bool
		}{},
		expected: []string{"One", "Three"},
		err:      nil,
	}}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			headers, err := GetStructHeaders(tc.input)
			assert.Equal(tc.expected, headers)
			assert.Equal(tc.err, err)
		})
	}
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    any
		expected []byte
		err      error
	}{{
		name:     "empty",
		input:    []struct{}{},
		expected: []byte{},
		err:      nil,
	}, {
		name:     "not a struct",
		input:    []int{1, 2, 3},
		expected: []byte{},
		err:      fmt.Errorf("GetStructHeaders: not a struct"),
	}, {
		name:     "not a slice",
		input:    34,
		expected: []byte{},
		err:      fmt.Errorf("Marshal: not a slice"),
	}, {
		name:     "simple",
		input:    []struct{ One string }{{One: "one"}},
		expected: []byte("One\none\n"),
		err:      nil,
	}, {
		name: "complex",
		input: []struct {
			One   string
			Two   int
			Three bool
		}{{One: "one", Two: 1, Three: true}},
		expected: []byte("One,Two,Three\none,1,true\n"),
		err:      nil,
	}, {
		name: "unexported",
		input: []struct {
			One   string
			two   int
			Three float64
		}{{One: "one", two: 2, Three: 67.3}},
		expected: []byte("One,Three\none,67.3\n"),
		err:      nil,
	}}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			out, err := Marshal(tc.input)
			assert.Equal(tc.expected, out)
			assert.Equal(tc.err, err)
		})
	}
}
