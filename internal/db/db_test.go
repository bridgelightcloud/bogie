package db

import (
	"testing"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/fixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type deserTestCase[T any] struct {
	dbValue       dynamodb.AttributeValue
	expectedValue T
}

func TestGetUUID(t *testing.T) {
	t.Parallel()

	tt := []deserTestCase[uuid.UUID]{
		{
			// nil
			dbValue:       nil,
			expectedValue: uuid.Nil,
		},
		{
			// nil uuid
			dbValue:       &dynamodb.AttributeValueMemberB{Value: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			expectedValue: uuid.Nil,
		},
		{
			// valid uuid
			dbValue:       &dynamodb.AttributeValueMemberB{Value: []byte{0x8d, 0x93, 0x0d, 0x82, 0x24, 0xe9, 0x4f, 0xc1, 0x82, 0x4b, 0x9e, 0x12, 0x53, 0xd4, 0xee, 0x02}},
			expectedValue: uuid.MustParse("8d930d82-24e9-4fc1-824b-9e1253d4ee02"),
		},
		{
			// empty uuid
			dbValue:       &dynamodb.AttributeValueMemberB{Value: []byte{}},
			expectedValue: uuid.Nil,
		},
		{
			// other AttributeValue type
			dbValue:       &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedValue: uuid.Nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedUUID := GetUUID(tc.dbValue)

			assert.Equal(tc.expectedValue, parsedUUID)
		})
	}
}

func TestGetString(t *testing.T) {
	t.Parallel()

	tt := []deserTestCase[string]{
		{
			// nil string
			dbValue:       nil,
			expectedValue: "",
		},
		{
			// empty string
			dbValue:       &dynamodb.AttributeValueMemberS{Value: ""},
			expectedValue: "",
		},
		{
			// valid string
			dbValue:       &dynamodb.AttributeValueMemberS{Value: "test"},
			expectedValue: "test",
		},
		{
			// other AttributeValue type
			dbValue:       &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedValue: "",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedStr := GetString(tc.dbValue)

			assert.Equal(tc.expectedValue, parsedStr)
		})
	}
}

func TestGetStringSlice(t *testing.T) {
	t.Parallel()

	tt := []deserTestCase[[]string]{
		{
			// nil slice
			dbValue:       nil,
			expectedValue: nil,
		},
		{
			// empty string slice
			dbValue:       &dynamodb.AttributeValueMemberSS{Value: []string{}},
			expectedValue: []string{},
		},
		{
			// valid string slice
			dbValue:       &dynamodb.AttributeValueMemberSS{Value: []string{"test1", "test2"}},
			expectedValue: []string{"test1", "test2"},
		},
		{
			// other AttributeValue type
			dbValue:       &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedValue: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedStringSlice := GetStringSlice(tc.dbValue)

			assert.Equal(tc.expectedValue, parsedStringSlice)
		})
	}
}

func TestGetTime(t *testing.T) {
	t.Parallel()

	tt := []deserTestCase[time.Time]{
		{
			// nil time
			dbValue:       nil,
			expectedValue: time.Time{},
		},
		{
			// empty time
			dbValue:       &dynamodb.AttributeValueMemberN{Value: ""},
			expectedValue: time.Time{},
		},
		{
			// valid time
			dbValue:       &dynamodb.AttributeValueMemberN{Value: "1616740000"},
			expectedValue: time.Unix(1616740000, 0),
		},
		{
			// other AttributeValue type
			dbValue:       &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedValue: time.Time{},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedTime := GetTime(tc.dbValue)

			assert.Equal(tc.expectedValue, parsedTime)
		})
	}
}

func TestGetIntPtr(t *testing.T) {
	t.Parallel()

	tt := []deserTestCase[*int]{
		{
			// nil
			dbValue:       nil,
			expectedValue: nil,
		},
		{
			// nil int pointer
			dbValue:       &dynamodb.AttributeValueMemberN{Value: ""},
			expectedValue: nil,
		},
		{
			// valid int pointer
			dbValue:       &dynamodb.AttributeValueMemberN{Value: "42"},
			expectedValue: fixtures.IntPtr(42),
		},
		{
			// other AttributeValue type
			dbValue:       &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedValue: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedIntPtr := GetIntPtr(tc.dbValue)

			assert.Equal(tc.expectedValue, parsedIntPtr)
		})
	}
}
