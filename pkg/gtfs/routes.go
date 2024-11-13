package gtfs

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
)

var (
	ErrEmptyRoutesFile      = fmt.Errorf("empty agency file")
	ErrInvalidRoutesHeaders = fmt.Errorf("invalid agency headers")
	ErrNoRoutesRecords      = fmt.Errorf("no agency records")
)

type Route struct {
	ID                string   `json:"routeId,omitempty" csv:"route_id,omitempty"`
	AgencyID          string   `json:"agencyId,omitempty" csv:"agency_id,omitempty"`
	ShortName         string   `json:"routeShortName" csv:"route_short_name"`
	LongName          string   `json:"routeLongName" csv:"route_long_name"`
	Desc              string   `json:"routeDesc" csv:"route_desc"`
	Type              string   `json:"routeType" csv:"route_type"`
	URL               string   `json:"routeUrl" csv:"route_url"`
	Color             string   `json:"routeColor" csv:"route_color"`
	TextColor         string   `json:"routeTextColor" csv:"route_text_color"`
	SortOrder         string   `json:"routeSortOrder" csv:"route_sort_order"`
	ContinuousPickup  string   `json:"continuousPickup" csv:"continuous_pickup"`
	ContinuousDropOff string   `json:"continuousDropOff" csv:"continuous_drop_off"`
	NetworkID         string   `json:"networkId" csv:"network_id"`
	Unused            []string `json:"-" csv:"-"`
}

func parseRoutes(file *zip.File) ([]Route, error) {
	rc, err := file.Open()
	if err != nil {
		return []Route{}, err
	}
	defer rc.Close()

	lines, err := csv.NewReader(rc).ReadAll()
	if len(lines) == 0 {
		return []Route{}, ErrEmptyRoutesFile
	}

	headers := lines[0]
	if err := validateRoutesHeader(headers); err != nil {
		return []Route{}, err
	}

	records := lines[1:]
	if len(records) == 0 {
		return []Route{}, ErrNoRoutesRecords
	}

	routes := make([]Route, len(records))
	for i, field := range headers {
		for j, record := range records {
			switch field {
			case "route_id":
				routes[j].ID = record[i]
			case "agency_id":
				routes[j].AgencyID = record[i]
			case "route_short_name":
				routes[j].ShortName = record[i]
			case "route_long_name":
				routes[j].LongName = record[i]
			case "route_desc":
				routes[j].Desc = record[i]
			case "route_type":
				routes[j].Type = record[i]
			case "route_url":
				routes[j].URL = record[i]
			case "route_color":
				routes[j].Color = record[i]
			case "route_text_color":
				routes[j].TextColor = record[i]
			case "route_sort_order":
				routes[j].SortOrder = record[i]
			case "continuous_pickup":
				routes[j].ContinuousPickup = record[i]
			case "continuous_drop_off":
				routes[j].ContinuousDropOff = record[i]
			case "network_id":
				routes[j].NetworkID = record[i]
			default:
				if routes[j].Unused == nil {
					routes[j].Unused = []string{record[i]}
				} else {
					routes[j].Unused = append(routes[j].Unused, record[i])
				}
			}
		}
	}

	return routes, nil
}

func validateRoutesHeader(headers []string) error {
	requiredFields := []struct {
		name  string
		found bool
	}{
		{
			name:  "route_id",
			found: false},
		{
			name:  "route_type",
			found: false,
		},
	}

	for _, field := range headers {
		for i, req := range requiredFields {
			if field == req.name {
				requiredFields[i].found = true
			}
		}
	}

	for _, req := range requiredFields {
		if !req.found {
			return ErrInvalidRoutesHeaders
		}
	}

	return nil
}
