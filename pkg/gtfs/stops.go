package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
)

type Stop struct {
	ID                 string `json:"stopId"`
	Code               string `json:"stopCode,omitempty"`
	Name               string `json:"stopName"`
	TTSName            string `json:"TTSStopName,omitempty"`
	Desc               string `json:"stopDesc,omitempty"`
	Coords             Coords `json:"coords"`
	ZoneID             string `json:"zoneId,omitempty"`
	URL                string `json:"stopUrl,omitempty"`
	LocationType       int    `json:"locationType,omitempty"`
	ParentStation      string `json:"parentStation"`
	Timezone           string `json:"stopTimezone,omitempty"`
	WheelchairBoarding string `json:"wheelchairBoarding,omitempty"`
	LevelID            string `json:"levelId,omitempty"`
	PlatformCode       string `json:"platformCode,omitempty"`
	unused             []string

	children map[string]bool
	errors   errorList
	warnings errorList
}

func (s *GTFSSchedule) parseStopsData(file *zip.File) error {
	s.Stops = make(map[string]Stop)

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors.add(fmt.Errorf("empty stops file"))
	}
	if err != nil {
		return err
	}

	cp := make(map[string]string)
	var record []string
	for i := 0; ; i++ {
		record, err = r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			continue
		}

		if len(record) > len(headers) {
			return fmt.Errorf("record has too many columns")
		}

		var st Stop
		for j, value := range record {
			switch headers[j] {
			case "stop_id":
				ParseString(value, &st.ID)
			case "stop_code":
				ParseString(value, &st.Code)
			case "stop_name":
				ParseString(value, &st.Name)
			case "tts_stop_name":
				ParseString(value, &st.TTSName)
			case "stop_desc":
				ParseString(value, &st.Desc)
			case "stop_lat":
				ParseLat(value, &st.Coords)
			case "stop_lon":
				ParseLon(value, &st.Coords)
			case "zone_id":
				ParseString(value, &st.ZoneID)
			case "stop_url":
				ParseString(value, &st.URL)
			case "location_type":
				if err := ParseEnum(value, LocationType, &st.LocationType); err != nil {
					return fmt.Errorf("invalid location_type at line %d: %w", i, err)
				}
			case "parent_station":
				ParseString(value, &st.ParentStation)
			case "stop_timezone":
				ParseString(value, &st.Timezone)
			case "wheelchair_boarding":
				ParseString(value, &st.WheelchairBoarding)
			case "level_id":
				ParseString(value, &st.LevelID)
			case "platform_code":
				ParseString(value, &st.PlatformCode)
			default:
				st.unused = append(st.unused, value)
			}
		}
		validateStop(&st)
		s.Stops[st.ID] = st

		if st.ParentStation != "" {
			cp[st.ID] = st.ParentStation
		}
	}

	if err != io.EOF {
		return err
	}

	if len(s.Stops) == 0 {
		s.errors.add(fmt.Errorf("no stop records found"))
	}

	for id, parentId := range cp {
		if p, ok := s.Stops[parentId]; ok {
			if p.children == nil {
				p.children = make(map[string]bool)
			}
			p.children[id] = true
		} else {
			return fmt.Errorf("Parent stop %s for stop %s not found", parentId, id)
		}
	}

	return nil
}

func validateStop(st *Stop) {
	if st.ID == "" {
		st.errors.add(fmt.Errorf("stop ID is required"))
	}

	// Code is optional

	if st.Name == "" {
		if st.LocationType == StopPlatform || st.LocationType == Station || st.LocationType == EntranceExit {
			st.errors.add(fmt.Errorf("stop name is required for location type %d", st.LocationType))
		}
	}

	// TTSName is optional

	// Desc is optional

	if !st.Coords.IsValid() {
		if st.LocationType == StopPlatform || st.LocationType == Station || st.LocationType == EntranceExit {
			st.errors.add(fmt.Errorf("invalid stop coordinates for location type %d", st.LocationType))
		}
	}

	// ZoneID is optional

	if st.LocationType < StopPlatform || st.LocationType > BoardingArea {
		st.errors.add(fmt.Errorf("invalid location type: %d", st.LocationType))
	}

	// ParentStation is validated in full stops list

	// Validate Timezone

	// WheelchairBoarding is validated in full stops list

	// LevelID is validated in full stops list

	// PlatformCode is optional
}
