package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
)

type Route struct {
	ID                string `json:"routeId"`
	AgencyID          string `json:"agencyId"`
	ShortName         string `json:"routeShortName" csv:"route_short_name"`
	LongName          string `json:"routeLongName" csv:"route_long_name"`
	Desc              string `json:"routeDesc,omitempty"`
	Type              string `json:"routeType"`
	URL               string `json:"routeUrl,omitempty"`
	Color             string `json:"routeColor,omitempty"`
	TextColor         string `json:"routeTextColor,omitempty"`
	SortOrder         string `json:"routeSortOrder,omitempty"`
	ContinuousPickup  string `json:"continuousPickup,omitempty"`
	ContinuousDropOff string `json:"continuousDropOff,omitempty"`
	NetworkID         string `json:"networkId,omitempty"`
	unused            []string
}

func (s *GTFSSchedule) parseRoutes(file *zip.File) {
	s.Routes = map[string]Route{}

	if s.Agencies == nil {
		s.errors.add(fmt.Errorf("Agencies must be parsed before Routes"))
	}

	rc, err := file.Open()
	if err != nil {
		s.errors.add(fmt.Errorf("error opening routes file: %w", err))
		return
	}
	defer rc.Close()

	r := csv.NewReader(rc)

	headers, err := r.Read()
	if err == io.EOF {
		s.errors.add(fmt.Errorf("empty routes file"))
		return
	}
	if err != nil {
		s.errors.add(fmt.Errorf("error reading routes file: %w", err))
		return
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
			s.errors.add(fmt.Errorf("record has too many columns"))
			continue
		}

		var route Route
		for j, value := range record {
			switch headers[j] {
			case "route_id":
				route.ID = value
			case "agency_id":
				route.AgencyID = value
			case "route_short_name":
				route.ShortName = value
			case "route_long_name":
				route.LongName = value
			case "route_desc":
				route.Desc = value
			case "route_type":
				route.Type = value
			case "route_url":
				route.URL = value
			case "route_color":
				route.Color = value
			case "route_text_color":
				route.TextColor = value
			case "route_sort_order":
				route.SortOrder = value
			case "continuous_pickup":
				route.ContinuousPickup = value
			case "continuous_drop_off":
				route.ContinuousDropOff = value
			case "network_id":
				route.NetworkID = value
			default:
				route.unused = append(route.unused, value)
			}
			s.Routes[route.ID] = route
		}
	}

	if err != io.EOF {
		s.errors.add(fmt.Errorf("error reading routes file: %w", err))
	}
}

func validateRoute(r Route) error {
	if r.ID == "" {
		return fmt.Errorf("route ID is required")
	}
	if r.AgencyID == "" {
		return fmt.Errorf("route agency ID is required")
	}
	if r.ShortName == "" {
		return fmt.Errorf("route short name is required")
	}
	if r.LongName == "" {
		return fmt.Errorf("route long name is required")
	}
	if r.Type == "" {
		return fmt.Errorf("route type is required")
	}

	return nil
}
