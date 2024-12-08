package csvmum

import (
	"bytes"
	"encoding/csv"
	"io"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestNewCSVUnmarshaler(t *testing.T) {
	t.Parallel()

	t.Run("empty file", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &bytes.Buffer{}
		c := csv.NewReader(b)
		m, err := NewCSVUnmarshaler[testType](*c)

		assert.NotNil(m)
		assert.Equal(io.EOF, err)
	})

	t.Run("simple", func(t *testing.T) {
		t.Parallel()

		type testType struct {
			First  string
			Second int
		}

		assert := assert.New(t)

		b := &bytes.Buffer{}
		b.WriteString("First,Second\n")

		c := csv.NewReader(b)
		m, err := NewCSVUnmarshaler[testType](*c)

		assert.Nil(err)
		assert.Equal(map[int]int{0: 0, 1: 1}, m.headerMap)
		assert.Equal(map[string]int{"First": 0, "Second": 1}, m.fieldMap)
	})

	t.Run("T is not a struct", func(t *testing.T) {
		t.Parallel()

		type testType int

		assert := assert.New(t)

		b := &bytes.Buffer{}
		c := csv.NewReader(b)
		m, err := NewCSVUnmarshaler[testType](*c)

		assert.NotNil(m)
		assert.EqualError(err, "cannot unmarshal: cannot get headers: not a struct")
	})

	t.Run("invalid delimiter", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &bytes.Buffer{}
		b.WriteString("First,Second\n")

		c := csv.NewReader(b)
		c.Comma = utf8.RuneError

		m, err := NewCSVUnmarshaler[testType](*c)

		assert.NotNil(m)
		assert.EqualError(err, "cannot unmarshal: csv: invalid field or comment delimiter")
	})
}

func TestNewUnmarshaler(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testType struct {
		First  string
		Second int
	}

	b := &bytes.Buffer{}
	b.WriteString("First,Second\none,2\n")

	m, err := NewUnmarshaler[testType](b)

	assert.NotNil(m)
	assert.Nil(err)

	assert.Equal(map[int]int{0: 0, 1: 1}, m.headerMap)
	assert.Equal(map[string]int{"First": 0, "Second": 1}, m.fieldMap)
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	t.Run("end of file", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &bytes.Buffer{}
		b.WriteString("First,Second\n")

		m, _ := NewUnmarshaler[testType](b)

		var record testType
		err := m.Unmarshal(&record)

		assert.Equal(io.EOF, err)
	})

	t.Run("simple", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &bytes.Buffer{}
		b.WriteString("First,Second\none,2\n")

		m, _ := NewUnmarshaler[testType](b)

		var record testType
		err := m.Unmarshal(&record)

		assert.Nil(err)
		assert.Equal(testType{First: "one", Second: 2}, record)
	})

	t.Run("invalid record: int", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			Int int
		}

		b := &bytes.Buffer{}
		b.WriteString("Int\none\n")

		m, _ := NewUnmarshaler[testType](b)

		var record testType
		err := m.Unmarshal(&record)

		assert.EqualError(err, "cannot unmarshal: error parsing int: strconv.ParseInt: parsing \"one\": invalid syntax")
	})

	t.Run("invalid record: bool", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			Bool bool
		}

		b := &bytes.Buffer{}
		b.WriteString("Bool\nblah\n")

		m, _ := NewUnmarshaler[testType](b)

		var record testType
		err := m.Unmarshal(&record)

		assert.EqualError(err, "cannot unmarshal: error parsing bool: strconv.ParseBool: parsing \"blah\": invalid syntax")
	})

	t.Run("invalid record: float64", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			Float64 float64
		}

		b := &bytes.Buffer{}
		b.WriteString("Float64\nblah\n")

		m, _ := NewUnmarshaler[testType](b)

		var record testType
		err := m.Unmarshal(&record)

		assert.EqualError(err, "cannot unmarshal: error parsing float64: strconv.ParseFloat: parsing \"blah\": invalid syntax")
	})

	t.Run("complex", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
			Third  bool
			Fourth float64
		}

		b := &bytes.Buffer{}
		b.WriteString("First,Second,Third,Fourth\none,2,true,3.14\n")

		m, _ := NewUnmarshaler[testType](b)

		var record testType
		err := m.Unmarshal(&record)

		assert.Nil(err)
		assert.Equal(testType{First: "one", Second: 2, Third: true, Fourth: 3.14}, record)
	})

	t.Run("custom unmarshaler", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			Field customMarshalAndUnmarshal
		}

		b := &bytes.Buffer{}
		b.WriteString("Field\n~one~\n")

		m, _ := NewUnmarshaler[testType](b)

		var record testType
		err := m.Unmarshal(&record)

		assert.Nil(err)
		assert.Equal(testType{Field: customMarshalAndUnmarshal{One: "one"}}, record)
	})

	t.Run("invalid custom unmarshaler", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			Field customMarshalAndUnmarshal
		}

		b := &bytes.Buffer{}
		b.WriteString("Field\n~\n")

		m, _ := NewUnmarshaler[testType](b)

		var record testType
		err := m.Unmarshal(&record)

		assert.EqualError(err, "cannot unmarshal: invalid text: ~")
	})

	t.Run("closed reader", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &closeReaderWriter{}
		m, _ := NewUnmarshaler[testType](b)

		b.Close()

		var record testType
		err := m.Unmarshal(&record)

		assert.EqualError(err, "cannot unmarshal: closed")
	})

}

