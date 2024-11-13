package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

var (
	ErrEmptyCalendarFile      = fmt.Errorf("empty calendar file")
	ErrInvalidCalendarHeaders = fmt.Errorf("invalid calendar headers")
	ErrNoCalendarRecords      = fmt.Errorf("no calendar records")
)

type Calendar struct {
	ServiceID string   `json:"serviceId" csv:"service_id"`
	Monday    string   `json:"monday" csv:"monday"`
	Tuesday   string   `json:"tuesday" csv:"tuesday"`
	Wednesday string   `json:"wednesday" csv:"wednesday"`
	Thursday  string   `json:"thursday" csv:"thursday"`
	Friday    string   `json:"friday" csv:"friday"`
	Saturday  string   `json:"saturday" csv:"saturday"`
	Sunday    string   `json:"sunday" csv:"sunday"`
	StartDate string   `json:"startDate" csv:"start_date"`
	EndDate   string   `json:"endDate" csv:"end_date"`
	Unused    []string `json:"-" csv:"-"`
}

func parseCalendar(file *zip.File) ([]Calendar, error) {
	rc, err := file.Open()
	if err != nil {
		return []Calendar{}, err
	}
	defer rc.Close()

	lines, err := csv.NewReader(rc).ReadAll()
	if len(lines) == 0 {
		return []Calendar{}, ErrEmptyCalendarFile
	}

	headers := lines[0]
	if err := validateCalendarHeader(headers); err != nil {
		return []Calendar{}, err
	}

	records := lines[1:]
	if len(records) == 0 {
		return []Calendar{}, ErrNoCalendarRecords
	}

	calendar := make([]Calendar, len(records))
	for i, field := range headers {
		for j, record := range records {
			switch field {
			case "service_id":
				calendar[j].ServiceID = record[i]
			case "monday":
				calendar[j].Monday = record[i]
			case "tuesday":
				calendar[j].Tuesday = record[i]
			case "wednesday":
				calendar[j].Wednesday = record[i]
			case "thursday":
				calendar[j].Thursday = record[i]
			case "friday":
				calendar[j].Friday = record[i]
			case "saturday":
				calendar[j].Saturday = record[i]
			case "sunday":
				calendar[j].Sunday = record[i]
			case "start_date":
				calendar[j].StartDate = record[i]
			case "end_date":
				calendar[j].EndDate = record[i]
			default:
				calendar[j].Unused = append(calendar[j].Unused, record[i])
			}
		}
	}

	return calendar, nil
}

func validateCalendarHeader(headers []string) error {
	requiredFields := []struct {
		name  string
		found bool
	}{{
		name:  "service_id",
		found: false,
	}, {
		name:  "monday",
		found: false,
	}, {
		name:  "tuesday",
		found: false,
	}, {
		name:  "wednesday",
		found: false,
	}, {
		name:  "thursday",
		found: false,
	}, {
		name:  "friday",
		found: false,
	}, {
		name:  "saturday",
		found: false,
	}, {
		name:  "sunday",
		found: false,
	}, {
		name:  "start_date",
		found: false,
	}, {
		name:  "end_date",
		found: false,
	}}

	for _, field := range headers {
		for i, req := range requiredFields {
			if field == req.name {
				requiredFields[i].found = true
			}
		}
	}

	for _, req := range requiredFields {
		if !req.found {
			return ErrInvalidCalendarHeaders
		}
	}

	return nil
}
