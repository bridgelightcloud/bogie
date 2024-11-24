package csvmum

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		data [][]string
		v    any
		res  any
		err  error
	}{{
		name: "empty data",
		data: [][]string{},
		v:    &[]struct{ A int }{},
		res:  &[]struct{ A int }{},
		err:  nil,
	}, {
		name: "not a pointer",
		data: [][]string{{"A"}, {"1"}},
		v:    []struct{}{},
		res:  []struct{}{},
		err:  fmt.Errorf("cannot unmarshal: not a pointer"),
	}, {
		name: "not a pointer to a slice",
		data: [][]string{{"A"}, {"1"}},
		v:    &struct{}{},
		res:  &struct{}{},
		err:  fmt.Errorf("cannot unmarshal: not a pointer to a slice"),
	}, {
		name: "not a pointer to a slice of structs",
		data: [][]string{{"A"}, {"1"}},
		v:    &[]int{},
		res:  &[]int{},
		err:  fmt.Errorf("cannot unmarshal: cannot get headers: not a struct"),
	}, {
		name: "no headers",
		data: [][]string{{}, {"1"}},
		v:    &[]struct{}{},
		res:  &[]struct{}{},
		err:  fmt.Errorf("cannot unmarshal: no headers"),
	}, {
		name: "no headers matched",
		data: [][]string{{"A"}, {"1"}},
		v:    &[]struct{ B int }{},
		res:  &[]struct{ B int }{},
		err:  fmt.Errorf("cannot unmarshal: no headers matched"),
	}, {
		name: "no records",
		data: [][]string{{"A"}},
		v:    &[]struct{ A int }{},
		res:  &[]struct{ A int }{},
		err:  nil,
	}, {
		name: "record has different number of fields",
		data: [][]string{{"A"}, {"1", "two"}},
		v:    &[]struct{ A int }{},
		res:  &[]struct{ A int }{},
		err:  nil,
	}, {
		name: "simple",
		data: [][]string{{"A"}, {"1"}},
		v:    &[]struct{ A int }{},
		res:  &[]struct{ A int }{{A: 1}},
		err:  nil,
	}, {
		name: "complex",
		data: [][]string{{"A", "B", "C"}, {"1", "2", "true"}},
		v: &[]struct {
			A int
			B int
			C bool
		}{},
		res: &[]struct {
			A int
			B int
			C bool
		}{{A: 1, B: 2, C: true}},
		err: nil,
	}, {
		name: "unexported",
		data: [][]string{{"A", "C"}, {"1", "67.3"}},
		v: &[]struct {
			A int
			b string
			C float64
		}{},
		res: &[]struct {
			A int
			b string
			C float64
		}{{A: 1, C: 67.3}},
		err: nil,
	}, {
		name: "mixed",
		data: [][]string{{"A", "B", "C", "D"}, {"test", "1", "true", "3.14"}},
		v: &[]struct {
			A string
			B int
			C bool
			D float64
		}{},
		res: &[]struct {
			A string
			B int
			C bool
			D float64
		}{{A: "test", B: 1, C: true, D: 3.14}},
		err: nil,
	}, {
		name: "empty records",
		data: [][]string{{"A", "B"}, {}, {"test", "1"}},
		v: &[]struct {
			A string
			B int
		}{},
		res: &[]struct {
			A string
			B int
		}{{A: "test", B: 1}},
		err: nil,
	}, {
		name: "invalid int",
		data: [][]string{{"A"}, {"one"}},
		v:    &[]struct{ A int }{},
		res:  &[]struct{ A int }{{A: 0}},
		err:  nil,
	}, {
		name: "invalid bool",
		data: [][]string{{"A"}, {"yes"}},
		v:    &[]struct{ A bool }{},
		res:  &[]struct{ A bool }{{A: false}},
		err:  nil,
	}, {
		name: "invalid float64",
		data: [][]string{{"A"}, {"three point one four"}},
		v:    &[]struct{ A float64 }{},
		res:  &[]struct{ A float64 }{{A: 0}},
		err:  nil,
	}, {
		name: "invalid time",
		data: [][]string{{"A"}, {"2024-14-35"}},
		v:    &[]struct{ A time.Time }{},
		res:  &[]struct{ A time.Time }{{A: time.Time{}}},
		err:  nil,
	}}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			err := Unmarshal(tc.data, tc.v)
			assert.Equal(tc.res, tc.v)
			assert.Equal(tc.err, err)
		})
	}
}
