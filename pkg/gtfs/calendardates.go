package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
)

var (
	ErrEmptyCalendarDatesFile      = fmt.Errorf("empty calendar dates file")
	ErrInvalidCalendarDatesHeaders = fmt.Errorf("invalid calendar dates headers")
	ErrNoCalendarDatesRecords      = fmt.Errorf("no calendar dates records")
)

type CalendarDate struct {
	ServiceID     string `json:"serviceId"`
	Date          Time   `json:"date"`
	ExceptionType Enum   `json:"exceptionType"`

	unused []string
}

func (s *GTFSSchedule) parseCalendarDates(file *zip.File) error {
	s.CalendarDates = map[string]CalendarDate{}

	rc, err := file.Open()
	if err != nil {
		s.errors = append(s.errors, err)
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors = append(s.errors, ErrEmptyCalendarDatesFile)
		return ErrEmptyCalendarDatesFile
	}
	if err != nil {
		s.errors = append(s.errors, err)
		return err
	}

	for i := 0; ; i++ {
		record, err := r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			s.errors = append(s.errors, fmt.Errorf("empty record at line %d", i))
			return ErrNoCalendarDatesRecords
		}

		var cd CalendarDate
		for j, v := range record {
			switch headers[j] {
			case "service_id":
				cd.ServiceID = v
			case "date":
				if err := cd.Date.parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid date at line %d: %w", i, err))
				}
			case "exception_type":
				if err := cd.ExceptionType.Parse(v, Accessibility); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid exception type at line %d: %w", i, err))
				}
			default:
				cd.unused = append(cd.unused, v)
			}
		}
		s.CalendarDates[cd.ServiceID] = cd
	}

	return nil
}
