package event

import (
	"errors"
	"strconv"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/db"
	"github.com/google/uuid"
)

type Event struct {
	Id            uuid.UUID  `json:"id"`
	Type          string     `json:"type"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	User          uuid.UUID  `json:"user,omitempty"`
	Carrier       string     `json:"carrier,omitempty"`
	Line          string     `json:"line,omitempty"`
	Trip          string     `json:"trip,omitempty"`
	UnitID        string     `json:"unitID,omitempty"`
	UnitCount     *int       `json:"unitCount,omitempty"`
	UnitPosition  *int       `json:"unitPosition,omitempty"`
	DepartureStop string     `json:"departureStop,omitempty"`
	ArrivalStop   string     `json:"arrivalStop,omitempty"`
	DepartureTime *time.Time `json:"departureTime,omitempty"`
	ArrivalTime   *time.Time `json:"arrivalTime,omitempty"`
	Notes         []string   `json:"notes,omitempty"`
}

func GetExampleEvent(id uuid.UUID, user uuid.UUID) Event {
	if id == uuid.Nil {
		id = uuid.New()
	}

	if user == uuid.Nil {
		user = uuid.New()
	}

	extime := time.Now().Truncate(time.Second)
	excount := 6
	exposition := 1

	return Event{
		Id:            id,
		Type:          db.EventDoc,
		CreatedAt:     &extime,
		UpdatedAt:     &extime,
		User:          user,
		Carrier:       "BART",
		Line:          "Red",
		Trip:          "123",
		UnitID:        "3095",
		UnitCount:     &excount,
		UnitPosition:  &exposition,
		DepartureStop: "Richmond",
		ArrivalStop:   "Millbrae",
		DepartureTime: &extime,
		ArrivalTime:   &extime,
		Notes: []string{
			"Very Full",
			"Stopped at 12th St. for 5 minutes",
		},
	}
}

func GetExampleEventArray(count int) []Event {
	evs := make([]Event, count)
	for i := 0; i < count; i++ {
		evs[i] = GetExampleEvent(uuid.Nil, uuid.Nil)
	}
	return evs
}

var ErrBadEventID = errors.New("bad event ID")
var ErrBadDocumentType = errors.New("bad event type")
var ErrBadCreatedAt = errors.New("bad created at time")
var ErrBadUpdatedAt = errors.New("bad updated at time")
var ErrBadUser = errors.New("bad user ID")

func (e Event) MarshalDynamoDB() (map[string]dynamodb.AttributeValue, error) {
	if e.Id == uuid.Nil {
		return nil, ErrBadEventID
	}

	data := map[string]dynamodb.AttributeValue{
		db.ID: &dynamodb.AttributeValueMemberB{Value: e.Id[:]},
	}

	if id, ok := db.NameMap[e.Type]; ok {
		data[db.Type] = &dynamodb.AttributeValueMemberB{Value: id[:]}
	} else {
		return nil, ErrBadDocumentType
	}

	if e.CreatedAt != nil && !e.CreatedAt.IsZero() {
		data[db.CreatedAt] = &dynamodb.AttributeValueMemberN{Value: strconv.FormatInt(e.CreatedAt.Unix(), 10)}
	} else {
		return nil, ErrBadCreatedAt
	}

	if e.UpdatedAt != nil && !e.UpdatedAt.IsZero() {
		data[db.UpdatedAt] = &dynamodb.AttributeValueMemberN{Value: strconv.FormatInt(e.UpdatedAt.Unix(), 10)}
	} else {
		return nil, ErrBadUpdatedAt
	}

	if e.User != uuid.Nil {
		data[db.UserID] = &dynamodb.AttributeValueMemberB{Value: e.User[:]}
	} else {
		return nil, ErrBadUser
	}

	if e.Carrier != "" {
		data[db.Carrier] = &dynamodb.AttributeValueMemberS{Value: e.Carrier}
	}

	if e.Line != "" {
		data[db.Line] = &dynamodb.AttributeValueMemberS{Value: e.Line}
	}

	if e.Trip != "" {
		data[db.Trip] = &dynamodb.AttributeValueMemberS{Value: e.Trip}
	}

	if e.UnitID != "" {
		data[db.UnitID] = &dynamodb.AttributeValueMemberS{Value: e.UnitID}
	}

	if e.UnitCount != nil {
		data[db.UnitCount] = &dynamodb.AttributeValueMemberN{Value: strconv.Itoa(*e.UnitCount)}
	}

	if e.UnitPosition != nil {
		data[db.UnitPosition] = &dynamodb.AttributeValueMemberN{Value: strconv.Itoa(*e.UnitPosition)}
	}

	if e.DepartureStop != "" {
		data[db.DepartureStop] = &dynamodb.AttributeValueMemberS{Value: e.DepartureStop}
	}

	if e.ArrivalStop != "" {
		data[db.ArrivalStop] = &dynamodb.AttributeValueMemberS{Value: e.ArrivalStop}
	}

	if e.DepartureTime != nil && !e.DepartureTime.IsZero() {
		data[db.DepartureTime] = &dynamodb.AttributeValueMemberN{Value: strconv.FormatInt(e.DepartureTime.Unix(), 10)}
	}

	if e.ArrivalTime != nil && !e.ArrivalTime.IsZero() {
		data[db.ArrivalTime] = &dynamodb.AttributeValueMemberN{Value: strconv.FormatInt(e.ArrivalTime.Unix(), 10)}
	}

	if len(e.Notes) > 0 {
		data[db.Notes] = &dynamodb.AttributeValueMemberSS{Value: e.Notes}
	}

	return data, nil
}

func (e *Event) UnmarshalDynamoDB(data map[string]dynamodb.AttributeValue) error {
	if id := db.GetUUID(data["id"]); id != uuid.Nil {
		e.Id = id
	} else {
		return ErrBadEventID
	}

	e.Type = db.IDMap[db.GetUUID(data[db.Type])]

	if t := db.GetTime(data[db.CreatedAt]); !t.IsZero() {
		e.CreatedAt = &t
	}

	if t := db.GetTime(data[db.UpdatedAt]); !t.IsZero() {
		e.UpdatedAt = &t
	}

	e.Carrier = db.GetString(data[db.Carrier])
	e.Line = db.GetString(data[db.Line])
	e.Trip = db.GetString(data[db.Trip])
	e.UnitID = db.GetString(data[db.UnitID])
	e.UnitCount = db.GetIntPtr(data[db.UnitCount])
	e.UnitPosition = db.GetIntPtr(data[db.UnitPosition])
	e.DepartureStop = db.GetString(data[db.DepartureStop])
	e.ArrivalStop = db.GetString(data[db.ArrivalStop])
	e.Notes = db.GetStringSlice(data[db.Notes])

	if t := db.GetTime(data[db.DepartureTime]); !t.IsZero() {
		e.DepartureTime = &t
	}

	if t := db.GetTime(data[db.ArrivalTime]); !t.IsZero() {
		e.ArrivalTime = &t
	}

	return nil
}
