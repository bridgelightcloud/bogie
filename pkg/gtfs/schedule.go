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

	sd, err := parseSchedule(r)
	if err != nil {
		return GTFSSchedule{}, err
	}

	return sd, nil
}

func parseSchedule(r *zip.ReadCloser) (GTFSSchedule, error) {
	s := GTFSSchedule{}

	files := make(map[string]*zip.File)
	for _, f := range r.File {
		files[f.Name] = f
	}

	if f, ok := files["agency.txt"]; !ok {
		return s, ErrMissingAgency
	} else if err := s.parseAgencies(f); err != nil {
		return s, err
	}

	if f, ok := files["stops.txt"]; !ok {
		return s, ErrMissingStops
	} else if err := s.parseStopsData(f); err != nil {
		return s, err
	}

	if f, ok := files["routes.txt"]; !ok {
		return s, ErrMissingRoutes
	} else if err := s.parseRoutes(f); err != nil {
		return s, err
	}

	if f, ok := files["calendar.txt"]; !ok {
		return s, ErrMissingCalendar
	} else if err := s.parseCalendar(f); err != nil {
		return s, err
	}

	if f, ok := files["calendar_dates.txt"]; !ok {
		return s, ErrMissingCalendarDates
	} else if err := s.parseCalendarDates(f); err != nil {
		return s, err
	}

	if f, ok := files["trips.txt"]; !ok {
		return s, ErrMissingTrips
	} else if err := s.parseTrips(f); err != nil {
		return s, err
	}

	if f, ok := files["stop_times.txt"]; !ok {
		return s, ErrMissingStopTimes
	} else if err := s.parseStopTimes(f); err != nil {
		return s, err
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

	return s, nil
}
