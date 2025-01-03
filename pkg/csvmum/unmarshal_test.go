package csvmum

import (
	"bytes"
	"encoding/csv"
	"io"
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TextNewUnmarshalerEmptyFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testType struct {
		First  string
		Second int
	}

	b := &bytes.Buffer{}
	c := csv.NewReader(b)
	m, err := NewCSVUnmarshaler[testType](c)

	assert.NotNil(m)
	assert.Equal(io.EOF, err)
}

func TestNewUnmarshalerHeaders(t *testing.T) {
	t.Parallel()

	type testType struct {
		First  string
		Second int
		third  bool   // unexported
		Fourth string `csv:"fourth"`
	}

	assert := assert.New(t)

	b := &bytes.Buffer{}
	b.WriteString("First,fourth,nope,Second\n")

	m, err := NewUnmarshaler[testType](b)

	assert.Nil(err)
	assert.Equal([]int{0, 3, -1, 1}, m.fieldList)
}

func TestNewMarshalerNotAStruct(t *testing.T) {

	t.Parallel()

	type testType int

	assert := assert.New(t)

	b := &bytes.Buffer{}
	m, err := NewUnmarshaler[testType](b)

	assert.NotNil(m)
	assert.ErrorAs(err, ptr[*UnmarshalError](nil))
}

func TestNewUnmarshalerInvalidDelimiter(t *testing.T) {
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
	m, err := NewCSVUnmarshaler[testType](c)

	assert.NotNil(m)
	assert.ErrorAs(err, ptr[*UnmarshalError](nil))
}

func TestUnmarshalClosedReader(t *testing.T) {
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

	assert.ErrorAs(err, ptr[*UnmarshalError](nil))
}

func TestUnmarshalEndOfFile(t *testing.T) {
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
}

func TestUnmarshalParseError(t *testing.T) {
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

	assert.ErrorAs(err, ptr[*UnmarshalError](nil))
	assert.ErrorAs(err, ptr[*strconv.NumError](nil))
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testType struct {
		First   string
		Second  int
		Third   bool
		Fourth  float64
		fifth   string
		Sixth   bool `csv:"-"`
		Seventh *int `csv:"seventh"`
		Eighth  *int
	}

	b := &bytes.Buffer{}
	b.WriteString("First,extra,Second,Fourth,Third,seventh,Eighth\none,,2,3.14,true,,7\n")

	m, _ := NewUnmarshaler[testType](b)

	var record testType
	err := m.Unmarshal(&record)

	assert.Nil(err)
	assert.Equal(testType{
		First:   "one",
		Second:  2,
		Third:   true,
		Fourth:  3.14,
		Seventh: nil,
		Eighth:  ptr(7),
	}, record)
}
