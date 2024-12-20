package gtfs

import (
	"fmt"
)

type Route struct {
	ID                string `json:"routeId" csv:"route_id"`
	AgencyID          string `json:"agencyId" csv:"agency_id"`
	ShortName         string `json:"routeShortName" csv:"route_short_name"`
	LongName          string `json:"routeLongName" csv:"route_long_name"`
	Desc              string `json:"routeDesc,omitempty" csv:"route_desc"`
	Type              string `json:"routeType" csv:"route_type"`
	URL               string `json:"routeUrl,omitempty" csv:"route_url"`
	Color             string `json:"routeColor,omitempty" csv:"route_color"`
	TextColor         string `json:"routeTextColor,omitempty" csv:"route_text_color"`
	SortOrder         string `json:"routeSortOrder,omitempty" csv:"route_sort_order"`
	ContinuousPickup  string `json:"continuousPickup,omitempty" csv:"continuous_pickup"`
	ContinuousDropOff string `json:"continuousDropOff,omitempty" csv:"continuous_drop_off"`
	NetworkID         string `json:"networkId,omitempty" csv:"network_id"`
}

func (r Route) key() string {
	return r.ID
}

func (r Route) validate() errorList {
	var errs errorList

	if r.ID == "" {
		errs.add(fmt.Errorf("route ID is required"))
	}
	if r.ShortName == "" {
		errs.add(fmt.Errorf("route short name is required"))
	}
	if r.LongName == "" {
		errs.add(fmt.Errorf("route long name is required"))
	}
	if r.Type == "" {
		errs.add(fmt.Errorf("route type is required"))
	}

	return errs
}
