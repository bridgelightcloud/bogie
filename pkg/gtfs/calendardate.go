package gtfs

import "fmt"

type CalendarDate struct {
	ServiceID     string `json:"serviceId" csv:"service_id"`
	Date          Date   `json:"date" csv:"date"`
	ExceptionType int    `json:"exceptionType" csv:"exception_type"`
}

func (c CalendarDate) key() string {
	return c.ServiceID
}

func (c CalendarDate) validate() errorList {
	var errs errorList

	if c.ServiceID == "" {
		errs.add(fmt.Errorf("service ID is required"))
	}
	if c.Date.IsZero() {
		errs.add(fmt.Errorf("date is required"))
	}
	if c.ExceptionType != 1 && c.ExceptionType != 2 {
		errs.add(fmt.Errorf("invalid exception type: %d", c.ExceptionType))
	}

	return errs
}
