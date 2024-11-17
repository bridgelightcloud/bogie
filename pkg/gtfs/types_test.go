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

func TestParseDate(t *testing.T) {
	t.Parallel()

	ct := time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC)
	zt := time.Time{}
	tt := []struct {
		value   string
		expErr  error
		expTime time.Time
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

			var d time.Time
			err := ParseDate(tc.value, &d)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expTime, d)
		})
	}
}

func TestParseTime(t *testing.T) {
	t.Parallel()

	ct := time.Date(0, time.January, 1, 15, 4, 5, 0, time.UTC)
	zt := time.Time{}
	tt := []struct {
		value   string
		expErr  error
		expTime time.Time
	}{{
		value:   "15:04:05",
		expErr:  nil,
		expTime: ct,
	}, {
		value:   "15:04:05 ",
		expErr:  nil,
		expTime: ct,
	}, {
		value:   " 15:04:05",
		expErr:  nil,
		expTime: ct,
	}, {
		value:   "15:04:05 ",
		expErr:  nil,
		expTime: ct,
	}, {
		value:   "15:04:05.000",
		expErr:  fmt.Errorf("invalid time format: %s", "15:04:05.000"),
		expTime: zt,
	}, {
		value:   "3:04:05",
		expErr:  nil,
		expTime: time.Date(0, time.January, 1, 3, 4, 5, 0, time.UTC),
	}, {
		value:   "03:4:05",
		expErr:  fmt.Errorf("invalid time format: %s", "03:4:05"),
		expTime: zt,
	}, {
		value:   "30:04:05",
		expErr:  fmt.Errorf("invalid time value: %s, parsing time \"%s\": hour out of range", "30:04:05", "30:04:05"),
		expTime: zt,
	}, {
		value:   "15:60:05",
		expErr:  fmt.Errorf("invalid time value: %s, parsing time \"%s\": minute out of range", "15:60:05", "15:60:05"),
		expTime: zt,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var d time.Time
			err := ParseTime(tc.value, &d)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expTime, d)
		})
	}
}

func TestParseEnum(t *testing.T) {
	t.Parallel()

	ze := 0
	tt := []struct {
		value   string
		u       enumBounds
		expErr  error
		expEnum int ``
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
		u:       WheelchairAccessible,
		expErr:  nil,
		expEnum: UnknownAccessibility,
	}, {
		value:   "1",
		u:       WheelchairAccessible,
		expErr:  nil,
		expEnum: AtLeastOneWheelchairAccomodated,
	}, {
		value:   "2",
		u:       WheelchairAccessible,
		expErr:  nil,
		expEnum: NoWheelchairsAccomodated,
	}, {
		value:   "3",
		u:       WheelchairAccessible,
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

			var e int
			err := ParseEnum(tc.value, tc.u, &e)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expEnum, e)
		})
	}
}

func TestParseInt(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value  string
		expErr error
		expInt int
	}{{
		value:  "-1",
		expErr: nil,
		expInt: -1,
	}, {
		value:  "0",
		expErr: nil,
		expInt: 0,
	}, {
		value:  "1",
		expErr: nil,
		expInt: 1,
	}, {
		value:  "2",
		expErr: nil,
		expInt: 2,
	}, {
		value:  "1.5",
		expErr: fmt.Errorf("invalid integer value: %s", "1.5"),
		expInt: 0,
	}, {
		value:  "1.",
		expErr: fmt.Errorf("invalid integer value: %s", "1."),
		expInt: 0,
	}, {
		value:  " 300",
		expErr: nil,
		expInt: 300,
	}, {
		value:  "300 ",
		expErr: nil,
		expInt: 300,
	}, {
		value:  "5a",
		expErr: fmt.Errorf("invalid integer value: %s", "5a"),
		expInt: 0,
	}, {
		value:  "a",
		expErr: fmt.Errorf("invalid integer value: %s", "a"),
		expInt: 0,
	}, {
		value:  "",
		expErr: fmt.Errorf("invalid integer value: %s", ""),
		expInt: 0,
	}, {
		value:  " ",
		expErr: fmt.Errorf("invalid integer value: %s", " "),
		expInt: 0,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var i int
			err := ParseInt(tc.value, &i)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expInt, i)
		})
	}
}

func TestParseFloat64(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value  string
		expErr error
		expFlt float64
	}{{
		value:  "-1",
		expErr: nil,
		expFlt: -1.0,
	}, {
		value:  "0",
		expErr: nil,
		expFlt: 0.0,
	}, {
		value:  "1",
		expErr: nil,
		expFlt: 1.0,
	}, {
		value:  "2",
		expErr: nil,
		expFlt: 2.0,
	}, {
		value:  "1.5",
		expErr: nil,
		expFlt: 1.5,
	}, {
		value:  "1.5 ",
		expErr: nil,
		expFlt: 1.5,
	}, {
		value:  " 1.5",
		expErr: nil,
		expFlt: 1.5,
	}, {
		value:  "1.5.5",
		expErr: fmt.Errorf("invalid float value: %s", "1.5.5"),
		expFlt: 0.0,
	}, {
		value:  "1.5a",
		expErr: fmt.Errorf("invalid float value: %s", "1.5a"),
		expFlt: 0.0,
	}, {
		value:  "1.",
		expErr: nil,
		expFlt: 1.0,
	}, {
		value:  "a",
		expErr: fmt.Errorf("invalid float value: %s", "a"),
		expFlt: 0.0,
	}, {
		value:  "",
		expErr: fmt.Errorf("invalid float value: %s", ""),
		expFlt: 0.0,
	}, {
		value:  " ",
		expErr: fmt.Errorf("invalid float value: %s", " "),
		expFlt: 0.0,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var f float64
			err := ParseFloat(tc.value, &f)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expFlt, f)
		})
	}
}

