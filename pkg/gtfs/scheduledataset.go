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

type gtfsScheduleZip struct {
	// Required files
	Agencies      *zip.File
	Routes        *zip.File
	Stops         *zip.File
	Trips         *zip.File
	StopTimes     *zip.File
	Calendar      *zip.File
	CalendarDates *zip.File

	// Optional files
	FareAttributes     *zip.File
	FareRules          *zip.File
	Timeframes         *zip.File
	FareMedia          *zip.File
	FareProducts       *zip.File
	FareLegRules       *zip.File
	FareTransferRules  *zip.File
	Areas              *zip.File
	StopAreas          *zip.File
	Networks           *zip.File
	RouteNetworks      *zip.File
	Shapes             *zip.File
	Frequencies        *zip.File
	Transfers          *zip.File
	Pathways           *zip.File
	Levels             *zip.File
	LocationGroups     *zip.File
	LocationGroupStops *zip.File
	LocationsGeojson   *zip.File
	BookingRules       *zip.File
	Translations       *zip.File
	FeedInfo           *zip.File
	Attributions       *zip.File

	// Additional files
	AdditionalFiles []*zip.File
}

type GTFSScheduleDataset struct {
	// Required files
	Agencies      []Agency
	Routes        []Route
	Stops         []Stop
	Trips         []Trip
	StopTimes     []StopTime
	Calendar      []Calendar
	CalendarDates []CalendarDate
}

type GTFSSchedule struct {
}

func OpenScheduleFromFile(fn string) (GTFSSchedule, error) {
	r, err := zip.OpenReader(fn)
	if err != nil {
		return GTFSSchedule{}, err
	}
	defer r.Close()

	gz, err := unzip(r)
	if err != nil {
		return GTFSSchedule{}, err
	}

	_, err = ParseSchedule(gz)
	if err != nil {
		return GTFSSchedule{}, err
	}

	return GTFSSchedule{}, nil
}

func unzip(r *zip.ReadCloser) (gtfsScheduleZip, error) {
	sz := gtfsScheduleZip{}

	for _, f := range r.File {
		switch f.Name {
		case "agency.txt":
			sz.Agencies = f
		case "routes.txt":
			sz.Routes = f
		case "stops.txt":
			sz.Stops = f
		case "trips.txt":
			sz.Trips = f
		case "stop_times.txt":
			sz.StopTimes = f
		case "calendar.txt":
			sz.Calendar = f
		case "calendar_dates.txt":
			sz.CalendarDates = f
		case "fare_attributes.txt":
			sz.FareAttributes = f
		case "fare_rules.txt":
			sz.FareRules = f
		case "timeframes.txt":
			sz.Timeframes = f
		case "fare_media.txt":
			sz.FareMedia = f
		case "fare_products.txt":
			sz.FareProducts = f
		case "fare_leg_rules.txt":
			sz.FareLegRules = f
		case "fare_transfer_rules.txt":
			sz.FareTransferRules = f
		case "areas.txt":
			sz.Areas = f
		case "stop_areas.txt":
			sz.StopAreas = f
		case "networks.txt":
			sz.Networks = f
		case "route_networks.txt":
			sz.RouteNetworks = f
		case "shapes.txt":
			sz.Shapes = f
		case "frequencies.txt":
			sz.Frequencies = f
		case "transfers.txt":
			sz.Transfers = f
		case "pathways.txt":
			sz.Pathways = f
		case "levels.txt":
			sz.Levels = f
		case "location_groups.txt":
			sz.LocationGroups = f
		case "location_group_stops.txt":
			sz.LocationGroupStops = f
		case "locations.geojson":
			sz.LocationsGeojson = f
		case "booking_rules.txt":
			sz.BookingRules = f
		case "translations.txt":
			sz.Translations = f
		case "feed_info.txt":
			sz.FeedInfo = f
		case "attributions.txt":
			sz.Attributions = f
		default:
			if sz.AdditionalFiles == nil {
				sz.AdditionalFiles = []*zip.File{f}
			} else {
				sz.AdditionalFiles = append(sz.AdditionalFiles, f)
			}
		}
	}

	// check that all required files are present
	if sz.Routes == nil {
		return sz, ErrMissingRoutes
	}
	if sz.Trips == nil {
		return sz, ErrMissingTrips
	}
	if sz.Stops == nil {
		return sz, ErrMissingStops
	}
	if sz.StopTimes == nil {
		return sz, ErrMissingStopTimes
	}
	if sz.Calendar == nil {
		return sz, ErrMissingCalendar
	}
	if sz.CalendarDates == nil {
		return sz, ErrMissingCalendarDates
	}

	return sz, nil
}

func ParseSchedule(sf gtfsScheduleZip) (GTFSScheduleDataset, error) {
	sd := GTFSScheduleDataset{}

	a, err := parseAgencies(sf.Agencies)
	if err != nil {
		return sd, err
	}
	sd.Agencies = a

	r, err := parseRoutes(sf.Routes)
	if err != nil {
		return sd, err
	}
	sd.Routes = r

	s, err := parseStops(sf.Stops)
	if err != nil {
		return sd, err
	}
	sd.Stops = s

	t, err := parseTrips(sf.Trips)
	if err != nil {
		return sd, err
	}
	sd.Trips = t

	st, err := parseStopTimes(sf.StopTimes)
	if err != nil {
		return sd, err
	}
	sd.StopTimes = st

	c, err := parseCalendar(sf.Calendar)
	if err != nil {
		return sd, err
	}
	sd.Calendar = c

	cd, err := parseCalendarDates(sf.CalendarDates)
	if err != nil {
		return sd, err
	}
	sd.CalendarDates = cd

	return sd, nil
}
