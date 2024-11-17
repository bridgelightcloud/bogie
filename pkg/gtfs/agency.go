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
	ID          String `json:"agencyId,omitempty"`
	Name        String `json:"agencyName"`
	URL         String `json:"agencyUrl"`
	Timezone    String `json:"agencyTimezone"`
	Lang        String `json:"agencyLang,omitempty"`
	Phone       String `json:"agencyPhone,omitempty"`
	FareURL     String `json:"agencyFareUrl,omitempty"`
	AgencyEmail String `json:"agencyEmail,omitempty"`
	unused      []String

	route []String
}

func (s *GTFSSchedule) parseAgencies(file *zip.File) error {
	s.Agencies = map[string]Agency{}

	rc, err := file.Open()
	if err != nil {
		s.errors = append(s.errors, fmt.Errorf("error opening agency file: %w", err))
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors = append(s.errors, ErrEmptyAgencyFile)
		return ErrEmptyAgencyFile
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
			s.errors = append(s.errors, fmt.Errorf("empty agency record"))
			continue
		}

		if len(record) > len(headers) {
			s.errors = append(s.errors, fmt.Errorf("record has too many columns"))
		}

		var a Agency
		for j, v := range record {
			v = strings.TrimSpace(v)
			switch headers[j] {
			case "agency_id":
				if err := a.ID.Parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid agency_id: %w", err))
				}
			case "agency_name":
				if err := a.Name.Parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid agency_name: %w", err))
				}
			case "agency_url":
				if err := a.URL.Parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid agency_url: %w", err))
				}
			case "agency_timezone":
				if err := a.Timezone.Parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid agency_timezone: %w", err))
				}
			case "agency_lang":
				if err := a.Lang.Parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid agency_lang: %w", err))
				}
			case "agency_phone":
				if err := a.Phone.Parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid agency_phone: %w", err))
				}
			case "agency_fare_url":
				if err := a.FareURL.Parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid agency_fare_url: %w", err))
				}
			case "agency_email":
				if err := a.AgencyEmail.Parse(v); err != nil {
					s.errors = append(s.errors, fmt.Errorf("invalid agency_email: %w", err))
				}
			default:
				a.unused = append(a.unused, String(strings.TrimSpace(v)))
			}
		}
		s.Agencies[string(a.ID)] = a
	}

	if err != io.EOF {
		s.errors = append(s.errors, err)
		return err
	}

	if len(s.Agencies) == 0 {
		s.errors = append(s.errors, ErrNoAgencyRecords)
		return ErrNoAgencyRecords
	}

	return nil
}
