package models

import (
	"strconv"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/db"
	"github.com/google/uuid"
)

type Event struct {
	Id            uuid.UUID  `json:"id"`                      // id
	Type          string     `json:"type"`                    // t
	Status        string     `json:"status,omitempty"`        // s
	CreatedAt     *time.Time `json:"createdAt"`               // ca
	UpdatedAt     *time.Time `json:"updatedAt"`               // ua
	User          uuid.UUID  `json:"user,omitempty"`          // uid
	Agency        string     `json:"agency,omitempty"`        // a
	Route         string     `json:"route,omitempty"`         // r
	Trip          string     `json:"trip,omitempty"`          // tr
	UnitID        string     `json:"unitID,omitempty"`        // u
	UnitCount     *int       `json:"unitCount,omitempty"`     // uc
	UnitPosition  *int       `json:"unitPosition,omitempty"`  // up
	DepartureStop string     `json:"departureStop,omitempty"` // ds
	ArrivalStop   string     `json:"arrivalStop,omitempty"`   // as
	DepartureTime *time.Time `json:"departureTime,omitempty"` // dt
	ArrivalTime   *time.Time `json:"arrivalTime,omitempty"`   // at
	Notes         []string   `json:"notes,omitempty"`         // n
}

func (e Event) MarshalDynamoDB() (db.DBDocument, error) {
	if e.Id == uuid.Nil {
		return nil, db.ErrBadDocID
	}

	data := db.DBDocument{
		db.ID: &dynamodb.AttributeValueMemberB{Value: e.Id[:]},
	}

	if id, ok := db.NameMap[e.Type]; ok {
		data[db.Type] = &dynamodb.AttributeValueMemberB{Value: id[:]}
	} else {
		return nil, db.ErrBadDocType
	}

	if id, ok := db.NameMap[e.Status]; ok {
		data[db.Status] = &dynamodb.AttributeValueMemberB{Value: id[:]}
	} else {
		return nil, db.ErrBadDocStatus
	}

	if e.CreatedAt != nil && !e.CreatedAt.IsZero() {
		data[db.CreatedAt] = &dynamodb.AttributeValueMemberN{
			Value: strconv.FormatInt(e.CreatedAt.Unix(), 10),
		}
	} else {
		return nil, db.ErrBadCreatedAt
	}

	if e.UpdatedAt != nil && !e.UpdatedAt.IsZero() {
		data[db.UpdatedAt] = &dynamodb.AttributeValueMemberN{
			Value: strconv.FormatInt(e.UpdatedAt.Unix(), 10),
		}
	} else {
		return nil, db.ErrBadUpdatedAt
	}

	if e.User != uuid.Nil {
		data[db.UserID] = &dynamodb.AttributeValueMemberB{Value: e.User[:]}
	} else {
		return nil, db.ErrBadUserID
	}

	if e.Agency != "" {
		data[db.Agency] = &dynamodb.AttributeValueMemberS{Value: e.Agency}
	}

	if e.Route != "" {
		data[db.Route] = &dynamodb.AttributeValueMemberS{Value: e.Route}
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

func (e *Event) UnmarshalDynamoDB(data db.DBDocument) error {
	if id := db.GetUUID(data["id"]); id != uuid.Nil {
		e.Id = id
	} else {
		return db.ErrBadDocID
	}

	if t, ok := db.IDMap[db.GetUUID(data[db.Type])]; ok {
		e.Type = t
	} else {
		return db.ErrBadDocType
	}

	if s, ok := db.IDMap[db.GetUUID(data[db.Status])]; ok {
		e.Status = s
	} else {
		return db.ErrBadDocStatus
	}

	if t := db.GetTime(data[db.CreatedAt]); !t.IsZero() {
		e.CreatedAt = &t
	} else {
		return db.ErrBadCreatedAt
	}

	if t := db.GetTime(data[db.UpdatedAt]); !t.IsZero() {
		e.UpdatedAt = &t
	} else {
		return db.ErrBadUpdatedAt
	}

	if u := db.GetUUID(data[db.UserID]); u != uuid.Nil {
		e.User = u
	}

	e.Agency = db.GetString(data[db.Agency])
	e.Route = db.GetString(data[db.Route])
	e.Trip = db.GetString(data[db.Trip])
	e.UnitID = db.GetString(data[db.UnitID])
	e.UnitCount = db.GetIntPtr(data[db.UnitCount])
	e.UnitPosition = db.GetIntPtr(data[db.UnitPosition])
	e.DepartureStop = db.GetString(data[db.DepartureStop])
	e.ArrivalStop = db.GetString(data[db.ArrivalStop])

	if t := db.GetTime(data[db.DepartureTime]); !t.IsZero() {
		e.DepartureTime = &t
	}

	if t := db.GetTime(data[db.ArrivalTime]); !t.IsZero() {
		e.ArrivalTime = &t
	}

	e.Notes = db.GetStringSlice(data[db.Notes])

	return nil
}
