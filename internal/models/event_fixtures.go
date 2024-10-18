package models

import (
	"strconv"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/db"
	"github.com/bridgelightcloud/bogie/internal/fixtures"
	"github.com/google/uuid"
)

var (
	defaultAgency        = "BART"
	defaultRoute         = "Red"
	defaultTrip          = "123"
	defaultUnitID        = "3095"
	defaultUnitCount     = 6
	defaultUnitPosition  = 1
	defaultDepartureStop = "Richmond"
	defaultArrivalStop   = "Millbrae"
	defaultNotes         = []string{
		"Very Full",
		"Stopped at 12th St. for 5 minutes",
	}
)

func GetExampleEvent(id uuid.UUID, user uuid.UUID) Event {
	fixtures.MaybeRefreshUUID(&id)
	fixtures.MaybeRefreshUUID(&user)

	exTime := time.Now().Truncate(time.Second)

	return Event{
		Id:            id,
		Type:          db.DocTypeEvent,
		Status:        db.StatusActive,
		CreatedAt:     &exTime,
		UpdatedAt:     &exTime,
		User:          user,
		Agency:        defaultAgency,
		Route:         defaultRoute,
		Trip:          defaultTrip,
		UnitID:        defaultUnitID,
		UnitCount:     &defaultUnitCount,
		UnitPosition:  &defaultUnitPosition,
		DepartureStop: defaultDepartureStop,
		ArrivalStop:   defaultArrivalStop,
		DepartureTime: &exTime,
		ArrivalTime:   &exTime,
		Notes:         defaultNotes,
	}
}

func GetExampleEventArray(count int) []Event {
	evs := make([]Event, count)
	for i := 0; i < count; i++ {
		evs[i] = GetExampleEvent(uuid.Nil, uuid.Nil)
	}
	return evs
}

func GetExampleEventDBDocument(id uuid.UUID, user uuid.UUID) db.DBDocument {
	fixtures.MaybeRefreshUUID(&id)
	fixtures.MaybeRefreshUUID(&user)

	exDocType := db.NameMap[db.DocTypeEvent]
	exStatus := db.NameMap[db.StatusActive]
	exTime := strconv.FormatInt(time.Now().Truncate(time.Second).Unix(), 10)

	return db.DBDocument{
		db.ID:            &dynamodb.AttributeValueMemberB{Value: id[:]},
		db.Type:          &dynamodb.AttributeValueMemberB{Value: exDocType[:]},
		db.Status:        &dynamodb.AttributeValueMemberB{Value: exStatus[:]},
		db.CreatedAt:     &dynamodb.AttributeValueMemberN{Value: exTime},
		db.UpdatedAt:     &dynamodb.AttributeValueMemberN{Value: exTime},
		db.UserID:        &dynamodb.AttributeValueMemberB{Value: user[:]},
		db.Agency:        &dynamodb.AttributeValueMemberS{Value: defaultAgency},
		db.Route:         &dynamodb.AttributeValueMemberS{Value: defaultRoute},
		db.Trip:          &dynamodb.AttributeValueMemberS{Value: defaultTrip},
		db.UnitID:        &dynamodb.AttributeValueMemberS{Value: defaultUnitID},
		db.UnitCount:     &dynamodb.AttributeValueMemberN{Value: strconv.Itoa(defaultUnitCount)},
		db.UnitPosition:  &dynamodb.AttributeValueMemberN{Value: strconv.Itoa(defaultUnitPosition)},
		db.DepartureStop: &dynamodb.AttributeValueMemberS{Value: defaultDepartureStop},
		db.ArrivalStop:   &dynamodb.AttributeValueMemberS{Value: defaultArrivalStop},
		db.DepartureTime: &dynamodb.AttributeValueMemberN{Value: exTime},
		db.ArrivalTime:   &dynamodb.AttributeValueMemberN{Value: exTime},
		db.Notes:         &dynamodb.AttributeValueMemberSS{Value: defaultNotes},
	}
}

func GetExampleEventDBDocumentArray(count int) []db.DBDocument {
	evDocs := make([]db.DBDocument, count)
	for i := 0; i < count; i++ {
		evDocs[i] = GetExampleEventDBDocument(uuid.Nil, uuid.Nil)
	}
	return evDocs
}
