package gtfs

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestColor(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value         string
		expectedErr   error
		expectedColor Color
	}{{
		value:         "000000",
		expectedErr:   nil,
		expectedColor: Color("000000"),
	}, {
		value:         "FFFFFF",
		expectedErr:   nil,
		expectedColor: Color("FFFFFF"),
	}, {
		value:         "123456",
		expectedErr:   nil,
		expectedColor: Color("123456"),
	}, {
		value:         "ABCDEF",
		expectedErr:   nil,
		expectedColor: Color("ABCDEF"),
	}, {
		value:         "abc123",
		expectedErr:   nil,
		expectedColor: Color("ABC123"),
	}, {
		value:         "abC14D",
		expectedErr:   nil,
		expectedColor: Color("ABC14D"),
	}, {
		value:         "1234567",
		expectedErr:   fmt.Errorf("invalid color: 1234567"),
		expectedColor: Color(""),
	}, {
		value:         "ABCDEF1",
		expectedErr:   fmt.Errorf("invalid color: ABCDEF1"),
		expectedColor: Color(""),
	}, {
		value:         "12345",
		expectedErr:   fmt.Errorf("invalid color: 12345"),
		expectedColor: Color(""),
	}, {
		value:         "ABCDE",
		expectedErr:   fmt.Errorf("invalid color: ABCDE"),
		expectedColor: Color(""),
	}, {
		value:         "12345G",
		expectedErr:   fmt.Errorf("invalid color: 12345G"),
		expectedColor: Color(""),
	}, {
		value:         "ABCDEG",
		expectedErr:   fmt.Errorf("invalid color: ABCDEG"),
		expectedColor: Color(""),
	}, {
		value:         "",
		expectedErr:   fmt.Errorf("invalid color: "),
		expectedColor: Color(""),
	}, {
		value:         " 04FE2B",
		expectedErr:   nil,
		expectedColor: Color("04FE2B"),
	}, {
		value:         "#A5FF32",
		expectedErr:   fmt.Errorf("invalid color: #A5FF32"),
		expectedColor: Color(""),
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var c Color
			err := c.parseColor(tc.value)

			assert.Equal(tc.expectedErr, err)
			assert.Equal(tc.expectedColor, c)
		})
	}
}

func TestCurrencyCode(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value        string
		expectedErr  error
		expectedCode currencyCode
	}{{
		value:        "USD",
		expectedErr:  nil,
		expectedCode: currencyCode("USD"),
	}, {
		value:        "usd",
		expectedErr:  nil,
		expectedCode: currencyCode("USD"),
	}, {
		value:        "uSd",
		expectedErr:  nil,
		expectedCode: currencyCode("USD"),
	}, {
		value:        "usd ",
		expectedErr:  nil,
		expectedCode: currencyCode("USD"),
	}, {
		value:        "USD1",
		expectedErr:  fmt.Errorf("invalid currency code: %s", "USD1"),
		expectedCode: currencyCode(""),
	}, {
		value:        " ",
		expectedErr:  fmt.Errorf("invalid currency code: %s", " "),
		expectedCode: currencyCode(""),
	}, {
		value:        "",
		expectedErr:  fmt.Errorf("invalid currency code: %s", ""),
		expectedCode: currencyCode(""),
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var c currencyCode
			err := c.parseCurrencyCode(tc.value)

			if tc.value == "USD1" {
				fmt.Println(c)
			}

			assert.Equal(tc.expectedErr, err)
			assert.Equal(tc.expectedCode, c)
		})
	}
}

func TestDate(t *testing.T) {
	t.Parallel()

	ct := Time(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	zt := Time(time.Time{})
	tt := []struct {
		value   string
		expErr  error
		expTime Time
	}{{
		value:   "20060102",
		expErr:  nil,
		expTime: ct,
	}, {
		value:   "2006-01-02",
		expErr:  fmt.Errorf("invalid date format: %s", "2006-01-02"),
		expTime: zt,
	}, {
		value:   "2006/01/02",
		expErr:  fmt.Errorf("invalid date format: %s", "2006/01/02"),
		expTime: zt,
	}, {
		value:   "20060102 ",
		expErr:  nil,
		expTime: ct,
	}, {
		value:   " 20060102",
		expErr:  nil,
		expTime: ct,
	}, {
		value:   "20060002",
		expErr:  fmt.Errorf("invalid date value: %s", "20060002"),
		expTime: zt,
	}, {
		value:   " ",
		expErr:  fmt.Errorf("invalid date format: %s", " "),
		expTime: zt,
	}, {
		value:   "",
		expErr:  fmt.Errorf("invalid date format: %s", ""),
		expTime: zt,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var d Time
			err := d.parse(tc.value)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expTime, d)
		})
	}
}

