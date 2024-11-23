package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/bridgelightcloud/bogie/pkg/csvmum"
)

type Calendar struct {
	ServiceID string    `json:"serviceId" csv:"service_id"`
	Monday    int       `json:"monday" csv:"monday"`
	Tuesday   int       `json:"tuesday" csv:"tuesday"`
	Wednesday int       `json:"wednesday" csv:"wednesday"`
	Thursday  int       `json:"thursday" csv:"thursday"`
	Friday    int       `json:"friday" csv:"friday"`
	Saturday  int       `json:"saturday" csv:"saturday"`
	Sunday    int       `json:"sunday" csv:"sunday"`
	StartDate time.Time `json:"startDate" csv:"start_date"`
	EndDate   time.Time `json:"endDate" csv:"end_date"`

	unused []string

	errors   errorList
	warnings errorList
}

func (s *GTFSSchedule) parseCalendar(file *zip.File) {
	rc, err := file.Open()
	if err != nil {
		s.errors.add(err)
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	data, err := r.ReadAll()
	if err != nil {
		s.errors.add(fmt.Errorf("error reading calendar file: %w", err))
		return
	}

	cs := []Calendar{}
	err = csvmum.Unmarshal(data, &cs)
	if err != nil {
		s.errors.add(fmt.Errorf("error unmarshalling calendar file: %w", err))
	}

	s.Calendar = make(map[string]Calendar, len(cs))
	for _, c := range cs {
		s.Calendar[c.ServiceID] = c
	}

	if len(s.Calendar) == 0 {
		s.errors.add(fmt.Errorf("no calendar records"))
	}
}
