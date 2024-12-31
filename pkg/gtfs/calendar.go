package gtfs

type Calendar struct {
	ServiceID string `json:"serviceId" csv:"service_id"`
	Monday    int    `json:"monday" csv:"monday"`
	Tuesday   int    `json:"tuesday" csv:"tuesday"`
	Wednesday int    `json:"wednesday" csv:"wednesday"`
	Thursday  int    `json:"thursday" csv:"thursday"`
	Friday    int    `json:"friday" csv:"friday"`
	Saturday  int    `json:"saturday" csv:"saturday"`
	Sunday    int    `json:"sunday" csv:"sunday"`
	StartDate Date   `json:"startDate" csv:"start_date"`
	EndDate   Date   `json:"endDate" csv:"end_date"`
}

func (c Calendar) key() string {
	return c.ServiceID
}

func (c Calendar) validate() errorList {
	var errs errorList
	return errs
}
