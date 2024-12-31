package gtfs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var validColor = regexp.MustCompile(`(?i)^[a-f\d]{6}$`)

func ParseColor(v string, c *string) error {
	f := strings.TrimSpace(v)
	if !validColor.MatchString(f) {
		return fmt.Errorf("invalid color: %s", v)
	}

	*c = strings.ToUpper(f)
	return nil
}

var validCurrencyCodes = map[string]int{
	"AED": 2,
	"AFN": 2,
	"ALL": 2,
	"AMD": 2,
	"ANG": 2,
	"AOA": 2,
	"ARS": 2,
	"AUD": 2,
	"AWG": 2,
	"AZN": 2,
	"BAM": 2,
	"BBD": 2,
	"BDT": 2,
	"BGN": 2,
	"BHD": 3,
	"BIF": 0,
	"BMD": 2,
	"BND": 2,
	"BOB": 2,
	"BOV": 2,
	"BRL": 2,
	"BSD": 2,
	"BTN": 2,
	"BWP": 2,
	"BYN": 2,
	"BZD": 2,
	"CAD": 2,
	"CDF": 2,
	"CHE": 2,
	"CHF": 2,
	"CHW": 2,
	"CLF": 4,
	"CLP": 0,
	"CNY": 2,
	"COP": 2,
	"COU": 2,
	"CRC": 2,
	"CUP": 2,
	"CVE": 2,
	"CZK": 2,
	"DJF": 0,
	"DKK": 2,
	"DOP": 2,
	"DZD": 2,
	"EGP": 2,
	"ERN": 2,
	"ETB": 2,
	"EUR": 2,
	"FJD": 2,
	"FKP": 2,
	"GBP": 2,
	"GEL": 2,
	"GHS": 2,
	"GIP": 2,
	"GMD": 2,
	"GNF": 0,
	"GTQ": 2,
	"GYD": 2,
	"HKD": 2,
	"HNL": 2,
	"HTG": 2,
	"HUF": 2,
	"IDR": 2,
	"ILS": 2,
	"INR": 2,
	"IQD": 3,
	"IRR": 2,
	"ISK": 0,
	"JMD": 2,
	"JOD": 3,
	"JPY": 0,
	"KES": 2,
	"KGS": 2,
	"KHR": 2,
	"KMF": 0,
	"KPW": 2,
	"KRW": 0,
	"KWD": 3,
	"KYD": 2,
	"KZT": 2,
	"LAK": 2,
	"LBP": 2,
	"LKR": 2,
	"LRD": 2,
	"LSL": 2,
	"LYD": 3,
	"MAD": 2,
	"MDL": 2,
	"MGA": 2,
	"MKD": 2,
	"MMK": 2,
	"MNT": 2,
	"MOP": 2,
	"MRU": 2,
	"MUR": 2,
	"MVR": 2,
	"MWK": 2,
	"MXN": 2,
	"MXV": 2,
	"MYR": 2,
	"MZN": 2,
	"NAD": 2,
	"NGN": 2,
	"NIO": 2,
	"NOK": 2,
	"NPR": 2,
	"NZD": 2,
	"OMR": 3,
	"PAB": 2,
	"PEN": 2,
	"PGK": 2,
	"PHP": 2,
	"PKR": 2,
	"PLN": 2,
	"PYG": 0,
	"QAR": 2,
	"RON": 2,
	"RSD": 2,
	"RUB": 2,
	"RWF": 0,
	"SAR": 2,
	"SBD": 2,
	"SCR": 2,
	"SDG": 2,
	"SEK": 2,
	"SGD": 2,
	"SHP": 2,
	"SLE": 2,
	"SOS": 2,
	"SRD": 2,
	"SSP": 2,
	"STN": 2,
	"SVC": 2,
	"SYP": 2,
	"SZL": 2,
	"THB": 2,
	"TJS": 2,
	"TMT": 2,
	"TND": 3,
	"TOP": 2,
	"TRY": 2,
	"TTD": 2,
	"TWD": 2,
	"TZS": 2,
	"UAH": 2,
	"UGX": 0,
	"USD": 2,
	"USN": 2,
	"UYI": 0,
	"UYU": 2,
	"UYW": 4,
	"UZS": 2,
	"VED": 2,
	"VES": 2,
	"VND": 0,
	"VUV": 0,
	"WST": 2,
	"YER": 2,
	"ZAR": 2,
	"ZMW": 2,
	"ZWG": 2,
}

