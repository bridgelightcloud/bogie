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
