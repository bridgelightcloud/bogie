package event

import "github.com/google/uuid"

type Event struct {
	Id            uuid.UUID `json:"id"`
	Type          string    `json:"type,omitempty"`
	Carrier       string    `json:"carrier,omitempty"`
	Line          string    `json:"line,omitempty"`
	DepartureStop string    `json:"departureStop,omitempty"`
	ArrivalStop   string    `json:"arrivalStop,omitempty"`
	DepartureTime string    `json:"departureTime,omitempty"`
	ArrivalTime   string    `json:"arrivalTime,omitempty"`
}
