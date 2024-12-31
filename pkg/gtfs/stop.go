package gtfs

import (
	"fmt"
)

type Stop struct {
	ID                 string `json:"stopId" csv:"stop_id"`
	Code               string `json:"stopCode,omitempty" csv:"stop_code"`
	Name               string `json:"stopName" csv:"stop_name"`
	TTSName            string `json:"TTSStopName,omitempty" csv:"tts_stop_name"`
	Desc               string `json:"stopDesc,omitempty" csv:"stop_desc"`
	Latitude           string `json:"latitude" csv:"stop_lat"`
	Longitude          string `json:"longitude" csv:"stop_lon"`
	ZoneID             string `json:"zoneId,omitempty" csv:"zone_id"`
	URL                string `json:"stopUrl,omitempty" csv:"stop_url"`
	LocationType       int    `json:"locationType,omitempty" csv:"location_type"`
	ParentStation      string `json:"parentStation" csv:"parent_station"`
	Timezone           string `json:"stopTimezone,omitempty" csv:"stop_timezone"`
	WheelchairBoarding string `json:"wheelchairBoarding,omitempty" csv:"wheelchair_boarding"`
	LevelID            string `json:"levelId,omitempty" csv:"level_id"`
	PlatformCode       string `json:"platformCode,omitempty" csv:"platform_code"`
}

func (s Stop) key() string {
	return s.ID
}

func (s Stop) validate() errorList {
	var errs errorList

	if s.ID == "" {
		errs.add(fmt.Errorf("stop ID is required"))
	}
	if s.Name == "" {
		if s.LocationType == StopPlatform || s.LocationType == Station || s.LocationType == EntranceExit {
			errs.add(fmt.Errorf("stop name is required for location type %d", s.LocationType))
		}
	}
	// if !s.Coords.IsValid() {
	// 	if s.LocationType == StopPlatform || s.LocationType == Station || s.LocationType == EntranceExit {
	// 		errs.add(fmt.Errorf("invalid stop coordinates for location type %d", s.LocationType))
	// 	}
	// }
	if s.LocationType < StopPlatform || s.LocationType > BoardingArea {
		errs.add(fmt.Errorf("invalid location type: %d", s.LocationType))
	}

	return errs
}
