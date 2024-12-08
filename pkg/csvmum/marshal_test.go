package csvmum

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    any
		expected [][]string
		err      error
	}{{
		name:     "empty",
		input:    []struct{}{},
		expected: [][]string{},
		err:      nil,
	}, {
		name:     "not a struct",
		input:    []int{1, 2, 3},
		expected: [][]string{},
		err:      fmt.Errorf("cannot marshal: cannot get headers: not a struct"),
	}, {
		name:     "not a slice",
		input:    34,
		expected: [][]string{},
		err:      fmt.Errorf("cannot marshal: not a slice"),
	}, {
		name:     "simple",
		input:    []struct{ One string }{{One: "one"}},
		expected: [][]string{{"One"}, {"one"}},
		err:      nil,
	}, {
		name: "complex",
		input: []struct {
			One   string
			Two   int
			Three bool
		}{{One: "one", Two: 1, Three: true}},
		expected: [][]string{{"One", "Two", "Three"}, {"one", "1", "true"}},
		err:      nil,
	}, {
		name: "unexported",
		input: []struct {
			One   string
			two   int
			Three float64
		}{{One: "one", two: 2, Three: 67.3}},
		expected: [][]string{{"One", "Three"}, {"one", "67.3"}},
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

type customMarshal struct {
	One string
}

func (c customMarshal) MarshalText() ([]byte, error) {
	if len(c.One) == 0 {
		return nil, fmt.Errorf("invalid text: %s", c.One)
	}

	return []byte(fmt.Sprintf("~%s~", c.One)), nil
}

func TestMarshalTextMarshaler(t *testing.T) {
	t.Parallel()

	type cs struct {
		Field customMarshal
	}

	tt := []struct {
		name     string
		input    []cs
		expected [][]string
		err      error
	}{{
		name:     "simple",
		input:    []cs{{Field: customMarshal{One: "one"}}},
		expected: [][]string{{"Field"}, {"~one~"}},
		err:      nil,
	}, {
		name:     "invalid",
		input:    []cs{{Field: customMarshal{One: ""}}},
		expected: [][]string{{"Field"}},
		err:      fmt.Errorf("cannot marshal: invalid text: "),
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
