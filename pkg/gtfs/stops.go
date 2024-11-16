package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
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

type StopLocationType int

var (
	StopLocationTypeStopPlatform StopLocationType = 0
	StopLocationTypeStation      StopLocationType = 1
	StopLocationTypeEntranceExit StopLocationType = 2
	StopLocationTypeGenericNode  StopLocationType = 3
	StopLocationTypeBoardingArea StopLocationType = 4
)

type (
	Latitude  float64
	Longitude float64
)

type Stop struct {
	ID                 string           `json:"stopId"`
	Code               string           `json:"stopCode,omitempty"`
	Name               string           `json:"stopName"`
	TTSName            string           `json:"TTSStopName,omitempty"`
	Desc               string           `json:"stopDesc,omitempty"`
	Lat                *Latitude        `json:"stopLat"`
	Lon                *Longitude       `json:"stopLon"`
	ZoneID             string           `json:"zoneId,omitempty"`
	URL                string           `json:"stopUrl,omitempty"`
	LocationType       StopLocationType `json:"locationType,omitempty"`
	ParentStation      string           `json:"parentStation"`
	Timezone           string           `json:"stopTimezone,omitempty"`
	WheelchairBoarding string           `json:"wheelchairBoarding,omitempty"`
	LevelID            string           `json:"levelId,omitempty"`
	PlatformCode       string           `json:"platformCode,omitempty"`
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
	for {
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

		var stop Stop
		for j, value := range record {
			value = strings.TrimSpace(value)
			switch headers[j] {
			case "stop_id":
				stop.ID = value
			case "stop_code":
				stop.Code = value
			case "stop_name":
				stop.Name = value
			case "tts_stop_name":
				stop.TTSName = value
			case "stop_desc":
				stop.Desc = value
			case "stop_lat":
				l, err := strconv.ParseFloat(value, 64)
				if err != nil {
					fmt.Printf("err: %s\n", err.Error())
					return ErrInvalidStopLat
				}
				p := Latitude(l)
				stop.Lat = &p
			case "stop_lon":
				l, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return ErrInvalidStopLon
				}
				p := Longitude(l)
				stop.Lon = &p
			case "zone_id":
				stop.ZoneID = value
			case "stop_url":
				stop.URL = value
			case "location_type":
				if value != "" {
					lt, err := strconv.Atoi(value)
					if err != nil {
						return ErrInvalidStopLocationType
					}
					stop.LocationType = StopLocationType(lt)
				}
			case "parent_station":
				stop.ParentStation = value
			case "stop_timezone":
				stop.Timezone = value
			case "wheelchair_boarding":
				stop.WheelchairBoarding = value
			case "level_id":
				stop.LevelID = value
			case "platform_code":
				stop.PlatformCode = value
			default:
				stop.unused = append(stop.unused, value)
			}
		}

		if err := stop.validateStop(); err != nil {
			return err
		}

		s.Stops[stop.ID] = stop

		if stop.ParentStation != "" {
			cp[stop.ID] = stop.ParentStation
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

func (s Stop) validateStop() error {
	if s.ID == "" {
		return ErrInvalidStopID
	}

	if s.Name == "" {
		rlt := map[StopLocationType]bool{
			StopLocationTypeStopPlatform: true,
			StopLocationTypeStation:      true,
			StopLocationTypeEntranceExit: true,
		}
		if _, ok := rlt[s.LocationType]; ok {
			fmt.Println(s)
			return fmt.Errorf("Invalid stop name \"%s\" for location type %d\n", s.Name, s.LocationType)
		}
	}

	if s.Lat == nil {
		rlt := map[StopLocationType]bool{
			StopLocationTypeStopPlatform: true,
			StopLocationTypeStation:      true,
			StopLocationTypeEntranceExit: true,
		}
		if _, ok := rlt[s.LocationType]; ok {
			fmt.Printf("invalid latitude %f for location type %d\n", *s.Lat, *&s.LocationType)
			return ErrInvalidStopLat
		}
	}

	if s.Lon == nil {
		rlt := map[StopLocationType]bool{
			StopLocationTypeStopPlatform: true,
			StopLocationTypeStation:      true,
			StopLocationTypeEntranceExit: true,
		}
		if _, ok := rlt[s.LocationType]; ok {
			return ErrInvalidStopLon
		}
	}

	if s.ParentStation == "" {
		rlt := map[StopLocationType]bool{
			StopLocationTypeEntranceExit: true,
			StopLocationTypeGenericNode:  true,
			StopLocationTypeBoardingArea: true,
		}
		if _, ok := rlt[s.LocationType]; ok {
			return ErrInvalidStopParentStation
		}
	} else {
		if s.LocationType == StopLocationTypeStation {
			return ErrInvalidStopParentStation
		}
	}

	return nil
}
