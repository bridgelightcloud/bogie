package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

var (
	ErrEmptyAgencyFile      = fmt.Errorf("empty agency file")
	ErrInvalidAgencyHeaders = fmt.Errorf("invalid agency headers")
	ErrNoAgencyRecords      = fmt.Errorf("no agency records")
)

type Agency struct {
	ID          string   `json:"agencyId,omitempty" csv:"agency_id,omitempty"`
	Name        string   `json:"agencyName" csv:"agency_name"`
	URL         string   `json:"agencyUrl" csv:"agency_url"`
	Timezone    string   `json:"agencyTimezone" csv:"agency_timezone"`
	Lang        string   `json:"agencyLang,omitempty" csv:"agency_lang,omitempty"`
	Phone       string   `json:"agencyPhone,omitempty" csv:"agency_phone,omitempty"`
	FareURL     string   `json:"agencyFareUrl,omitempty" csv:"agency_fare_url,omitempty"`
	AgencyEmail string   `json:"agencyEmail,omitempty" csv:"agency_email,omitempty"`
	Unused      []string `json:"-" csv:"-"`
}

func parseAgencies(file *zip.File) ([]Agency, error) {
	rc, err := file.Open()
	if err != nil {
		return []Agency{}, err
	}
	defer rc.Close()

	lines, err := csv.NewReader(rc).ReadAll()
	if len(lines) == 0 {
		return []Agency{}, ErrEmptyAgencyFile
	}

	headers := lines[0]
	if err := validateAgenciesHeader(headers); err != nil {
		return []Agency{}, err
	}

	records := lines[1:]
	if len(records) == 0 {
		return []Agency{}, ErrNoAgencyRecords
	}

	agencies := make([]Agency, len(records))
	for i, field := range headers {
		for j, record := range records {
			switch field {
			case "agency_id":
				agencies[j].ID = record[i]
			case "agency_name":
				agencies[j].Name = record[i]
			case "agency_url":
				agencies[j].URL = record[i]
			case "agency_timezone":
				agencies[j].Timezone = record[i]
			case "agency_lang":
				agencies[j].Lang = record[i]
			case "agency_phone":
				agencies[j].Phone = record[i]
			case "agency_fare_url":
				agencies[j].FareURL = record[i]
			case "agency_email":
				agencies[j].AgencyEmail = record[i]
			default:
				if agencies[j].Unused == nil {
					agencies[j].Unused = []string{record[i]}
				} else {
					agencies[j].Unused = append(agencies[j].Unused, record[i])
				}
			}
		}
	}

	return agencies, nil
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
