package gtfs

import (
	"archive/zip"
	"fmt"
)

// Errors
var (
	ErrBadScheduleFile      = fmt.Errorf("bad schedule file")
	ErrMissingAgency        = fmt.Errorf("missing agency file")
	ErrMissingRoutes        = fmt.Errorf("missing routes file")
	ErrMissingTrips         = fmt.Errorf("missing trips file")
	ErrMissingStops         = fmt.Errorf("missing stops file")
	ErrMissingStopTimes     = fmt.Errorf("missing stop times file")
	ErrMissingCalendar      = fmt.Errorf("missing calendar file")
	ErrMissingCalendarDates = fmt.Errorf("missing calendar dates file")
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

	unusedFiles []string
	errors      errorList
}

func OpenScheduleFromFile(fn string) (GTFSSchedule, error) {
	r, err := zip.OpenReader(fn)
	if err != nil {
		return GTFSSchedule{}, err
	}
	defer r.Close()

	sd := parseSchedule(r)

	return sd, nil
}

func parseSchedule(r *zip.ReadCloser) GTFSSchedule {
	s := GTFSSchedule{}

	files := make(map[string]*zip.File)
	for _, f := range r.File {
		files[f.Name] = f
	}

	if f, ok := files["agency.txt"]; !ok {
		s.errors.add(ErrMissingAgency)
	} else if err := s.parseAgencies(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["stops.txt"]; !ok {
		s.errors.add(ErrMissingStops)
	} else if err := s.parseStopsData(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["routes.txt"]; !ok {
		s.errors.add(ErrMissingRoutes)
	} else if err := s.parseRoutes(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["calendar.txt"]; !ok {
		s.errors.add(ErrMissingCalendar)
	} else if err := s.parseCalendar(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["calendar_dates.txt"]; !ok {
		s.errors.add(ErrMissingCalendarDates)
	} else if err := s.parseCalendarDates(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["trips.txt"]; !ok {
		s.errors.add(ErrMissingTrips)
	} else if err := s.parseTrips(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["stop_times.txt"]; !ok {
		s.errors.add(ErrMissingStopTimes)
	} else if err := s.parseStopTimes(f); err != nil {
		s.errors.add(err)
	}

	// f, ok = files["trips.txt"]
	// f, ok = files["stop_times.txt"]
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
