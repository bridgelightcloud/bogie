package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
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
	DirectionID          int      `json:"directionId" csv:"direction_id"`
	BlockID              string   `json:"blockId" csv:"block_id"`
	ShapeID              string   `json:"shapeId" csv:"shape_id"`
	WheelchairAccessible int      `json:"wheelchairAccessible" csv:"wheelchair_accessible"`
	BikesAllowed         int      `json:"bikesAllowed" csv:"bikes_allowed"`
	unused               []string `json:"-" csv:"-"`
}

func (s *GTFSSchedule) parseTrips(file *zip.File) error {
	s.Trips = map[string]Trip{}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors.add(ErrEmptyTripsFile)
		return ErrEmptyTripsFile
	}

	if err != nil {
		s.errors.add(err)
		return err
	}

	var record []string
	for i := 0; ; i++ {
		record, err = r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			s.errors.add(fmt.Errorf("empty record at line %d", i))
			continue
		}

		t := Trip{}
		for j, v := range record {
			switch headers[j] {
			case "route_id":
				ParseString(v, &t.RouteID)
			case "service_id":
				ParseString(v, &t.ServiceID)
			case "trip_id":
				ParseString(v, &t.ID)
			case "trip_headsign":
				ParseString(v, &t.Headsign)
			case "trip_short_name":
				ParseString(v, &t.ShortName)
			case "direction_id":
				if err := ParseEnum(v, DirectionID, &t.DirectionID); err != nil {
					s.errors.add(fmt.Errorf("invalid direction id at line %d: %w", i, err))
				}
			case "block_id":
				ParseString(v, &t.BlockID)
			case "shape_id":
				ParseString(v, &t.ShapeID)
			case "wheelchair_accessible":
				if err := ParseEnum(v, WheelchairAccessible, &t.WheelchairAccessible); err != nil {
					s.errors.add(fmt.Errorf("invalid wheelchair accessible at line %d: %w", i, err))
				}
			case "bikes_allowed":
				if err := ParseEnum(v, BikesAllowed, &t.BikesAllowed); err != nil {
					s.errors.add(fmt.Errorf("invalid bikes allowed at line %d: %w", i, err))
				}
			default:
				t.unused = append(t.unused, v)
			}
		}
		s.Trips[t.ID] = t
	}

	return nil
}
