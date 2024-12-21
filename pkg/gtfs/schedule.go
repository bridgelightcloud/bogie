package gtfs

import (
	"archive/zip"
	"fmt"
)

type GTFSSchedule struct {
	// Required files
	Agencies      map[string]Agency
	Stops         map[string]Stop
	Routes        map[string]Route
	Calendar      map[string]Calendar
	CalendarDates map[string]CalendarDate
	Trips         map[string]Trip
	StopTimes     map[string]StopTime
	Levels        map[string]Level

	unusedFiles []string
	errors      errorList
	warning     errorList
}

func (s GTFSSchedule) Errors() errorList {
	return s.errors
}

type gtfsSpec[R record] struct {
	setter func(*GTFSSchedule, map[string]R)
}

func (s gtfsSpec[R]) Parse(f *zip.File, schedule *GTFSSchedule, errors *errorList) {
	r, err := f.Open()
	if err != nil {
		errors.add(fmt.Errorf("error opening file: %w", err))
		return
	}
	defer r.Close()

	records := make(map[string]R)

	parse(r, records, errors)

	s.setter(schedule, records)
}

type parseableGtfs interface {
	Parse(*zip.File, *GTFSSchedule, *errorList)
}

var gtfsSpecs = map[string]parseableGtfs{
	"agency.txt":         gtfsSpec[Agency]{setter: func(s *GTFSSchedule, r map[string]Agency) { s.Agencies = r }},
	"stops.txt":          gtfsSpec[Stop]{setter: func(s *GTFSSchedule, r map[string]Stop) { s.Stops = r }},
	"routes.txt":         gtfsSpec[Route]{setter: func(s *GTFSSchedule, r map[string]Route) { s.Routes = r }},
	"calendar.txt":       gtfsSpec[Calendar]{setter: func(s *GTFSSchedule, r map[string]Calendar) { s.Calendar = r }},
	"calendar_dates.txt": gtfsSpec[CalendarDate]{setter: func(s *GTFSSchedule, r map[string]CalendarDate) { s.CalendarDates = r }},
	"trips.txt":          gtfsSpec[Trip]{setter: func(s *GTFSSchedule, r map[string]Trip) { s.Trips = r }},
	"stop_times.txt":     gtfsSpec[StopTime]{setter: func(s *GTFSSchedule, r map[string]StopTime) { s.StopTimes = r }},
	"levels.txt":         gtfsSpec[Level]{setter: func(s *GTFSSchedule, r map[string]Level) { s.Levels = r }},
}

func OpenScheduleFromZipFile(fn string) (GTFSSchedule, error) {
	r, err := zip.OpenReader(fn)
	if err != nil {
		return GTFSSchedule{}, err
	}
	defer r.Close()

	sd := parseSchedule(r)

	return sd, nil
}

func parseSchedule(r *zip.ReadCloser) GTFSSchedule {
	var s GTFSSchedule

	for _, f := range r.File {
		if spec, ok := gtfsSpecs[f.Name]; ok {
			spec.Parse(f, &s, &s.errors)
		}
	}

	// f, ok = files["fare_attributes.txt"]
	// f, ok = files["fare_rules.txt"]
	// f, ok = files["timeframes.txt"]
	// f, ok = files["fare_media.txt"]
	// f, ok = files["fare_products.txt"]
	// f, ok = files["fare_leg_rules.txt"]
	// f, ok = files["fare_transfer_rules.txt"]
	// f, ok = files["areas.txt"]
	// f, ok = files["stop_areas.txt"]
	// f, ok = files["networks.txt"]
	// f, ok = files["route_networks.txt"]
	// f, ok = files["shapes.txt"]
	// f, ok = files["frequencies.txt"]
	// f, ok = files["transfers.txt"]
	// f, ok = files["pathways.txt"]
	// f, ok = files["levels.txt"]
	// f, ok = files["location_groups.txt"]
	// f, ok = files["location_group_stops.txt"]
	// f, ok = files["locations.geojson"]
	// f, ok = files["booking_rules.txt"]
	// f, ok = files["translations.txt"]
	// f, ok = files["feed_info.txt"]
	// f, ok = files["attributions.txt"]

	return s
}
