package gtfs

import (
	"fmt"

	"github.com/google/uuid"
)

func Overview(c map[string]GTFSSchedule) string {
	var o string

	for sid, s := range c {
		o += fmt.Sprintf("Schedule %s\n", sid[0:4])
		o += fmt.Sprintf("  %d agencies\n", len(s.Agencies))
		o += fmt.Sprintf("  %d stops\n", len(s.Stops))
		o += fmt.Sprintf("  %d routes\n", len(s.Routes))
		o += fmt.Sprintf("  %d calendar entries\n", len(s.Calendar))
		o += fmt.Sprintf("  %d calendar dates\n", len(s.CalendarDates))
		o += fmt.Sprintf("  %d trips\n", len(s.Trips))
		o += fmt.Sprintf("  %d stop times\n", len(s.StopTimes))
		o += fmt.Sprintf("  %d levels\n", len(s.Levels))
		o += fmt.Sprintf("  %d errors\n", len(s.errors))
		o += "\n"
	}

	return o
}

func CreateGTFSCollection(zipFiles []string) (map[string]GTFSSchedule, error) {
	sc := make(map[string]GTFSSchedule)

	for _, path := range zipFiles {
		s, err := OpenScheduleFromZipFile(path)
		if err != nil {
			return sc, err
		}

		sc[uuid.NewString()] = s
	}

	return sc, nil
}
