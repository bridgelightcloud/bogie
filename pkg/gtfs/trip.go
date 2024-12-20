package gtfs

import (
	"fmt"
)

type Trip struct {
	RouteID              string `json:"routeId,omitempty" csv:"route_id,omitempty"`
	ServiceID            string `json:"serviceId,omitempty" csv:"service_id,omitempty"`
	ID                   string `json:"tripId" csv:"trip_id"`
	Headsign             string `json:"tripHeadsign" csv:"trip_headsign"`
	ShortName            string `json:"tripShortName" csv:"trip_short_name"`
	DirectionID          int    `json:"directionId" csv:"direction_id"`
	BlockID              string `json:"blockId" csv:"block_id"`
	ShapeID              string `json:"shapeId" csv:"shape_id"`
	WheelchairAccessible int    `json:"wheelchairAccessible" csv:"wheelchair_accessible"`
	BikesAllowed         int    `json:"bikesAllowed" csv:"bikes_allowed"`
}

func (t Trip) key() string {
	return t.ID
}

func (t Trip) validate() errorList {
	var errs errorList

	if t.ID == "" {
		errs.add(fmt.Errorf("trip ID is required"))
	}
	if t.ServiceID == "" {
		errs.add(fmt.Errorf("trip service id is required"))
	}
	if t.ID == "" {
		errs.add(fmt.Errorf("trip ID is required"))
	}

	return errs
}