func TestEnum(t *testing.T) {
	t.Parallel()

	ze := Enum(0)
	tt := []struct {
		value   string
		u       int
		expErr  error
		expEnum Enum
	}{{
		value:   "-1",
		u:       Availability,
		expErr:  fmt.Errorf("enum out of bounds: %d", -1),
		expEnum: ze,
	}, {
		value:   "0",
		u:       Availability,
		expErr:  nil,
		expEnum: Available,
	}, {
		value:   "1",
		u:       Availability,
		expErr:  nil,
		expEnum: Unavailable,
	}, {
		value:   "2",
		u:       Availability,
		expErr:  fmt.Errorf("enum out of bounds: %d", 2),
		expEnum: ze,
	}, {
		value:   "0",
		u:       Accessibility,
		expErr:  nil,
		expEnum: UnknownAccessibility,
	}, {
		value:   "1",
		u:       Accessibility,
		expErr:  nil,
		expEnum: AccessibeForAtLeastOne,
	}, {
		value:   "2",
		u:       Accessibility,
		expErr:  nil,
		expEnum: NotAccessible,
	}, {
		value:   "3",
		u:       Accessibility,
		expErr:  fmt.Errorf("enum out of bounds: %d", 3),
		expEnum: ze,
	}, {
		value:   "0",
		u:       ContinuousPickup,
		expErr:  nil,
		expEnum: RegularlyScheduled,
	}, {
		value:   "1",
		u:       ContinuousPickup,
		expErr:  nil,
		expEnum: NoneAvailable,
	}, {
		value:   "2",
		u:       ContinuousPickup,
		expErr:  nil,
		expEnum: MustPhoneAgency,
	}, {
		value:   "3",
		u:       ContinuousPickup,
		expErr:  nil,
		expEnum: MustCoordinate,
	}, {
		value:   "4",
		u:       ContinuousPickup,
		expErr:  fmt.Errorf("enum out of bounds: %d", 4),
		expEnum: ze,
	}, {

		value:   "0",
		u:       Timepoint,
		expErr:  nil,
		expEnum: ApproximateTime,
	}, {
		value:   "1",
		u:       Timepoint,
		expErr:  nil,
		expEnum: ExactTime,
	}, {
		value:   "2",
		u:       Timepoint,
		expErr:  fmt.Errorf("enum out of bounds: %d", 2),
		expEnum: ze,
	}, {
		value:   "",
		u:       Timepoint,
		expErr:  fmt.Errorf("invalid enum value: %s", ""),
		expEnum: ze,
	}, {
		value:   " ",
		u:       Timepoint,
		expErr:  fmt.Errorf("invalid enum value: %s", " "),
		expEnum: ze,
	}, {
		value:   "a",
		u:       Timepoint,
		expErr:  fmt.Errorf("invalid enum value: %s", "a"),
		expEnum: ze,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var e Enum
			err := e.Parse(tc.value, tc.u)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expEnum, e)
		})
	}
}

func TestInt(t *testing.T) {
	t.Parallel()

	zi := Int(0)
	tt := []struct {
		value  string
		expErr error
		expInt Int
	}{{
		value:  "-1",
		expErr: nil,
		expInt: Int(-1),
	}, {
		value:  "0",
		expErr: nil,
		expInt: Int(0),
	}, {
		value:  "1",
		expErr: nil,
		expInt: Int(1),
	}, {
		value:  "2",
		expErr: nil,
		expInt: Int(2),
	}, {
		value:  "a",
		expErr: fmt.Errorf("invalid integer value: %s", "a"),
		expInt: zi,
	}, {
		value:  "",
		expErr: fmt.Errorf("invalid integer value: %s", ""),
		expInt: zi,
	}, {
		value:  " ",
		expErr: fmt.Errorf("invalid integer value: %s", " "),
		expInt: zi,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var i Int
			err := i.Parse(tc.value)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expInt, i)
		})
	}
}

func TestFloat64(t *testing.T) {
	t.Parallel()

	zf := Float64(0)
	tt := []struct {
		value  string
		expErr error
		expFlt Float64
	}{{
		value:  "-1",
		expErr: nil,
		expFlt: Float64(-1),
	}, {
		value:  "0",
		expErr: nil,
		expFlt: Float64(0),
	}, {
		value:  "1",
		expErr: nil,
		expFlt: Float64(1),
	}, {
		value:  "2",
		expErr: nil,
		expFlt: Float64(2),
	}, {
		value:  "1.5",
		expErr: nil,
		expFlt: Float64(1.5),
	}, {
		value:  "1.5 ",
		expErr: nil,
		expFlt: Float64(1.5),
	}, {
		value:  " 1.5",
		expErr: nil,
		expFlt: Float64(1.5),
	}, {
		value:  "1.5.5",
		expErr: fmt.Errorf("invalid float value: %s", "1.5.5"),
		expFlt: zf,
	}, {
		value:  "1.5a",
		expErr: fmt.Errorf("invalid float value: %s", "1.5a"),
		expFlt: zf,
	}, {
		value:  "1.",
		expErr: nil,
		expFlt: Float64(1),
	}, {
		value:  "a",
		expErr: fmt.Errorf("invalid float value: %s", "a"),
		expFlt: zf,
	}, {
		value:  "",
		expErr: fmt.Errorf("invalid float value: %s", ""),
		expFlt: zf,
	}, {
		value:  " ",
		expErr: fmt.Errorf("invalid float value: %s", " "),
		expFlt: zf,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var f Float64
			err := f.Parse(tc.value)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expFlt, f)
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
