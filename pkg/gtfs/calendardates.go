package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

var (
	ErrEmptyCalendarDatesFile      = fmt.Errorf("empty calendar dates file")
	ErrInvalidCalendarDatesHeaders = fmt.Errorf("invalid calendar dates headers")
	ErrNoCalendarDatesRecords      = fmt.Errorf("no calendar dates records")
)

type CalendarDate struct {
	ServiceID     string   `json:"serviceId" csv:"service_id"`
	Date          string   `json:"date" csv:"date"`
	ExceptionType string   `json:"exceptionType" csv:"exception_type"`
	Unused        []string `json:"-" csv:"-"`
}

func parseCalendarDates(file *zip.File) ([]CalendarDate, error) {
	rc, err := file.Open()
	if err != nil {
		return []CalendarDate{}, err
	}
	defer rc.Close()

	lines, err := csv.NewReader(rc).ReadAll()
	if len(lines) == 0 {
		return []CalendarDate{}, ErrEmptyCalendarDatesFile
	}

	headers := lines[0]
	if err := validateCalendarDatesHeader(headers); err != nil {
		return []CalendarDate{}, err
	}

	records := lines[1:]
	if len(records) == 0 {
		return []CalendarDate{}, ErrNoCalendarDatesRecords
	}

	calendarDates := make([]CalendarDate, len(records))
	for i, field := range headers {
		for j, record := range records {
			switch field {
			case "service_id":
				calendarDates[j].ServiceID = record[i]
			case "date":
				calendarDates[j].Date = record[i]
			case "exception_type":
				calendarDates[j].ExceptionType = record[i]
			default:
				calendarDates[j].Unused = append(calendarDates[j].Unused, record[i])
			}
		}
	}

	return calendarDates, nil
}

func validateCalendarDatesHeader(headers []string) error {
	requiredFields := []struct {
		name  string
		found bool
	}{{
		name:  "service_id",
		found: false,
	}, {
		name:  "date",
		found: false,
	}, {
		name:  "exception_type",
		found: false,
	}}

	for _, h := range headers {
		for i, f := range requiredFields {
			if h == f.name {
				requiredFields[i].found = true
			}
		}
	}

	for _, f := range requiredFields {
		if !f.found {
			return ErrInvalidCalendarDatesHeaders
		}
	}

	return nil
}
