package csvmum

import (
	"bytes"
	"encoding/csv"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestNewCSVMarshaler(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &bytes.Buffer{}
		c := csv.NewWriter(b)
		m, err := NewCSVMarshaler[testType](c)

		assert.NotNil(m)
		assert.Nil(err)

		m.Flush()
		assert.Equal([]byte("First,Second\n"), b.Bytes())
	})

	t.Run("T is not a struct", func(t *testing.T) {
		t.Parallel()

		type testType int

		assert := assert.New(t)

		b := &bytes.Buffer{}
		c := csv.NewWriter(b)
		m, err := NewCSVMarshaler[testType](c)

		assert.NotNil(m)
		assert.EqualError(err, "cannot marshal: cannot get headers: not a struct")

		assert.Equal([]byte(nil), b.Bytes())
	})

	t.Run("invalid delimiter", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		b := &bytes.Buffer{}
		c := csv.NewWriter(b)
		c.Comma = utf8.RuneError

		m, err := NewCSVMarshaler[struct{}](c)

		assert.NotNil(m)
		assert.EqualError(err, "cannot marshal: csv: invalid field or comment delimiter")

	})

	t.Run("short write", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		var bw shortWriter
		c := csv.NewWriter(bw)
		m, err := NewCSVMarshaler[struct{}](c)

		assert.NotNil(m)

		err = m.Flush()
		assert.EqualError(err, "cannot marshal: short write")
	})
}

func TestNewMarshaler(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	type testType struct {
		First  string
		Second int
	}

	b := &bytes.Buffer{}
	m, err := NewMarshaler[testType](b)

	assert.NotNil(m)
	assert.Nil(err)

	m.Flush()
	assert.Equal([]byte("First,Second\n"), b.Bytes())
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		type testType struct{}

		b := &bytes.Buffer{}
		m, _ := NewMarshaler[testType](b)

		err := m.Marshal(testType{})
		assert.Nil(err)

		m.Flush()

		assert.Equal([]byte("\n\n"), b.Bytes())
	})

	t.Run("simple", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &bytes.Buffer{}
		m, _ := NewMarshaler[testType](b)

		err := m.Marshal(testType{"one", 1})
		assert.Nil(err)

		m.Flush()
		assert.Equal([]byte("First,Second\none,1\n"), b.Bytes())
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
		m, _ := NewMarshaler[testType](b)

		err := m.Marshal(testType{"one", 1, true, 3.14})
		assert.Nil(err)

		m.Flush()
		assert.Equal([]byte("First,Second,Third,Fourth\none,1,true,3.14\n"), b.Bytes())
	})

	t.Run("unexported", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		type testType struct {
			First  string
			second int
			Third  float64
		}

		b := &bytes.Buffer{}
		m, _ := NewMarshaler[testType](b)

		err := m.Marshal(testType{"one", 2, 67.3})
		assert.Nil(err)

		m.Flush()
		assert.Equal([]byte("First,Third\none,67.3\n"), b.Bytes())
	})

	t.Run("text marshaler", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		type testType struct {
			Field customMarshalAndUnmarshal
		}

		b := &bytes.Buffer{}
		m, _ := NewMarshaler[testType](b)

		err := m.Marshal(testType{Field: customMarshalAndUnmarshal{One: "one"}})
		assert.Nil(err)

		m.Flush()
		assert.Equal([]byte("Field\n~one~\n"), b.Bytes())
	})

	t.Run("invalid text marshaler", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		type testType struct {
			Field customMarshalAndUnmarshal
		}

		b := &bytes.Buffer{}
		m, _ := NewMarshaler[testType](b)

		err := m.Marshal(testType{Field: customMarshalAndUnmarshal{One: ""}})
		assert.EqualError(err, "cannot marshal: invalid text: ")

		m.Flush()
		assert.Equal([]byte("Field\n"), b.Bytes())
	})

	t.Run("closed writer", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &closeReaderWriter{}
		m, _ := NewMarshaler[testType](b)

		b.Close()

		m.Marshal(testType{"one", 1})

		err := m.Flush()
		assert.EqualError(err, "cannot marshal: closed")
	})

	t.Run("invalid delimiter", func(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)

		type testType struct {
			First  string
			Second int
		}

		b := &bytes.Buffer{}
		m, _ := NewMarshaler[testType](b)

		m.writer.Comma = utf8.RuneError

		err := m.Marshal(testType{"one", 1})
		assert.EqualError(err, "cannot marshal: csv: invalid field or comment delimiter")

		m.Flush()
		assert.Equal([]byte("First,Second\n"), b.Bytes())
	})
}
