package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
)

var (
	ErrEmptyStopsFile      = fmt.Errorf("empty stops file")
	ErrInvalidStopsHeaders = fmt.Errorf("invalid stops headers")
	ErrNoStopsRecords      = fmt.Errorf("no stops records")
)

var (
	ErrInvalidStopID                 = fmt.Errorf("invalid stop ID")
	ErrInvalidStopCode               = fmt.Errorf("invalid stop code")
	ErrInvalidStopName               = fmt.Errorf("invalid stop name")
	ErrInvalidStopTTSName            = fmt.Errorf("invalid stop TTS name")
	ErrInvalidStopDesc               = fmt.Errorf("invalid stop description")
	ErrInvalidStopLat                = fmt.Errorf("invalid stop latitude")
	ErrInvalidStopLon                = fmt.Errorf("invalid stop longitude")
	ErrInvalidStopZoneID             = fmt.Errorf("invalid stop zone ID")
	ErrInvalidStopURL                = fmt.Errorf("invalid stop URL")
	ErrInvalidStopLocationType       = fmt.Errorf("invalid stop location type")
	ErrInvalidStopParentStation      = fmt.Errorf("invalid stop parent station")
	ErrInvalidStopTimezone           = fmt.Errorf("invalid stop timezone")
	ErrInvalidStopWheelchairBoarding = fmt.Errorf("invalid stop wheelchair boarding")
	ErrInvalidStopLevelID            = fmt.Errorf("invalid stop level ID")
	ErrInvalidStopPlatformCode       = fmt.Errorf("invalid stop platform code")
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
}

func (s *GTFSSchedule) parseStopsData(file *zip.File) error {
	s.Stops = make(map[string]Stop)

	cp := make(map[string]string)
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		return ErrEmptyStopsFile
	}
	if err != nil {
		return err
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

		s.Stops[st.ID] = st

		if st.ParentStation != "" {
			cp[st.ID] = st.ParentStation
		}
	}

	if err != io.EOF {
		return err
	}

	if len(s.Stops) == 0 {
		return ErrNoStopsRecords
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
