package models

import (
	"time"

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
		Type:      db.UnitDoc,
		Status:    db.ActiveStatus,
		CreatedAt: &t,
		UpdatedAt: &t,
		Agency:    "BART",
		UnitID:    "1234",
	}
}

func (u Unit) MarshalDynamoDB() (DDBDocument, error) {
	return nil, nil
}

func (u *Unit) UnmarshalDynamoDB(item DDBDocument) error {
	return nil
}
