package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

var (
	ErrEmptyCalendarDatesFile      = fmt.Errorf("empty calendar dates file")
	ErrInvalidCalendarDatesHeaders = fmt.Errorf("invalid calendar dates headers")
	ErrNoCalendarDatesRecords      = fmt.Errorf("no calendar dates records")
)

type CalendarDate struct {
	ServiceID     string    `json:"serviceId"`
	Date          time.Time `json:"date"`
	ExceptionType int       `json:"exceptionType"`

	unused []string
}

func (s *GTFSSchedule) parseCalendarDates(file *zip.File) error {
	s.CalendarDates = map[string]CalendarDate{}

	rc, err := file.Open()
	if err != nil {
		s.errors.add( err)
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors.add(ErrEmptyCalendarDatesFile)
		return ErrEmptyCalendarDatesFile
	}
	if err != nil {
		s.errors.add( err)
		return err
	}

	for i := 0; ; i++ {
		record, err := r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			s.errors.add(fmt.Errorf("empty record at line %d", i))
			return ErrNoCalendarDatesRecords
		}

		var cd CalendarDate
		for j, v := range record {
			switch headers[j] {
			case "service_id":
				ParseString(v, &cd.ServiceID)
			case "date":
				if err := ParseDate(v, &cd.Date); err != nil {
					s.errors.add( fmt.Errorf("invalid date at line %d: %w", i, err))
				}
			case "exception_type":
				if err := ParseEnum(v, ExceptionType, &cd.ExceptionType); err != nil {
					s.errors.add(fmt.Errorf("invalid exception_type at line %d: %w", i, err))
				}
			default:
				cd.unused = append(cd.unused, v)
			}
		}
		s.CalendarDates[cd.ServiceID] = cd
	}

	return nil
}
