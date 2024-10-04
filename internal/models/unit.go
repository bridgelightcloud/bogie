package models

import (
	"strconv"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/db"
	"github.com/google/uuid"
)

type Unit struct {
	Id        uuid.UUID  `json:"id"`
	Type      string     `json:"type"`
	Status    string     `json:"status,omitempty"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	Agency    string     `json:"agency,omitempty"`
	UnitID    string     `json:"unitID,omitempty"`
	Notes     []string   `json:"notes,omitempty"`
}

func GetExampleUnit() Unit {
	t := time.Now().Truncate(time.Second)

	return Unit{
		Id:        uuid.New(),
		Type:      db.DocTypeUnit,
		Status:    db.StatusActive,
		CreatedAt: &t,
		UpdatedAt: &t,
		Agency:    "BART",
		UnitID:    "1234",
	}
}

func (u Unit) MarshalDynamoDB() (DBDocument, error) {
	if u.Id == uuid.Nil {
		return nil, db.ErrBadDocID
	}

	data := DBDocument{
		db.ID: &dynamodb.AttributeValueMemberB{Value: u.Id[:]},
	}

	if id, ok := db.NameMap[u.Type]; ok {
		data[db.Type] = &dynamodb.AttributeValueMemberB{Value: id[:]}
	} else {
		return nil, db.ErrBadDocType
	}

	if id, ok := db.NameMap[u.Status]; ok {
		data[db.Status] = &dynamodb.AttributeValueMemberB{Value: id[:]}
	} else {
		return nil, db.ErrBadDocStatus
	}

	if u.CreatedAt != nil && !u.CreatedAt.IsZero() {
		data[db.CreatedAt] = &dynamodb.AttributeValueMemberN{
			Value: strconv.FormatInt(u.CreatedAt.Unix(), 10),
		}
	} else {
		return nil, db.ErrBadCreatedAt
	}

	if u.UpdatedAt != nil && !u.UpdatedAt.IsZero() {
		data[db.UpdatedAt] = &dynamodb.AttributeValueMemberN{
			Value: strconv.FormatInt(u.UpdatedAt.Unix(), 10),
		}
	} else {
		return nil, db.ErrBadUpdatedAt
	}

	if u.Agency != "" {
		data[db.Agency] = &dynamodb.AttributeValueMemberS{Value: u.Agency}
	}

	if u.UnitID != "" {
		data[db.UnitID] = &dynamodb.AttributeValueMemberS{Value: u.Agency}
	}

	if len(u.Notes) > 0 {
		data[db.Notes] = &dynamodb.AttributeValueMemberSS{Value: u.Notes}
	}

	return data, nil
}

func (u *Unit) UnmarshalDynamoDB(data DBDocument) error {
	if id := db.GetUUID(data[db.ID]); id != uuid.Nil {
		u.Id = id
	} else {
		return db.ErrBadDocID
	}

	if t, ok := db.IDMap[db.GetUUID(data[db.Type])]; ok {
		u.Type = t
	} else {
		return db.ErrBadDocType
	}

	if s, ok := db.IDMap[db.GetUUID(data[db.Status])]; ok {
		u.Status = s
	} else {
		return db.ErrBadDocStatus
	}

	if t := db.GetTime(data[db.CreatedAt]); !t.IsZero() {
		u.CreatedAt = &t
	} else {
		return db.ErrBadCreatedAt
	}

	if t := db.GetTime(data[db.UpdatedAt]); !t.IsZero() {
		u.UpdatedAt = &t
	} else {
		return db.ErrBadUpdatedAt
	}

	u.Agency = db.GetString(data[db.Agency])
	u.UnitID = db.GetString(data[db.UnitID])
	u.Notes = db.GetStringSlice(data[db.Notes])

	return nil
}
