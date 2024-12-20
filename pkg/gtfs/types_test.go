package gtfs

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseColor(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value         string
		expectedErr   error
		expectedColor string
	}{{
		value:         "000000",
		expectedErr:   nil,
		expectedColor: "000000",
	}, {
		value:         "FFFFFF",
		expectedErr:   nil,
		expectedColor: "FFFFFF",
	}, {
		value:         "123456",
		expectedErr:   nil,
		expectedColor: "123456",
	}, {
		value:         "ABCDEF",
		expectedErr:   nil,
		expectedColor: "ABCDEF",
	}, {
		value:         "abc123",
		expectedErr:   nil,
		expectedColor: "ABC123",
	}, {
		value:         "abC14D",
		expectedErr:   nil,
		expectedColor: "ABC14D",
	}, {
		value:         "1234567",
		expectedErr:   fmt.Errorf("invalid color: 1234567"),
		expectedColor: "",
	}, {
		value:         "ABCDEF1",
		expectedErr:   fmt.Errorf("invalid color: ABCDEF1"),
		expectedColor: "",
	}, {
		value:         "12345",
		expectedErr:   fmt.Errorf("invalid color: 12345"),
		expectedColor: "",
	}, {
		value:         "ABCDE",
		expectedErr:   fmt.Errorf("invalid color: ABCDE"),
		expectedColor: "",
	}, {
		value:         "12345G",
		expectedErr:   fmt.Errorf("invalid color: 12345G"),
		expectedColor: "",
	}, {
		value:         "ABCDEG",
		expectedErr:   fmt.Errorf("invalid color: ABCDEG"),
		expectedColor: "",
	}, {
		value:         "",
		expectedErr:   fmt.Errorf("invalid color: "),
		expectedColor: "",
	}, {
		value:         " 04FE2B",
		expectedErr:   nil,
		expectedColor: "04FE2B",
	}, {
		value:         "#A5FF32",
		expectedErr:   fmt.Errorf("invalid color: #A5FF32"),
		expectedColor: "",
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var c string
			err := ParseColor(tc.value, &c)

			assert.Equal(tc.expectedErr, err)
			assert.Equal(tc.expectedColor, c)
		})
	}
}

func TestParseCurrencyCode(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value        string
		expectedErr  error
		expectedCode string
	}{{
		value:        "USD",
		expectedErr:  nil,
		expectedCode: "USD",
	}, {
		value:        "usd",
		expectedErr:  nil,
		expectedCode: "USD",
	}, {
		value:        "uSd",
		expectedErr:  nil,
		expectedCode: "USD",
	}, {
		value:        "usd ",
		expectedErr:  nil,
		expectedCode: "USD",
	}, {
		value:        "USD1",
		expectedErr:  fmt.Errorf("invalid currency code: %s", "USD1"),
		expectedCode: "",
	}, {
		value:        " ",
		expectedErr:  fmt.Errorf("invalid currency code: %s", " "),
		expectedCode: "",
	}, {
		value:        "",
		expectedErr:  fmt.Errorf("invalid currency code: %s", ""),
		expectedCode: "",
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var c string
			err := ParseCurrencyCode(tc.value, &c)

			assert.Equal(tc.expectedErr, err)
			assert.Equal(tc.expectedCode, c)
		})
	}
}

func TestDateMarshalText(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		date Date
		out  []byte
		err  error
	}{{
		name: "valid date",
		date: Date{Time: time.Date(2004, 11, 27, 0, 0, 0, 0, time.UTC)},
		out:  []byte("20041127"),
		err:  nil,
	}, {
		name: "zero date",
		date: Date{Time: time.Time{}},
		out:  []byte("00010101"),
		err:  nil,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			mt, err := tc.date.MarshalText()
			assert.Equal(tc.out, mt)
			assert.Equal(tc.err, err)

		})
	}
}

func TestDateUnmarshalText(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		in   []byte
		date Date
		err  error
	}{{
		name: "valid date",
		in:   []byte("20241127"),
		date: Date{Time: time.Date(2024, 11, 27, 0, 0, 0, 0, time.UTC)},
		err:  nil,
	}, {
		name: "zero date?",
		in:   []byte("00010101"),
		date: Date{Time: time.Time{}},
		err:  nil,
	}, {
		name: "invalid date",
		in:   []byte("Nov 27, 2024"),
		date: Date{Time: time.Time{}},
		err:  fmt.Errorf("invalid date value: Nov 27, 2024"),
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var d Date
			err := d.UnmarshalText(tc.in)

			assert.Equal(tc.date, d)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDateMarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		date Date
		out  []byte
		err  error
	}{{
		name: "valid date",
		date: Date{Time: time.Date(2024, 11, 27, 0, 0, 0, 0, time.UTC)},
		out:  []byte("1732665600"),
		err:  nil,
	}, {
		name: "zero date",
		date: Date{Time: time.Time{}},
		out:  []byte("-62135596800"),
		err:  nil,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			dm, err := tc.date.MarshalJSON()

			assert.Equal(tc.out, dm)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDateUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		in   []byte
		date Date
		err  error
	}{{
		name: "valid date",
		in:   []byte("1732665600"),
		date: Date{Time: time.Date(2024, 11, 27, 0, 0, 0, 0, time.Local)},
		err:  nil,
	}, {
		name: "zero date",
		in:   []byte("-62135596800"),
		date: Date{Time: time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)},
		err:  nil,
	}, {
		name: "invalid date",
		in:   []byte("x"),
		// date: Date{Time: time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)},
		date: Date{Time: time.Time{}},
		err:  fmt.Errorf("invalid date value: x"),
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var d Date
			err := d.UnmarshalJSON(tc.in)

			assert.Equal(tc.date, d)
			assert.Equal(tc.err, err)
		})
	}
}

