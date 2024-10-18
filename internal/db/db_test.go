package db

import (
	"testing"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/fixtures"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUUID(t *testing.T) {
	t.Parallel()

	tt := []struct {
		dbValue      dynamodb.AttributeValue
		expectedUuid uuid.UUID
	}{
		{
			// nil uuid
			dbValue:      &dynamodb.AttributeValueMemberB{Value: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			expectedUuid: uuid.Nil,
		},
		{
			// valid uuid
			dbValue:      &dynamodb.AttributeValueMemberB{Value: []byte{0x8d, 0x93, 0x0d, 0x82, 0x24, 0xe9, 0x4f, 0xc1, 0x82, 0x4b, 0x9e, 0x12, 0x53, 0xd4, 0xee, 0x02}},
			expectedUuid: uuid.MustParse("8d930d82-24e9-4fc1-824b-9e1253d4ee02"),
		},
		{
			// empty uuid
			dbValue:      &dynamodb.AttributeValueMemberB{Value: []byte{}},
			expectedUuid: uuid.Nil,
		},
		{
			// other AttributeValue type
			dbValue:      &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedUuid: uuid.Nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedUUID := GetUUID(tc.dbValue)

			assert.Equal(tc.expectedUuid, parsedUUID)
		})
	}
}

func TestGetString(t *testing.T) {
	t.Parallel()

	tt := []struct {
		dbValue     dynamodb.AttributeValue
		expectedStr string
	}{
		{
			// nil/empty string
			dbValue:     &dynamodb.AttributeValueMemberS{Value: ""},
			expectedStr: "",
		},
		{
			// valid string
			dbValue:     &dynamodb.AttributeValueMemberS{Value: "test"},
			expectedStr: "test",
		},
		{
			// other AttributeValue type
			dbValue:     &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedStr: "",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedStr := GetString(tc.dbValue)

			assert.Equal(tc.expectedStr, parsedStr)
		})
	}
}

func TestGetStringSlice(t *testing.T) {
	t.Parallel()

	tt := []struct {
		dbValue             dynamodb.AttributeValue
		expectedStringSlice []string
	}{
		{
			// nil string slice
			dbValue:             &dynamodb.AttributeValueMemberSS{Value: nil},
			expectedStringSlice: nil,
		},
		{
			// empty string slice
			dbValue:             &dynamodb.AttributeValueMemberSS{Value: []string{}},
			expectedStringSlice: []string{},
		},
		{
			// valid string slice
			dbValue:             &dynamodb.AttributeValueMemberSS{Value: []string{"test1", "test2"}},
			expectedStringSlice: []string{"test1", "test2"},
		},
		{
			// other AttributeValue type
			dbValue:             &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedStringSlice: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedStringSlice := GetStringSlice(tc.dbValue)

			assert.Equal(tc.expectedStringSlice, parsedStringSlice)
		})
	}
}

func TestGetTime(t *testing.T) {
	t.Parallel()

	tt := []struct {
		dbValue    dynamodb.AttributeValue
		expectedTm time.Time
	}{
		{
			// nil time
			dbValue:    &dynamodb.AttributeValueMemberN{Value: ""},
			expectedTm: time.Time{},
		},
		{
			// valid time
			dbValue:    &dynamodb.AttributeValueMemberN{Value: "1616740000"},
			expectedTm: time.Unix(1616740000, 0),
		},
		{
			// other AttributeValue type
			dbValue:    &dynamodb.AttributeValueMemberBOOL{Value: true},
			expectedTm: time.Time{},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedTime := GetTime(tc.dbValue)

			assert.Equal(tc.expectedTm, parsedTime)
		})
	}
}

func TestGetIntPtr(t *testing.T) {
	t.Parallel()

	tt := []struct {
		dbValue  dynamodb.AttributeValue
		expected *int
	}{
		{
			// nil int pointer
			dbValue:  &dynamodb.AttributeValueMemberN{Value: ""},
			expected: nil,
		},
		{
			// valid int pointer
			dbValue:  &dynamodb.AttributeValueMemberN{Value: "42"},
			expected: fixtures.IntPtr(42),
		},
		{
			// other AttributeValue type
			dbValue:  &dynamodb.AttributeValueMemberBOOL{Value: true},
			expected: nil,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(t.Name(), func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)

			parsedIntPtr := GetIntPtr(tc.dbValue)

			assert.Equal(tc.expected, parsedIntPtr)
		})
	}
}
