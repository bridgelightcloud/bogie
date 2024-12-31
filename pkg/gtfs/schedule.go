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
	warnings    errorList
}

func (s GTFSSchedule) Errors() errorList {
	return s.errors
}

type gtfsSpec[R record] struct {
	set func(*GTFSSchedule, map[string]R)
}

type fileParser interface {
	parseFile(*zip.File, *GTFSSchedule, *errorList)
}

func (spec gtfsSpec[R]) parseFile(f *zip.File, schedule *GTFSSchedule, errors *errorList) {
	r, err := f.Open()
	if err != nil {
		errors.add(fmt.Errorf("error opening file: %w", err))
		return
	}
	defer r.Close()

	records := make(map[string]R)

	parse(r, records, errors)

	spec.set(schedule, records)
}

var gtfsSpecs = map[string]fileParser{
	"agency.txt":         gtfsSpec[Agency]{set: func(s *GTFSSchedule, r map[string]Agency) { s.Agencies = r }},
	"stops.txt":          gtfsSpec[Stop]{set: func(s *GTFSSchedule, r map[string]Stop) { s.Stops = r }},
	"routes.txt":         gtfsSpec[Route]{set: func(s *GTFSSchedule, r map[string]Route) { s.Routes = r }},
	"calendar.txt":       gtfsSpec[Calendar]{set: func(s *GTFSSchedule, r map[string]Calendar) { s.Calendar = r }},
	"calendar_dates.txt": gtfsSpec[CalendarDate]{set: func(s *GTFSSchedule, r map[string]CalendarDate) { s.CalendarDates = r }},
	"trips.txt":          gtfsSpec[Trip]{set: func(s *GTFSSchedule, r map[string]Trip) { s.Trips = r }},
	"stop_times.txt":     gtfsSpec[StopTime]{set: func(s *GTFSSchedule, r map[string]StopTime) { s.StopTimes = r }},
	"levels.txt":         gtfsSpec[Level]{set: func(s *GTFSSchedule, r map[string]Level) { s.Levels = r }},
}

func OpenScheduleFromZipFile(fn string) (GTFSSchedule, error) {
	r, err := zip.OpenReader(fn)
	if err != nil {
		return GTFSSchedule{}, err
	}
	defer r.Close()

	s := parseSchedule(r)

	return s, nil
}

func parseSchedule(r *zip.ReadCloser) GTFSSchedule {
	var s GTFSSchedule

	for _, f := range r.File {
		spec := gtfsSpecs[f.Name]
		if spec == nil {
			s.unusedFiles = append(s.unusedFiles, f.Name)
			s.warnings.add(fmt.Errorf("unused file: %s", f.Name))
			continue
		}
		spec.parseFile(f, &s, &s.errors)
	}

	return s
}
