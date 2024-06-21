package event

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
	assert.Equal(evt.Type, db.Event)
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
	typeUUID := db.NameMap[db.Event]

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
			name: "ErrBadEventID",
			event: Event{
				Id: uuid.Nil,
			},
			expectedError: ErrBadEventID,
		},
		{
			name: "ErrBadEventType",
			event: Event{
				Id:   fixtures.GetTestUUID(),
				Type: "unknown-type",
			},
			expectedError: ErrBadDocumentType,
		},
		{
			name: "ErrBadCreatedAt",
			event: Event{
				Id:   fixtures.GetTestUUID(),
				Type: db.Event,
			},
			expectedError: ErrBadCreatedAt,
		},
		{
			name: "ErrBadUpdatedAt",
			event: Event{
				Id:        fixtures.GetTestUUID(),
				Type:      db.Event,
				CreatedAt: fixtures.GetTestTimePtr(),
			},
			expectedError: ErrBadUpdatedAt,
		},
		{
			name: "ErrBadUser",
			event: Event{
				Id:        fixtures.GetTestUUID(),
				Type:      db.Event,
				CreatedAt: fixtures.GetTestTimePtr(),
				UpdatedAt: fixtures.GetTestTimePtr(),
			},
			expectedError: ErrBadUser,
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
