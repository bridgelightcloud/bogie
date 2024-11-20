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
		s.errors.add(fmt.Errorf("missing agency.txt"))
	} else {
		s.parseAgencies(f)
	}

	if f, ok := files["levels.txt"]; !ok {
		s.errors.add(fmt.Errorf("missing levels.txt"))
	} else {
		s.parseLevels(f)
	}

	if f, ok := files["stops.txt"]; !ok {
		s.errors.add(fmt.Errorf("missing stops.txt"))
	} else {
		s.parseStopsData(f)
	}

	if f, ok := files["routes.txt"]; !ok {
		s.errors.add(fmt.Errorf("missing routes.txt"))
	} else {
		s.parseRoutes(f)
	}

	if f, ok := files["calendar.txt"]; !ok {
		s.errors.add(fmt.Errorf("missing calendar.txt"))
	} else if err := s.parseCalendar(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["calendar_dates.txt"]; !ok {
		s.errors.add(fmt.Errorf("missing calendar_dates.txt"))
	} else if err := s.parseCalendarDates(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["trips.txt"]; !ok {
		s.errors.add(fmt.Errorf("missing trips.txt"))
	} else if err := s.parseTrips(f); err != nil {
		s.errors.add(err)
	}

	if f, ok := files["stop_times.txt"]; !ok {
		s.errors.add(fmt.Errorf("missing stop_times.txt"))
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