// func TestUnmarshal(t *testing.T) {
// 	t.Parallel()

// 	tt := []struct {
// 		name string
// 		data [][]string
// 		v    any
// 		res  any
// 		err  error
// 	}{{
// 		name: "empty data",
// 		data: [][]string{},
// 		v:    &[]struct{ A int }{},
// 		res:  &[]struct{ A int }{},
// 		err:  nil,
// 	}, {
// 		name: "not a pointer",
// 		data: [][]string{{"A"}, {"1"}},
// 		v:    []struct{}{},
// 		res:  []struct{}{},
// 		err:  fmt.Errorf("cannot unmarshal: not a pointer"),
// 	}, {
// 		name: "not a pointer to a slice",
// 		data: [][]string{{"A"}, {"1"}},
// 		v:    &struct{}{},
// 		res:  &struct{}{},
// 		err:  fmt.Errorf("cannot unmarshal: not a pointer to a slice"),
// 	}, {
// 		name: "not a pointer to a slice of structs",
// 		data: [][]string{{"A"}, {"1"}},
// 		v:    &[]int{},
// 		res:  &[]int{},
// 		err:  fmt.Errorf("cannot unmarshal: cannot get headers: not a struct"),
// 	}, {
// 		name: "no headers",
// 		data: [][]string{{}, {"1"}},
// 		v:    &[]struct{}{},
// 		res:  &[]struct{}{},
// 		err:  fmt.Errorf("cannot unmarshal: no headers"),
// 	}, {
// 		name: "no headers matched",
// 		data: [][]string{{"A"}, {"1"}},
// 		v:    &[]struct{ B int }{},
// 		res:  &[]struct{ B int }{},
// 		err:  fmt.Errorf("cannot unmarshal: no headers matched"),
// 	}, {
// 		name: "no records",
// 		data: [][]string{{"A"}},
// 		v:    &[]struct{ A int }{},
// 		res:  &[]struct{ A int }{},
// 		err:  nil,
// 	}, {
// 		name: "record has different number of fields",
// 		data: [][]string{{"A"}, {"1", "two"}},
// 		v:    &[]struct{ A int }{},
// 		res:  &[]struct{ A int }{},
// 		err:  nil,
// 	}, {
// 		name: "simple",
// 		data: [][]string{{"A"}, {"1"}},
// 		v:    &[]struct{ A int }{},
// 		res:  &[]struct{ A int }{{A: 1}},
// 		err:  nil,
// 	}, {
// 		name: "complex",
// 		data: [][]string{{"A", "B", "C"}, {"1", "2", "true"}},
// 		v: &[]struct {
// 			A int
// 			B int
// 			C bool
// 		}{},
// 		res: &[]struct {
// 			A int
// 			B int
// 			C bool
// 		}{{A: 1, B: 2, C: true}},
// 		err: nil,
// 	}, {
// 		name: "unexported",
// 		data: [][]string{{"A", "C"}, {"1", "67.3"}},
// 		v: &[]struct {
// 			A int
// 			b string
// 			C float64
// 		}{},
// 		res: &[]struct {
// 			A int
// 			b string
// 			C float64
// 		}{{A: 1, C: 67.3}},
// 		err: nil,
// 	}, {
// 		name: "mixed",
// 		data: [][]string{{"A", "B", "C", "D"}, {"test", "1", "true", "3.14"}},
// 		v: &[]struct {
// 			A string
// 			B int
// 			C bool
// 			D float64
// 		}{},
// 		res: &[]struct {
// 			A string
// 			B int
// 			C bool
// 			D float64
// 		}{{A: "test", B: 1, C: true, D: 3.14}},
// 		err: nil,
// 	}, {
// 		name: "empty records",
// 		data: [][]string{{"A", "B"}, {}, {"test", "1"}},
// 		v: &[]struct {
// 			A string
// 			B int
// 		}{},
// 		res: &[]struct {
// 			A string
// 			B int
// 		}{{A: "test", B: 1}},
// 		err: nil,
// 	}, {
// 		name: "invalid int",
// 		data: [][]string{{"A"}, {"one"}},
// 		v:    &[]struct{ A int }{},
// 		res:  &[]struct{ A int }{{A: 0}},
// 		err:  nil,
// 	}, {
// 		name: "invalid bool",
// 		data: [][]string{{"A"}, {"yes"}},
// 		v:    &[]struct{ A bool }{},
// 		res:  &[]struct{ A bool }{{A: false}},
// 		err:  nil,
// 	}, {
// 		name: "invalid float64",
// 		data: [][]string{{"A"}, {"three point one four"}},
// 		v:    &[]struct{ A float64 }{},
// 		res:  &[]struct{ A float64 }{{A: 0}},
// 		err:  nil,
// 	}}