func TestParseString(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value  string
		expStr string
	}{{
		value:  "string",
		expStr: "string",
	}, {
		value:  " string",
		expStr: "string",
	}, {
		value:  "string ",
		expStr: "string",
	}, {
		value:  " string ",
		expStr: "string",
	}, {
		value:  " ",
		expStr: "",
	}, {
		value:  "",
		expStr: "",
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			var s string
			ParseString(tc.value, &s)

			assert.Equal(tc.expStr, s)
		})
	}
}

func TestParseLat(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value  string
		expErr error
		expLat float64
	}{{
		value:  "-1",
		expErr: nil,
		expLat: -1.0,
	}, {
		value:  "0",
		expErr: nil,
		expLat: 0.0,
	}, {
		value:  "1",
		expErr: nil,
		expLat: 1.0,
	}, {
		value:  "2",
		expErr: nil,
		expLat: 2.0,
	}, {
		value:  "1.5",
		expErr: nil,
		expLat: 1.5,
	}, {
		value:  "1.5 ",
		expErr: nil,
		expLat: 1.5,
	}, {
		value:  " 1.5",
		expErr: nil,
		expLat: 1.5,
	}, {
		value:  "1.5.5",
		expErr: fmt.Errorf("invalid latitude value: %s", "1.5.5"),
		expLat: 0.0,
	}, {
		value:  "1.5a",
		expErr: fmt.Errorf("invalid latitude value: %s", "1.5a"),
		expLat: 0.0,
	}, {
		value:  "1.",
		expErr: nil,
		expLat: 1.0,
	}, {
		value:  "a",
		expErr: fmt.Errorf("invalid latitude value: %s", "a"),
		expLat: 0.0,
	}, {
		value:  "",
		expErr: fmt.Errorf("invalid latitude value: %s", ""),
		expLat: 0.0,
	}, {
		value:  " ",
		expErr: fmt.Errorf("invalid latitude value: %s", " "),
		expLat: 0.0,
	}, {
		value:  "90",
		expErr: nil,
		expLat: 90.0,
	}, {
		value:  "-90",
		expErr: nil,
		expLat: -90.0,
	}, {
		value:  "90.1",
		expErr: fmt.Errorf("latitude out of bounds: %f", 90.1),
		expLat: 0.0,
	}, {
		value:  "-90.1",
		expErr: fmt.Errorf("latitude out of bounds: %f", -90.1),
		expLat: 0.0,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			c := Coords{}
			err := ParseLat(tc.value, &c)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expLat, c.Lat)
		})
	}
}

func TestParseLon(t *testing.T) {
	t.Parallel()

	tt := []struct {
		value  string
		expErr error
		expLon float64
	}{{
		value:  "-1",
		expErr: nil,
		expLon: -1.0,
	}, {
		value:  "0",
		expErr: nil,
		expLon: 0.0,
	}, {
		value:  "1",
		expErr: nil,
		expLon: 1.0,
	}, {
		value:  "2",
		expErr: nil,
		expLon: 2.0,
	}, {
		value:  "1.5",
		expErr: nil,
		expLon: 1.5,
	}, {
		value:  "1.5 ",
		expErr: nil,
		expLon: 1.5,
	}, {
		value:  " 1.5",
		expErr: nil,
		expLon: 1.5,
	}, {
		value:  "1.5.5",
		expErr: fmt.Errorf("invalid longitude value: %s", "1.5.5"),
		expLon: 0.0,
	}, {
		value:  "1.5a",
		expErr: fmt.Errorf("invalid longitude value: %s", "1.5a"),
		expLon: 0.0,
	}, {
		value:  "1.",
		expErr: nil,
		expLon: 1.0,
	}, {
		value:  "a",
		expErr: fmt.Errorf("invalid longitude value: %s", "a"),
		expLon: 0.0,
	}, {
		value:  "",
		expErr: fmt.Errorf("invalid longitude value: %s", ""),
		expLon: 0.0,
	}, {
		value:  " ",
		expErr: fmt.Errorf("invalid longitude value: %s", " "),
		expLon: 0.0,
	}, {
		value:  "180",
		expErr: nil,
		expLon: 180.0,
	}, {
		value:  "-180",
		expErr: nil,
		expLon: -180.0,
	}, {
		value:  "180.1",
		expErr: fmt.Errorf("longitude out of bounds: %f", 180.1),
		expLon: 0.0,
	}, {
		value:  "-180.1",
		expErr: fmt.Errorf("longitude out of bounds: %f", -180.1),
		expLon: 0.0,
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			c := Coords{}
			err := ParseLon(tc.value, &c)

			assert.Equal(tc.expErr, err)
			assert.Equal(tc.expLon, c.Lon)
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

func TestAppendParsedString(t *testing.T) {
	t.Parallel()

	tt := []struct {
		stringSlice []string
		value       string
		expStr      []string
	}{{
		stringSlice: []string{},
		value:       "string",
		expStr:      []string{"string"},
	}, {
		stringSlice: []string{"string"},
		value:       "string",
		expStr:      []string{"string", "string"},
	}, {
		stringSlice: []string{"string"},
		value:       " string",
		expStr:      []string{"string", "string"},
	}, {
		stringSlice: []string{"string"},
		value:       "string ",
		expStr:      []string{"string", "string"},
	}, {
		stringSlice: []string{"string"},
		value:       "",
		expStr:      []string{"string", ""},
	}, {
		stringSlice: []string{"string"},
		value:       " ",
		expStr:      []string{"string", ""},
	}}

	for _, tc := range tt {
		tc := tc

		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			str := tc.stringSlice
			appendParsedString(tc.value, &str)

			assert.Equal(tc.expStr, str)
		})
	}
}
