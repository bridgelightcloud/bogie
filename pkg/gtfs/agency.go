package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

var (
	ErrEmptyAgencyFile      = fmt.Errorf("empty agency file")
	ErrInvalidAgencyHeaders = fmt.Errorf("invalid agency headers")
	ErrNoAgencyRecords      = fmt.Errorf("no agency records")
)

type Agency struct {
	ID          string `json:"agencyId,omitempty"`
	Name        string `json:"agencyName"`
	URL         string `json:"agencyUrl"`
	Timezone    string `json:"agencyTimezone"`
	Lang        string `json:"agencyLang,omitempty"`
	Phone       string `json:"agencyPhone,omitempty"`
	FareURL     string `json:"agencyFareUrl,omitempty"`
	AgencyEmail string `json:"agencyEmail,omitempty"`
	unused      []string

	route []string
}

func (s *GTFSSchedule) parseAgencies(file *zip.File) error {
	s.Agencies = map[string]Agency{}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		return ErrEmptyAgencyFile
	}
	if err != nil {
		return err
	}

	var record []string
	for {
		record, err = r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			continue
		}

		if len(record) > len(headers) {
			return fmt.Errorf("record has too many columns")
		}

		var agency Agency
		for j, value := range record {
			value = strings.TrimSpace(value)
			switch headers[j] {
			case "agency_id":
				agency.ID = value
			case "agency_name":
				agency.Name = value
			case "agency_url":
				agency.URL = value
			case "agency_timezone":
				agency.Timezone = value
			case "agency_lang":
				agency.Lang = value
			case "agency_phone":
				agency.Phone = value
			case "agency_fare_url":
				agency.FareURL = value
			case "agency_email":
				agency.AgencyEmail = value
			default:
				agency.unused = append(agency.unused, value)
			}
		}
		s.Agencies[agency.ID] = agency
	}

	if err != io.EOF {
		return err
	}

	if len(s.Agencies) == 0 {
		return ErrNoAgencyRecords
	}

	return nil
}

func validateAgenciesHeader(fields []string) error {
	requiredFields := []struct {
		name  string
		found bool
	}{{
		name:  "agency_name",
		found: false},
		{
			name:  "agency_url",
			found: false,
		},
		{
			name:  "agency_timezone",
			found: false,
		},
	}

	for _, field := range fields {
		for i, req := range requiredFields {
			if field == req.name {
				requiredFields[i].found = true
			}
		}
	}

	for _, req := range requiredFields {
		if !req.found {
			return ErrInvalidAgencyHeaders
		}
	}

	return nil
}
