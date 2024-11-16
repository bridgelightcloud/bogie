package gtfs

import (
	"fmt"

	"github.com/google/uuid"
)

func Overview(c map[string]GTFSSchedule) string {
	o := ""

	for sid, s := range c {
		o += fmt.Sprintf("Schedule %s\n", sid[0:4])
		o += fmt.Sprintf("    %d agencies\n", len(s.Agencies))
		o += fmt.Sprintf("    %d stops\n", len(s.Stops))
		o += fmt.Sprintf("    %d routes\n", len(s.Routes))
		o += "\n"
	}

	return o
}

func CreateGTFSCollection(zipFiles []string) (map[string]GTFSSchedule, error) {
	sc := make(map[string]GTFSSchedule)

	for _, path := range zipFiles {
		s, err := OpenScheduleFromFile(path)
		if err != nil {
			return sc, err
		}

		sc[uuid.NewString()] = s
	}

	return sc, nil
}
