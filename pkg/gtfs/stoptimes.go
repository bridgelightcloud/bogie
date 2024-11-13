package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

var (
	ErrEmptyStopTimesFile      = fmt.Errorf("empty stop times file")
	ErrInvalidStopTimesHeaders = fmt.Errorf("invalid stop times headers")
	ErrNoStopTimesRecords      = fmt.Errorf("no stop times records")
)

type StopTime struct {
	TripID                   string   `json:"tripId" csv:"trip_id"`
	ArrivalTime              string   `json:"arrivalTime,omitempty" csv:"arrival_time,omitempty"`
	DepartureTime            string   `json:"departureTime,omitempty" csv:"departure_time,omitempty"`
	StopID                   string   `json:"stopId" csv:"stop_id"`
	LocationGroupID          string   `json:"locationGroupId" csv:"location_group_id"`
	LocationID               string   `json:"locationId" csv:"location_id"`
	StopSequence             string   `json:"stopSequence" csv:"stop_sequence"`
	StopHeadsign             string   `json:"stopHeadsign" csv:"stop_headsign"`
	StartPickupDropOffWindow string   `json:"startPickupDropOffWindow" csv:"start_pickup_drop_off_window"`
	EndPickupDropOffWindow   string   `json:"endPickupDropOffWindow" csv:"end_pickup_drop_off_window"`
	PickupType               string   `json:"pickupType" csv:"pickup_type"`
	DropOffType              string   `json:"dropOffType" csv:"drop_off_type"`
	ContinuousPickup         string   `json:"continuousPickup" csv:"continuous_pickup"`
	ContinuousDropOff        string   `json:"continuousDropOff" csv:"continuous_drop_off"`
	ShapeDistTraveled        string   `json:"shapeDistTraveled" csv:"shape_dist_traveled"`
	Timepoint                string   `json:"timepoint" csv:"timepoint"`
	PickupBookingRuleId      string   `json:"pickupBookingRuleId" csv:"pickup_booking_rule_id"`
	DropOffBookingRuleId     string   `json:"dropOffBookingRuleId" csv:"drop_off_booking_rule_id"`
	Unused                   []string `json:"-" csv:"-"`
}

func parseStopTimes(file *zip.File) ([]StopTime, error) {
	rc, err := file.Open()
	if err != nil {
		return []StopTime{}, err
	}
	defer rc.Close()

	lines, err := csv.NewReader(rc).ReadAll()
	if len(lines) == 0 {
		return []StopTime{}, ErrEmptyStopTimesFile
	}

	headers := lines[0]
	if err := validateStopTimesHeader(headers); err != nil {
		return []StopTime{}, err
	}

	records := lines[1:]
	if len(records) == 0 {
		return []StopTime{}, ErrNoStopTimesRecords
	}

	stopTimes := make([]StopTime, len(records))
	for i, field := range headers {
		for j, record := range records {
			switch field {
			case "trip_id":
				stopTimes[j].TripID = record[i]
			case "arrival_time":
				stopTimes[j].ArrivalTime = record[i]
			case "departure_time":
				stopTimes[j].DepartureTime = record[i]
			case "stop_id":
				stopTimes[j].StopID = record[i]
			case "location_group_id":
				stopTimes[j].LocationGroupID = record[i]
			case "location_id":
				stopTimes[j].LocationID = record[i]
			case "stop_sequence":
				stopTimes[j].StopSequence = record[i]
			case "stop_headsign":
				stopTimes[j].StopHeadsign = record[i]
			case "start_pickup_drop_off_window":
				stopTimes[j].StartPickupDropOffWindow = record[i]
			case "end_pickup_drop_off_window":
				stopTimes[j].EndPickupDropOffWindow = record[i]
			case "pickup_type":
				stopTimes[j].PickupType = record[i]
			case "drop_off_type":
				stopTimes[j].DropOffType = record[i]
			case "continuous_pickup":
				stopTimes[j].ContinuousPickup = record[i]
			case "continuous_drop_off":
				stopTimes[j].ContinuousDropOff = record[i]
			case "shape_dist_traveled":
				stopTimes[j].ShapeDistTraveled = record[i]
			case "timepoint":
				stopTimes[j].Timepoint = record[i]
			case "pickup_booking_rule_id":
				stopTimes[j].PickupBookingRuleId = record[i]
			case "drop_off_booking_rule_id":
				stopTimes[j].DropOffBookingRuleId = record[i]
			default:
				stopTimes[j].Unused = append(stopTimes[j].Unused, record[i])
			}
		}
	}

	return stopTimes, nil
}

func validateStopTimesHeader(fields []string) error {
	requiredFields := []struct {
		name  string
		found bool
	}{{
		name:  "trip_id",
		found: false,
	}, {
		name:  "stop_sequence",
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
			return ErrInvalidStopTimesHeaders
		}
	}

	return nil
}
