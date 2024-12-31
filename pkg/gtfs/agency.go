package gtfs

import (
	"fmt"
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
}

func (a Agency) key() string {
	return a.ID
}

func (a Agency) validate() errorList {
	var errs errorList

	if a.Name == "" {
		errs.add(fmt.Errorf("agency name is required"))
	}
	if a.URL == "" {
		errs.add(fmt.Errorf("agency URL is required"))
	}
	if a.Timezone == "" {
		errs.add(fmt.Errorf("agency timezone is required"))
	}

	return errs
}