func ParseCurrencyCode(v string, c *string) error {
	f := strings.TrimSpace(v)
	f = strings.ToUpper(f)
	if _, ok := validCurrencyCodes[f]; !ok {
		return fmt.Errorf("invalid currency code: %s", v)
	}

	*c = f
	return nil
}

type Date struct {
	time.Time
}

var dateFormat = "20060102"

func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.Format(dateFormat)), nil
}

func (d *Date) UnmarshalText(text []byte) error {
	p, err := time.Parse(dateFormat, string(text))
	if err != nil {
		return fmt.Errorf("invalid date value: %s", text)
	}
	d.Time = p
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", d.Unix())), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid date value: %s", string(data))
	}
	*d = Date{time.Unix(i, 0)}
	return nil
}

type Time struct {
	time.Time
}

var timeFormat = "15:04:05"

func (t Time) MarshalText() ([]byte, error) {
	timeStr := t.Format(timeFormat)

	if d := t.Time.Day(); d > 1 {
		hrs := strconv.Itoa(t.Hour() + 24)
		return []byte(hrs + timeStr[2:]), nil
	}

	return []byte(timeStr), nil
}

func (t *Time) UnmarshalText(text []byte) error {
	timeStr := string(text)

	p, err := time.Parse(timeFormat, timeStr)

	if err != nil {
		if len(timeStr) < 8 {
			return fmt.Errorf("invalid time value: %s", text)
		}
		hrs := timeStr[:2]
		h, err := strconv.Atoi(hrs)
		if err != nil || h < 24 {
			return fmt.Errorf("invalid time value: %s", text)
		}

		timeStr = strconv.Itoa(h-24) + timeStr[2:]

		p, err = time.Parse(timeFormat, timeStr)

		if err != nil {
			return fmt.Errorf("invalid time value: %s", text)
		}

		t.Time = p.AddDate(0, 0, 1)
	} else {
		t.Time = p
	}

	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if str := string(data); str == "null" {
		t.Time = time.Time{}
	} else if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		*t = Time{Time: time.Unix(i, 0)}
	} else {
		return fmt.Errorf("invalid time value: %s", str)
	}

	return nil
}

type enumBounds struct {
	L int
	U int
}

var (
	Availability enumBounds = enumBounds{0, 1}
	Available    int        = 0
	Unavailable  int        = 1

	BikesAllowed                 enumBounds = enumBounds{0, 2}
	NoInfo                       int        = 0
	AtLeastOneBicycleAccomodated int        = 1
	NoBicyclesAllowed            int        = 2

	ContinuousPickup   enumBounds = enumBounds{0, 3}
	ContinuousDropOff  enumBounds = enumBounds{0, 3}
	DropOffType        enumBounds = enumBounds{0, 3}
	PickupType         enumBounds = enumBounds{0, 3}
	RegularlyScheduled int        = 0
	NoneAvailable      int        = 1
	MustPhoneAgency    int        = 2
	MustCoordinate     int        = 3

	DirectionID       enumBounds = enumBounds{0, 1}
	OneDirection      int        = 0
	OppositeDirection int        = 1

	ExceptionType enumBounds = enumBounds{1, 2}
	Added         int        = 1
	Removed       int        = 2

	Timepoint       enumBounds = enumBounds{0, 1}
	ApproximateTime int        = 0
	ExactTime       int        = 1

	LocationType enumBounds = enumBounds{0, 4}
	StopPlatform int        = 0
	Station      int        = 1
	EntranceExit int        = 2
	GenericNode  int        = 3
	BoardingArea int        = 4

	WheelchairAccessible            enumBounds = enumBounds{0, 2}
	UnknownAccessibility            int        = 0
	AtLeastOneWheelchairAccomodated int        = 1
	NoWheelchairsAccomodated        int        = 2
)

type errorList []error

func (e *errorList) add(err error) error {
	if err == nil {
		return err
	}
	*e = append(*e, err)
	return err
}
