package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

var (
	ErrEmptyCalendarFile      = fmt.Errorf("empty calendar file")
	ErrInvalidCalendarHeaders = fmt.Errorf("invalid calendar headers")
	ErrNoCalendarRecords      = fmt.Errorf("no calendar records")
)

type Calendar struct {
	ServiceID string    `json:"serviceId"`
	Monday    int       `json:"monday"`
	Tuesday   int       `json:"tuesday"`
	Wednesday int       `json:"wednesday"`
	Thursday  int       `json:"thursday"`
	Friday    int       `json:"friday"`
	Saturday  int       `json:"saturday"`
	Sunday    int       `json:"sunday"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	unused []string
}

func (s *GTFSSchedule) parseCalendar(file *zip.File) error {
	s.Calendar = map[string]Calendar{}

	rc, err := file.Open()
	if err != nil {
		s.errors.add(err)
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors.add(ErrEmptyCalendarFile)
		return ErrEmptyCalendarFile
	}
	if err != nil {
		s.errors.add(err)
		return err
	}

	var record []string
	for {
		record, err = r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			s.errors.add(fmt.Errorf("empty calendar record"))
			continue
		}

		if len(record) > len(headers) {
			s.errors.add(fmt.Errorf("invalid calendar record: %v", record))
			continue
		}

		var c Calendar
		for j, value := range record {
			switch headers[j] {
			case "service_id":
				c.ServiceID = value
			case "monday":
				if err := ParseEnum(value, Availability, &c.Monday); err != nil {
					s.errors.add(err)
				}
			case "tuesday":
				if err := ParseEnum(value, Availability, &c.Tuesday); err != nil {
					s.errors.add(err)
				}
			case "wednesday":
				if err := ParseEnum(value, Availability, &c.Wednesday); err != nil {
					s.errors.add(err)
				}
			case "thursday":
				if err := ParseEnum(value, Availability, &c.Thursday); err != nil {
					s.errors.add(err)
				}
			case "friday":
				if err := ParseEnum(value, Availability, &c.Friday); err != nil {
					s.errors.add(err)
				}
			case "saturday":
				if err := ParseEnum(value, Availability, &c.Saturday); err != nil {
					s.errors.add(err)
				}
			case "sunday":
				if err := ParseEnum(value, Availability, &c.Sunday); err != nil {
					s.errors.add(err)
				}
			case "start_date":
				if err := ParseDate(value, &c.StartDate); err != nil {
					s.errors.add(err)
				}
			case "end_date":
				if err := ParseDate(value, &c.EndDate); err != nil {
					s.errors.add(err)
				}
			default:
				appendParsedString(value, &c.unused)
			}
		}
		if _, ok := s.Calendar[c.ServiceID]; ok {
			s.errors.add(fmt.Errorf("duplicate calendar record: %s", c.ServiceID))
			continue
		} else {
			s.Calendar[c.ServiceID] = c
		}
	}

	if err != io.EOF {
		s.errors.add(err)
		return err
	}

	if len(s.Calendar) == 0 {
		s.errors.add(ErrNoCalendarRecords)
	}

	return nil
}