func TestTimeMarshalText(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		time Time
		out  []byte
		err  error
	}{{
		name: "time under 24 hrs",
		time: Time{Time: time.Date(1, 1, 1, 12, 55, 30, 0, time.UTC)},
		out:  []byte("12:55:30"),
		err:  nil,
	}, {
		name: "time over 24 hrs",
		time: Time{Time: time.Date(1, 1, 2, 1, 34, 22, 0, time.UTC)},
		out:  []byte("25:34:22"),
		err:  nil,
	}, {
		name: "zero time",
		time: Time{Time: time.Time{}},
		out:  []byte("00:00:00"),
		err:  nil,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			mt, err := tc.time.MarshalText()
			assert.Equal(tc.out, mt)
			assert.Equal(tc.err, err)

		})
	}
}

func TestTimeUnmarshalText(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		in   []byte
		time Time
		err  error
	}{{
		name: "time under 24 hrs",
		in:   []byte("17:23:22"),
		time: Time{Time: time.Date(0, 1, 1, 17, 23, 22, 0, time.UTC)},
		err:  nil,
	}, {
		name: "time over 24 hrs",
		in:   []byte("25:34:22"),
		time: Time{Time: time.Date(0, 1, 2, 1, 34, 22, 0, time.UTC)},
	}, {
		name: "zero time",
		in:   []byte("00:00:00"),
		time: Time{Time: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)},
		// time: Time{Time: time.Time{}},
		err: nil,
	}, {
		name: "invalid time",
		in:   []byte("09:34 AM"),
		time: Time{Time: time.Time{}},
		err:  fmt.Errorf("invalid time value: 09:34 AM"),
	}, {
		name: "invalid time over 24 hrs",
		in:   []byte("24:77:22"),
		time: Time{Time: time.Time{}},
		err:  fmt.Errorf("invalid time value: 24:77:22"),
	}, {
		name: "invalid time under 48 hrs",
		in:   []byte("48:34:22"),
		time: Time{Time: time.Time{}},
		err:  fmt.Errorf("invalid time value: 48:34:22"),
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var time Time
			err := time.UnmarshalText(tc.in)

			assert.Equal(tc.time, time)
			assert.Equal(tc.err, err)
		})
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		time Time
		out  []byte
		err  error
	}{{
		name: "valid time",
		time: Time{Time: time.Date(1, 1, 1, 12, 57, 44, 0, time.UTC)},
		out:  []byte("-62135550136"),
		err:  nil,
	}, {
		name: "zero time",
		time: Time{Time: time.Time{}},
		out:  []byte("null"),
		err:  nil,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			dm, err := tc.time.MarshalJSON()

			assert.Equal(tc.out, dm)
			assert.Equal(tc.err, err)
		})
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		in   []byte
		time Time
		err  error
	}{{
		name: "valid time",
		in:   []byte("-62135550136"),
		time: Time{Time: time.Date(1, 1, 1, 12, 57, 44, 0, time.Local)}, err: nil,
	}, {
		name: "zero time",
		in:   []byte("null"),
		time: Time{Time: time.Time{}},
		err:  nil,
	}, {
		name: "invalid time",
		in:   []byte("x"),
		time: Time{Time: time.Time{}},
		err:  fmt.Errorf("invalid time value: x"),
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var d Time
			err := d.UnmarshalJSON(tc.in)

			assert.Equal(tc.time, d)
			assert.Equal(tc.err, err)
		})
	}
}

func TestErrorList(t *testing.T) {
	t.Parallel()

	tt := []struct {
		errList errorList
		err     error
		expList errorList
	}{{
		errList: errorList{},
		err:     nil,
		expList: errorList{},
	}, {
		errList: errorList{fmt.Errorf("error 1")},
		err:     nil,
		expList: errorList{fmt.Errorf("error 1")},
	}, {
		errList: errorList{},
		err:     fmt.Errorf("error 1"),
		expList: errorList{fmt.Errorf("error 1")},
	}, {
		errList: errorList{fmt.Errorf("error 1")},
		err:     fmt.Errorf("error 2"),
		expList: errorList{fmt.Errorf("error 1"), fmt.Errorf("error 2")},
	}, {
		errList: errorList{fmt.Errorf("error 1"), fmt.Errorf("error 2")},
		err:     fmt.Errorf("error 3"),
		expList: errorList{fmt.Errorf("error 1"), fmt.Errorf("error 2"), fmt.Errorf("error 3")},
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			err := tc.errList.add(tc.err)

			assert.Equal(tc.err, err)
			assert.Equal(tc.expList, tc.errList)
		})
	}
}
