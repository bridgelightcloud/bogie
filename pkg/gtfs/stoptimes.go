package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

type StopTime struct {
	TripID                   string    `json:"tripId"`
	ArrivalTime              time.Time `json:"arrivalTime,omitempty"`
	DepartureTime            time.Time `json:"departureTime,omitempty"`
	StopID                   string    `json:"stopId"`
	LocationGroupID          string    `json:"locationGroupId"`
	LocationID               string    `json:"locationId"`
	StopSequence             int       `json:"stopSequence"`
	StopHeadsign             string    `json:"stopHeadsign"`
	StartPickupDropOffWindow time.Time `json:"startPickupDropOffWindow"`
	EndPickupDropOffWindow   time.Time `json:"endPickupDropOffWindow"`
	PickupType               int       `json:"pickupType"`
	DropOffType              int       `json:"dropOffType"`
	ContinuousPickup         int       `json:"continuousPickup"`
	ContinuousDropOff        int       `json:"continuousDropOff"`
	ShapeDistTraveled        float64   `json:"shapeDistTraveled"`
	Timepoint                int       `json:"timepoint"`
	PickupBookingRuleId      string    `json:"pickupBookingRuleId"`
	DropOffBookingRuleId     string    `json:"dropOffBookingRuleId"`

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
				ParseString(v, &st.TripID)
			case "arrival_time":
				t, err := time.Parse(timeFormat, v)
				if err != nil {
					s.errors.add(fmt.Errorf("invalid arrival time at line %d: %w", i, err))
				}
				st.ArrivalTime = t
			case "departure_time":
				t, err := time.Parse(timeFormat, v)
				if err != nil {
					s.errors.add(fmt.Errorf("invalid departure time at line %d: %w", i, err))
				}
				st.DepartureTime = t
			case "stop_id":
				ParseString(v, &st.StopID)
			case "location_group_id":
				ParseString(v, &st.LocationGroupID)
			case "location_id":
				ParseString(v, &st.LocationID)
			case "stop_sequence":
				if err := ParseInt(v, &st.StopSequence); err != nil {
					s.errors.add(fmt.Errorf("invalid stop sequence at line %d: %w", i, err))
				}
			case "stop_headsign":
				ParseString(v, &st.StopHeadsign)
			case "start_pickup_drop_off_window":
				t, err := time.Parse(timeFormat, v)
				if err != nil {
					s.errors.add(fmt.Errorf("invalid start pickup drop off window at line %d: %w", i, err))
				}
				st.StartPickupDropOffWindow = t
			case "end_pickup_drop_off_window":
				t, err := time.Parse(timeFormat, v)
				if err != nil {
					s.errors.add(fmt.Errorf("invalid end pickup drop off window at line %d: %w", i, err))
				}
				st.EndPickupDropOffWindow = t
			case "pickup_type":
				if err := ParseEnum(v, PickupType, &st.PickupType); err != nil {
					s.errors.add(fmt.Errorf("invalid pickup type at line %d: %w", i, err))
				}
			case "drop_off_type":
				if err := ParseEnum(v, DropOffType, &st.DropOffType); err != nil {
					s.errors.add(fmt.Errorf("invalid drop off type at line %d: %w", i, err))
				}
			case "continuous_pickup":
				if err := ParseEnum(v, ContinuousPickup, &st.ContinuousPickup); err != nil {
					s.errors.add(fmt.Errorf("invalid continuous pickup at line %d: %w", i, err))
				}
			case "continuous_drop_off":
				if err := ParseEnum(v, ContinuousDropOff, &st.ContinuousDropOff); err != nil {
					s.errors.add(fmt.Errorf("invalid continuous drop off at line %d: %w", i, err))
				}
			case "shape_dist_traveled":
				if err := ParseFloat(v, &st.ShapeDistTraveled); err != nil {
					s.errors.add(fmt.Errorf("invalid shape dist traveled at line %d: %w", i, err))
				}
			case "timepoint":
				if err := ParseEnum(v, Timepoint, &st.Timepoint); err != nil {
					s.errors.add(fmt.Errorf("invalid timepoint at line %d: %w", i, err))
				}
			case "pickup_booking_rule_id":
				ParseString(v, &st.PickupBookingRuleId)
			case "drop_off_booking_rule_id":
				ParseString(v, &st.DropOffBookingRuleId)
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
		s.errors.add(fmt.Errorf("no stop times found"))
	}
	return nil
}
