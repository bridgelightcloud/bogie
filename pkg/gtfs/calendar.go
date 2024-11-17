package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

var (
	ErrEmptyCalendarFile      = fmt.Errorf("empty calendar file")
	ErrInvalidCalendarHeaders = fmt.Errorf("invalid calendar headers")
	ErrNoCalendarRecords      = fmt.Errorf("no calendar records")
)

type Calendar struct {
	ServiceID string `json:"serviceId"`
	Monday    Enum   `json:"monday"`
	Tuesday   Enum   `json:"tuesday"`
	Wednesday Enum   `json:"wednesday"`
	Thursday  Enum   `json:"thursday"`
	Friday    Enum   `json:"friday"`
	Saturday  Enum   `json:"saturday"`
	Sunday    Enum   `json:"sunday"`
	StartDate Time   `json:"startDate"`
	EndDate   Time   `json:"endDate"`

	unused []string
}

func (s *GTFSSchedule) parseCalendar(file *zip.File) error {
	s.Calendar = map[string]Calendar{}

	rc, err := file.Open()
	if err != nil {
		s.errors = append(s.errors, err)
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors = append(s.errors, ErrEmptyCalendarFile)
		return ErrEmptyCalendarFile
	}
	if err != nil {
		s.errors = append(s.errors, err)
		return err
	}

	var record []string
	for {
		record, err = r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			s.errors = append(s.errors, fmt.Errorf("empty calendar record"))
			continue
		}

		if len(record) > len(headers) {
			s.errors = append(s.errors, fmt.Errorf("invalid calendar record: %v", record))
			continue
		}

		var c Calendar
		for j, value := range record {
			value = strings.TrimSpace(value)
			switch headers[j] {
			case "service_id":
				c.ServiceID = value
			case "monday":
				if err := c.Monday.Parse(value, Availability); err != nil {
					s.errors = append(s.errors, err)
				}
			case "tuesday":
				if err := c.Tuesday.Parse(value, Availability); err != nil {
					s.errors = append(s.errors, err)
				}
			case "wednesday":
				if err := c.Wednesday.Parse(value, Availability); err != nil {
					s.errors = append(s.errors, err)
				}
			case "thursday":
				if err := c.Thursday.Parse(value, Availability); err != nil {
					s.errors = append(s.errors, err)
				}
			case "friday":
				if err := c.Friday.Parse(value, Availability); err != nil {
					s.errors = append(s.errors, err)
				}
			case "saturday":
				if err := c.Saturday.Parse(value, Availability); err != nil {
					s.errors = append(s.errors, err)
				}
			case "sunday":
				if err := c.Sunday.Parse(value, Availability); err != nil {
					s.errors = append(s.errors, err)
				}
			case "start_date":
				if err := c.StartDate.parse(value); err != nil {
					s.errors = append(s.errors, err)
				}
			case "end_date":
				if err := c.EndDate.parse(value); err != nil {
					s.errors = append(s.errors, err)
				}
			default:
				c.unused = append(c.unused, value)
			}
		}
		if _, ok := s.Calendar[c.ServiceID]; ok {
			s.errors = append(s.errors, fmt.Errorf("duplicate calendar record: %s", c.ServiceID))
			continue
		} else {
			s.Calendar[c.ServiceID] = c
		}
	}

	if err != io.EOF {
		s.errors = append(s.errors, err)
		return err
	}

	if len(s.Calendar) == 0 {
		s.errors = append(s.errors, ErrNoCalendarRecords)
	}

	return nil
}

// func validateCalendarHeader(headers []string) error {
// 	requiredFields := []struct {
// 		name  string
// 		found bool
// 	}{{
// 		name:  "service_id",
// 		found: false,
// 	}, {
// 		name:  "monday",
// 		found: false,
// 	}, {
// 		name:  "tuesday",
// 		found: false,
// 	}, {
// 		name:  "wednesday",
// 		found: false,
// 	}, {
// 		name:  "thursday",
// 		found: false,
// 	}, {
// 		name:  "friday",
// 		found: false,
// 	}, {
// 		name:  "saturday",
// 		found: false,
// 	}, {
// 		name:  "sunday",
// 		found: false,
// 	}, {
// 		name:  "start_date",
// 		found: false,
// 	}, {
// 		name:  "end_date",
// 		found: false,
// 	}}

// 	for _, field := range headers {
// 		for i, req := range requiredFields {
// 			if field == req.name {
// 				requiredFields[i].found = true
// 			}
// 		}
// 	}

// 	for _, req := range requiredFields {
// 		if !req.found {
// 			return ErrInvalidCalendarHeaders
// 		}
// 	}

// 	return nil
// }
