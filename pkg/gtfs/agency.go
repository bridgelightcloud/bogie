package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
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

	errors   errorList
	warnings errorList
}

func (a Agency) IsValid() bool {
	return len(a.errors) == 0
}

func (s *GTFSSchedule) parseAgencies(file *zip.File) {
	s.Agencies = map[string]Agency{}

	rc, err := file.Open()
	if err != nil {
		s.errors.add(fmt.Errorf("error opening agency file: %w", err))
		return
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors.add(fmt.Errorf("empty agency file"))
		return
	}
	if err != nil {
		s.errors.add(err)
		return
	}

	var record []string
	for {
		record, err = r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			s.errors.add(fmt.Errorf("empty agency record"))
			continue
		}

		if len(record) > len(headers) {
			s.errors.add(fmt.Errorf("record has too many columns"))
		}

		var a Agency
		for j, v := range record {
			v = strings.TrimSpace(v)
			switch headers[j] {
			case "agency_id":
				ParseString(v, &a.ID)
			case "agency_name":
				ParseString(v, &a.Name)
			case "agency_url":
				ParseString(v, &a.URL)
			case "agency_timezone":
				ParseString(v, &a.Timezone)
			case "agency_lang":
				ParseString(v, &a.Lang)
			case "agency_phone":
				ParseString(v, &a.Phone)
			case "agency_fare_url":
				ParseString(v, &a.FareURL)
			case "agency_email":
				ParseString(v, &a.AgencyEmail)
			default:
				a.unused = append(a.unused, strings.TrimSpace(v))
			}
		}
		validateAgency(&a)
		s.Agencies[a.ID] = a
	}

	if err != io.EOF {
		s.errors.add(err)
		return
	}

	if len(s.Agencies) == 0 {
		s.errors.add(fmt.Errorf("no agency records"))
		return
	}
}

func validateAgency(a *Agency) {
	if a.Name == "" {
		a.errors.add(fmt.Errorf("agency name is required"))
	}

	if a.URL == "" {
		a.errors.add(fmt.Errorf("agency URL is required"))
	}

	if a.Timezone == "" {
		a.errors.add(fmt.Errorf("agency timezone is required"))
	} else {

	}

	if a.Lang != "" {
		// validate language code
	}

	if a.Phone != "" {
		// validate phone number
	}

	if a.FareURL != "" {
		// validate URL
	}

	if a.AgencyEmail != "" {
		// validate email
	}
}