// 	for _, tc := range tt {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			assert := assert.New(t)

// 			err := Unmarshal(tc.data, tc.v)
// 			assert.Equal(tc.res, tc.v)
// 			assert.Equal(tc.err, err)
// 		})
// 	}
// }

// type customUnmarshal struct {
// 	One string
// }

// func (c *customUnmarshal) UnmarshalText(text []byte) error {
// 	if len(text) < 2 {
// 		return fmt.Errorf("invalid text: %s", string(text))
// 	}
// 	c.One = string(text[1 : len(text)-1])
// 	return nil
// }

// func TestUnmarshalTextMarshaler(t *testing.T) {
// 	t.Parallel()

// 	type cs struct {
// 		Field customUnmarshal
// 	}

// 	tt := []struct {
// 		name     string
// 		input    [][]string
// 		expected []cs
// 		err      error
// 	}{{
// 		name:     "simple",
// 		input:    [][]string{{"Field"}, {"~one~"}},
// 		expected: []cs{{Field: customUnmarshal{One: "one"}}},
// 		err:      nil,
// 	}, {
// 		name:     "invalid text",
// 		input:    [][]string{{"Field"}, {"~"}},
// 		expected: []cs{},
// 		err:      fmt.Errorf("cannot unmarshal: invalid text: ~"),
// 	}}

// 	for _, tc := range tt {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()

// 			assert := assert.New(t)

// 			out := []cs{}
// 			err := Unmarshal(tc.input, &out)
// 			assert.Equal(tc.expected, out)
// 			assert.Equal(tc.err, err)
// 		})
// 	}
// }
