package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

var (
	ErrEmptyTripsFile      = fmt.Errorf("empty trips file")
	ErrInvalidTripsHeaders = fmt.Errorf("invalid trips headers")
	ErrNoTripsRecords      = fmt.Errorf("no trips records")
)

type Trip struct {
	RouteID              string   `json:"routeId,omitempty" csv:"route_id,omitempty"`
	ServiceID            string   `json:"serviceId,omitempty" csv:"service_id,omitempty"`
	ID                   string   `json:"tripId" csv:"trip_id"`
	Headsign             string   `json:"tripHeadsign" csv:"trip_headsign"`
	ShortName            string   `json:"tripShortName" csv:"trip_short_name"`
	DirectionID          string   `json:"directionId" csv:"direction_id"`
	BlockID              string   `json:"blockId" csv:"block_id"`
	ShapeID              string   `json:"shapeId" csv:"shape_id"`
	WheelchairAccessible string   `json:"wheelchairAccessible" csv:"wheelchair_accessible"`
	BikesAllowed         string   `json:"bikesAllowed" csv:"bikes_allowed"`
	Unused               []string `json:"-" csv:"-"`
}

func parseTrips(file *zip.File) ([]Trip, error) {
	rc, err := file.Open()
	if err != nil {
		return []Trip{}, err
	}
	defer rc.Close()

	lines, err := csv.NewReader(rc).ReadAll()
	if len(lines) == 0 {
		return []Trip{}, ErrEmptyTripsFile
	}

	headers := lines[0]
	if err := validateTripsHeader(headers); err != nil {
		return []Trip{}, err
	}

	records := lines[1:]
	if len(records) == 0 {
		return []Trip{}, ErrNoTripsRecords
	}

	trips := make([]Trip, len(records))
	for i, field := range headers {
		for j, record := range records {
			switch field {
			case "trip_id":
				trips[j].ID = record[i]
			case "route_id":
				trips[j].RouteID = record[i]
			case "service_id":
				trips[j].ServiceID = record[i]
			case "trip_headsign":
				trips[j].Headsign = record[i]
			case "trip_short_name":
				trips[j].ShortName = record[i]
			case "direction_id":
				trips[j].DirectionID = record[i]
			case "block_id":
				trips[j].BlockID = record[i]
			case "shape_id":
				trips[j].ShapeID = record[i]
			case "wheelchair_accessible":
				trips[j].WheelchairAccessible = record[i]
			case "bikes_allowed":
				trips[j].BikesAllowed = record[i]
			default:
				trips[j].Unused = append(trips[j].Unused, record[i])
			}
		}
	}

	return trips, nil
}

func validateTripsHeader(headers []string) error {
	requiredFields := []struct {
		name  string
		found bool
	}{{
		name:  "route_id",
		found: false,
	}, {
		name:  "service_id",
		found: false,
	}, {
		name:  "trip_id",
		found: false,
	}}

	for _, field := range headers {
		for i, rf := range requiredFields {
			if field == rf.name {
				requiredFields[i].found = true
			}
		}
	}

	for _, rf := range requiredFields {
		if !rf.found {
			return ErrInvalidTripsHeaders
		}
	}

	return nil
}
