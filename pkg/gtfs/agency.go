package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"

	"github.com/bridgelightcloud/bogie/pkg/csvmum"
)

type Agency struct {
	ID          string `json:"agencyId,omitempty" csv:"agency_id"`
	Name        string `json:"agencyName" csv:"agency_name"`
	URL         string `json:"agencyUrl" csv:"agency_url"`
	Timezone    string `json:"agencyTimezone" csv:"agency_timezone"`
	Lang        string `json:"agencyLang,omitempty" csv:"agency_lang"`
	Phone       string `json:"agencyPhone,omitempty" csv:"agency_phone"`
	FareURL     string `json:"agencyFareUrl,omitempty" csv:"agency_fare_url"`
	AgencyEmail string `json:"agencyEmail,omitempty" csv:"agency_email"`

	unused []string

	errors   errorList
	warnings errorList
}

func (a Agency) IsValid() bool {
	return len(a.errors) == 0
}

func (s *GTFSSchedule) parseAgencies(file *zip.File) {
	rc, err := file.Open()
	if err != nil {
		s.errors.add(fmt.Errorf("error opening agency file: %w", err))
		return
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	data, err := r.ReadAll()
	if err != nil {
		s.errors.add(fmt.Errorf("error reading agency file: %w", err))
		return
	}

	as := []Agency{}
	err = csvmum.Unmarshal(data, &as)
	if err != nil {
		s.errors.add(fmt.Errorf("error unmarshalling agency file: %w", err))
	}

	s.Agencies = make(map[string]Agency, len(as))
	for _, a := range as {
		validateAgency(&a)
		s.Agencies[a.ID] = a
	}

	if len(s.Agencies) == 0 {
		s.errors.add(fmt.Errorf("no agency records"))
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
