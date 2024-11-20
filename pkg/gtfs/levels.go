package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"math"
)

type Level struct {
	ID    string  `json:"levelId"`
	Index float64 `json:"levelIndex"`
	Name  string  `json:"levelName,omitempty"`

	unused   []string
	errors   errorList
	warnings errorList
}

func (s *GTFSSchedule) parseLevels(file *zip.File) {
	s.Levels = make(map[string]Level)

	rc, err := file.Open()
	if err != nil {
		s.errors.add(fmt.Errorf("error opening levels file: %w", err))
		return
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors.add(fmt.Errorf("empty levels file"))
		return
	}
	if err != nil {
		s.errors.add(err)
		return
	}

	var record []string
	for i := 0; ; i++ {
		record, err = r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			continue
		}

		var l Level
		for j, v := range record {
			switch headers[j] {
			case "level_id":
				ParseString(v, &l.ID)
			case "level_index":
				ParseFloat(v, &l.Index)
			case "level_name":
				ParseString(v, &l.Name)
			default:
				l.unused = append(l.unused, headers[j])
			}
		}

		validateLevel(l)
		s.Levels[l.ID] = l
	}
}

func validateLevel(l Level) {
	if l.ID == "" {
		l.errors.add(fmt.Errorf("missing level_id"))
	}

	if l.Index == math.Inf(-1) {
		l.errors.add(fmt.Errorf("invalid index valie"))
	}

	// Name is optional
}
