package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
)

var (
	ErrEmptyStopTimesFile      = fmt.Errorf("empty stop times file")
	ErrInvalidStopTimesHeaders = fmt.Errorf("invalid stop times headers")
	ErrNoStopTimesRecords      = fmt.Errorf("no stop times records")
)

type StopTime struct {
	TripID                   string  `json:"tripId"`
	ArrivalTime              Time    `json:"arrivalTime,omitempty"`
	DepartureTime            Time    `json:"departureTime,omitempty"`
	StopID                   string  `json:"stopId"`
	LocationGroupID          string  `json:"locationGroupId"`
	LocationID               string  `json:"locationId"`
	StopSequence             Int     `json:"stopSequence"`
	StopHeadsign             string  `json:"stopHeadsign"`
	StartPickupDropOffWindow Time    `json:"startPickupDropOffWindow"`
	EndPickupDropOffWindow   Time    `json:"endPickupDropOffWindow"`
	PickupType               Enum    `json:"pickupType"`
	DropOffType              Enum    `json:"dropOffType"`
	ContinuousPickup         Enum    `json:"continuousPickup"`
	ContinuousDropOff        Enum    `json:"continuousDropOff"`
	ShapeDistTraveled        Float64 `json:"shapeDistTraveled"`
	Timepoint                Enum    `json:"timepoint"`
	PickupBookingRuleId      string  `json:"pickupBookingRuleId"`
	DropOffBookingRuleId     string  `json:"dropOffBookingRuleId"`

	primaryKey string
	unused     []string
}

func (s *GTFSSchedule) parseStopTimes(file *zip.File) error {
	s.StopTimes = map[string]StopTime{}

	rc, err := file.Open()
	if err != nil {
		return s.errors.add(fmt.Errorf("error opening stop times file: %w", err))
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		return s.errors.add(fmt.Errorf("empty stop times file"))
	}
	if err != nil {
		return s.errors.add(fmt.Errorf("error reading stop times headers: %w", err))
	}

	record := []string{}
	for i := 0; ; i++ {
		record, err = r.Read()
		if err != nil {
			break
		}

		if len(record) == 0 {
			s.errors.add(fmt.Errorf("empty record at line %d", i))
			continue
		}

		if len(record) > len(headers) {
			s.errors.add(fmt.Errorf("invalid record at line %d: %v", i, record))
			continue
		}

		var st StopTime
		for j, v := range record {
			switch headers[j] {
			case "trip_id":
				st.TripID = v
			case "arrival_time":
				if err := st.ArrivalTime.parse(v); err != nil {
					s.errors.add(fmt.Errorf("invalid arrival time at line %d: %w", i, err))
				}
			case "departure_time":
				if err := st.DepartureTime.parse(v); err != nil {
					s.errors.add(fmt.Errorf("invalid departure time at line %d: %w", i, err))
				}
			case "stop_id":
				st.StopID = v
			case "location_group_id":
				st.LocationGroupID = v
			case "location_id":
				st.LocationID = v
			case "stop_sequence":
				if err := st.StopSequence.Parse(v); err != nil {
					s.errors.add(fmt.Errorf("invalid stop sequence at line %d: %w", i, err))
				}
			case "stop_headsign":
				st.StopHeadsign = v
			case "start_pickup_drop_off_window":
				if err := st.StartPickupDropOffWindow.parse(v); err != nil {
					s.errors.add(fmt.Errorf("invalid start pickup drop off window at line %d: %w", i, err))
				}
			case "end_pickup_drop_off_window":
				if err := st.EndPickupDropOffWindow.parse(v); err != nil {
					s.errors.add(fmt.Errorf("invalid end pickup drop off window at line %d: %w", i, err))
				}
			case "pickup_type":
				if err := st.PickupType.Parse(v, PickupType); err != nil {
					s.errors.add(fmt.Errorf("invalid pickup type at line %d: %w", i, err))
				}
			case "drop_off_type":
				if err := st.DropOffType.Parse(v, DropOffType); err != nil {
					s.errors.add(fmt.Errorf("invalid drop off type at line %d: %w", i, err))
				}
			case "continuous_pickup":
				if err := st.ContinuousPickup.Parse(v, ContinuousPickup); err != nil {
					s.errors.add(fmt.Errorf("invalid continuous pickup at line %d: %w", i, err))
				}
			case "continuous_drop_off":
				if err := st.ContinuousDropOff.Parse(v, ContinuousDropOff); err != nil {
					s.errors.add(fmt.Errorf("invalid continuous drop off at line %d: %w", i, err))
				}
			case "shape_dist_traveled":
				if err := st.ShapeDistTraveled.Parse(v); err != nil {
					s.errors.add(fmt.Errorf("invalid shape dist traveled at line %d: %w", i, err))
				}
			case "timepoint":
				if err := st.Timepoint.Parse(v, Timepoint); err != nil {
					s.errors.add(fmt.Errorf("invalid timepoint at line %d: %w", i, err))
				}
			case "pickup_booking_rule_id":
				st.PickupBookingRuleId = v
			case "drop_off_booking_rule_id":
				st.DropOffBookingRuleId = v
			default:
				st.unused = append(st.unused, v)
			}
		}
		primaryKey := fmt.Sprintf("%s.%d", st.TripID, st.StopSequence)
		if _, ok := s.StopTimes[primaryKey]; ok {
			fmt.Println(s.errors.add(fmt.Errorf("duplicate stop time record at line %d", i)))
		}
		s.StopTimes[primaryKey] = st
	}

	if err != io.EOF {
		s.errors.add(fmt.Errorf("error reading stop times file: %w", err))
	}

	if len(s.StopTimes) == 0 {
		s.errors.add(ErrNoStopTimesRecords)
	}
	return nil
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
