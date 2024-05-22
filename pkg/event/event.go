package event

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id            uuid.UUID `json:"id"`
	Type          string    `json:"type,omitempty"`
	Carrier       string    `json:"carrier,omitempty"`
	Line          string    `json:"line,omitempty"`
	DepartureStop string    `json:"departureStop,omitempty"`
	ArrivalStop   string    `json:"arrivalStop,omitempty"`
	DepartureTime time.Time `json:"departureTime,omitempty"`
	ArrivalTime   time.Time `json:"arrivalTime,omitempty"`
}

func GetExampleEvent(id uuid.UUID) Event {
	if id == uuid.Nil {
		id = uuid.New()
	}

	return Event{
		Id:            id,
		Type:          "train",
		Carrier:       "BART",
		Line:          "Red",
		DepartureStop: "Richmond",
		ArrivalStop:   "Millbrae",
		DepartureTime: time.Now().Truncate(time.Second),
		ArrivalTime:   time.Now().Truncate(time.Second),
	}
}

func GetExampleEventArray(count int) []Event {
	evs := make([]Event, count)
	for i := 0; i < count; i++ {
		evs[i] = GetExampleEvent(uuid.Nil)
	}
	return evs
}
