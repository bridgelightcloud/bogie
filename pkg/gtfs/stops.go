package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
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
	ID                 string   `json:"stopId,omitempty" csv:"stop_id,omitempty"`
	Code               string   `json:"stopCode,omitempty" csv:"stop_code,omitempty"`
	Name               string   `json:"stopName" csv:"stop_name"`
	TTSName            string   `json:"TTSStopName" csv:"tts_stop_name"`
	Desc               string   `json:"stopDesc" csv:"stop_desc"`
	Lat                string   `json:"stopLat" csv:"stop_lat"`
	Lon                string   `json:"stopLon" csv:"stop_lon"`
	ZoneID             string   `json:"zoneId" csv:"zone_id"`
	URL                string   `json:"stopUrl" csv:"stop_url"`
	LocationType       string   `json:"locationType" csv:"location_type"`
	ParentStation      string   `json:"parentStation" csv:"parent_station"`
	Timezone           string   `json:"stopTimezone" csv:"stop_timezone"`
	WheelchairBoarding string   `json:"wheelchairBoarding" csv:"wheelchair_boarding"`
	LevelID            string   `json:"levelId" csv:"level_id"`
	PlatformCode       string   `json:"platformCode" csv:"platform_code"`
	Unused             []string `json:"-" csv:"-"`
}

func parseStops(file *zip.File) ([]Stop, error) {
	rc, err := file.Open()
	if err != nil {
		return []Stop{}, err
	}
	defer rc.Close()

	lines, err := csv.NewReader(rc).ReadAll()
	if len(lines) == 0 {
		return []Stop{}, ErrEmptyStopsFile
	}

	headers := lines[0]
	if err := validateStopsHeader(headers); err != nil {
		return []Stop{}, err
	}

	records := lines[1:]
	if len(records) == 0 {
		return []Stop{}, ErrNoStopsRecords
	}

	stops := make([]Stop, len(records))
	for i, field := range headers {
		for j, record := range records {
			switch field {
			case "stop_id":
				stops[j].ID = record[i]
			case "stop_code":
				stops[j].Code = record[i]
			case "stop_name":
				stops[j].Name = record[i]
			case "tts_stop_name":
				stops[j].TTSName = record[i]
			case "stop_desc":
				stops[j].Desc = record[i]
			case "stop_lat":
				stops[j].Lat = record[i]
			case "stop_lon":
				stops[j].Lon = record[i]
			case "zone_id":
				stops[j].ZoneID = record[i]
			case "stop_url":
				stops[j].URL = record[i]
			case "location_type":
				stops[j].LocationType = record[i]
			case "parent_station":
				stops[j].ParentStation = record[i]
			case "stop_timezone":
				stops[j].Timezone = record[i]
			case "wheelchair_boarding":
				stops[j].WheelchairBoarding = record[i]
			case "level_id":
				stops[j].LevelID = record[i]
			case "platform_code":
				stops[j].PlatformCode = record[i]
			default:
				stops[j].Unused = append(stops[j].Unused, record[i])
			}
		}
	}

	return stops, nil
}

func validateStopsHeader(fields []string) error {
	requiredFields := []struct {
		name  string
		found bool
	}{{
		name:  "stop_id",
		found: false,
	}}

	for _, field := range fields {
		for i, f := range requiredFields {
			if field == f.name {
				requiredFields[i].found = true
			}
		}
	}

	for _, f := range requiredFields {
		if !f.found {
			return ErrInvalidStopsHeaders
		}
	}

	return nil
}

func (s Stop) validateStop() {

}

func buildStopHierarchy(stops []Stop) map[string][]Stop {
	hierarchy := make(map[string][]Stop)
	for _, stop := range stops {
		if stop.ParentStation != "" {
			if _, ok := hierarchy[stop.ParentStation]; !ok {
				hierarchy[stop.ParentStation] = []Stop{}
			}
			hierarchy[stop.ParentStation] = append(hierarchy[stop.ParentStation], stop)
		}
	}
	return hierarchy
}
