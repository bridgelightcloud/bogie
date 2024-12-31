package gtfs

import (
	"fmt"
)

type StopTime struct {
	TripID                   string   `json:"tripId" csv:"trip_id"`
	ArrivalTime              Time     `json:"arrivalTime,omitempty" csv:"arrival_time"`
	DepartureTime            Time     `json:"departureTime,omitempty" csv:"departure_time"`
	StopID                   string   `json:"stopId" csv:"stop_id"`
	LocationGroupID          string   `json:"locationGroupId" csv:"location_group_id"`
	LocationID               string   `json:"locationId" csv:"location_id"`
	StopSequence             int      `json:"stopSequence" csv:"stop_sequence"`
	StopHeadsign             string   `json:"stopHeadsign" csv:"stop_headsign"`
	StartPickupDropOffWindow Time     `json:"startPickupDropOffWindow" csv:"start_pickup_drop_off_window"`
	EndPickupDropOffWindow   Time     `json:"endPickupDropOffWindow" csv:"end_pickup_drop_off_window"`
	PickupType               *int     `json:"pickupType" csv:"pickup_type"`
	DropOffType              *int     `json:"dropOffType" csv:"drop_off_type"`
	ContinuousPickup         *int     `json:"continuousPickup" csv:"continuous_pickup"`
	ContinuousDropOff        *int     `json:"continuousDropOff" csv:"continuous_drop_off"`
	ShapeDistTraveled        *float64 `json:"shapeDistTraveled" csv:"shape_dist_traveled"`
	Timepoint                *int     `json:"timepoint" csv:"timepoint"`
	PickupBookingRuleId      string   `json:"pickupBookingRuleId" csv:"pickup_booking_rule_id"`
	DropOffBookingRuleId     string   `json:"dropOffBookingRuleId" csv:"drop_off_booking_rule_id"`
}

func (st StopTime) key() string {
	return fmt.Sprintf("%s-%d", st.TripID, st.StopSequence)
}

func (st StopTime) validate() errorList {
	var errs errorList

	if st.TripID == "" {
		errs.add(fmt.Errorf("trip ID is required"))
	}
	if st.StopSequence < 0 {
		errs.add(fmt.Errorf("stop sequence must be greater than or equal to 0"))
	}
	if st.StopID == "" {
		errs.add(fmt.Errorf("stop ID is required"))
	}

	return errs
}
