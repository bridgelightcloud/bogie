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
	DirectionID          Enum     `json:"directionId" csv:"direction_id"`
	BlockID              string   `json:"blockId" csv:"block_id"`
	ShapeID              string   `json:"shapeId" csv:"shape_id"`
	WheelchairAccessible Enum     `json:"wheelchairAccessible" csv:"wheelchair_accessible"`
	BikesAllowed         string   `json:"bikesAllowed" csv:"bikes_allowed"`
	Unused               []string `json:"-" csv:"-"`
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
		s.errors = append(s.errors, ErrEmptyTripsFile)
		return ErrEmptyTripsFile
	}

	if err != nil {
		s.errors = append(s.errors, err)
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
				t.RouteID = v
			case "service_id":
				t.ServiceID = v
			case "trip_id":
				t.ID = v
			case "trip_headsign":
				t.Headsign = v
			case "trip_short_name":
				t.ShortName = v
			case "direction_id":
				if err := t.DirectionID.Parse(v, Availability); err != nil {
					s.errors.add(fmt.Errorf("invalid direction_id at line %d: %w", i, err))
				}
			case "block_id":
				t.BlockID = v
			case "shape_id":
				t.ShapeID = v
			case "wheelchair_accessible":
				t.WheelchairAccessible = 0
			case "bikes_allowed":
				t.BikesAllowed = v
			default:
				t.Unused = append(t.Unused, v)
			}
		}
		s.Trips[t.ID] = t
	}

	return nil
}
