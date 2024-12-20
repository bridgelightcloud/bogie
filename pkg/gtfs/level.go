package gtfs

import (
	"fmt"
	"math"
)

type Level struct {
	ID    string  `json:"levelId" csv:"level_id"`
	Index float64 `json:"levelIndex" csv:"level_index"`
	Name  string  `json:"levelName,omitempty" csv:"level_name"`
}

func (l Level) key() string {
	return l.ID
}

func (l Level) validate() errorList {
	var errs errorList

	if l.ID == "" {
		errs.add(fmt.Errorf("missing level_id"))
	}
	if l.Index == math.Inf(-1) {
		errs.add(fmt.Errorf("invalid index valie"))
	}

	return errs
}
