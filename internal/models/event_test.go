package models

import (
	"testing"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/db"
	"github.com/bridgelightcloud/bogie/internal/fixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetExampleEvent(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	evt := GetExampleEvent(uuid.Nil, uuid.Nil)

	assert.True(evt.Id != uuid.Nil)
	assert.True(evt.User != uuid.Nil)
	assert.Equal(evt.Type, db.DocTypeEvent)
}

func TestGetExampleEventArray(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	evts := GetExampleEventArray(5)

	assert.Len(evts, 5)
}

func TestMarshalDynamoDB(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	uuid1 := fixtures.GetTestUUID()
	uuid2 := fixtures.GetTestUUID()
	typeUUID := db.NameMap[db.DocTypeEvent]

	evt := GetExampleEvent(uuid1, uuid2)

	dbEvt, err := evt.MarshalDynamoDB()

	assert.Nil(err)
	assert.Equal(dbEvt["id"], &dynamodb.AttributeValueMemberB{Value: uuid1[:]})
	assert.Equal(dbEvt["t"], &dynamodb.AttributeValueMemberB{Value: typeUUID[:]})
}

func TestMarshDynamoDBErr(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name          string
		event         Event
		expectedError error
	}{
		{
			name: "ErrBadDocID",
			event: Event{
				Id: uuid.Nil,
			},
			expectedError: db.ErrBadDocID,
		},
		{
			name: "ErrBadDocType",
			event: Event{
				Id:   fixtures.GetTestUUID(),
				Type: "unknown-type",
			},
			expectedError: db.ErrBadDocType,
		},
		{
			name: "ErrBadDocStatus",
			event: Event{
				Id:     fixtures.GetTestUUID(),
				Type:   db.DocTypeEvent,
				Status: "unknown-status",
			},
			expectedError: db.ErrBadDocStatus,
		},
		{
			name: "ErrBadCreatedAt",
			event: Event{
				Id:     fixtures.GetTestUUID(),
				Type:   db.DocTypeEvent,
				Status: db.StatusActive,
			},
			expectedError: db.ErrBadCreatedAt,
		},
		{
			name: "ErrBadUpdatedAt",
			event: Event{
				Id:        fixtures.GetTestUUID(),
				Type:      db.DocTypeEvent,
				Status:    db.StatusActive,
				CreatedAt: fixtures.GetTestTimePtr(),
			},
			expectedError: db.ErrBadUpdatedAt,
		},
		{
			name: "ErrBadUser",
			event: Event{
				Id:        fixtures.GetTestUUID(),
				Type:      db.DocTypeEvent,
				Status:    db.StatusActive,
				CreatedAt: fixtures.GetTestTimePtr(),
				UpdatedAt: fixtures.GetTestTimePtr(),
			},
			expectedError: db.ErrBadUserID,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			dbEvt, err := tc.event.MarshalDynamoDB()

			assert.Nil(dbEvt)
			assert.Equal(tc.expectedError, err)
		})
	}
}
