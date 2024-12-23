package db

import (
	"errors"
	"strconv"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

// Document types
const (
	DocTypeEvent = "event"
	DocTypeUnit  = "unit"
	DocTypeUser  = "user"
)

// Document statuses
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

var NameMap = map[string]uuid.UUID{
	DocTypeEvent:   uuid.MustParse("88c2333e-2bc2-4063-b865-719c24211d2c"),
	DocTypeUnit:    uuid.MustParse("915ddb34-93ba-4e2f-99f2-ea814bb2790d"),
	DocTypeUser:    uuid.MustParse("8d930d82-24e9-4fc1-824b-9e1253d4ee02"),
	StatusActive:   uuid.MustParse("5c96e882-112a-42e3-adf3-941f28ff9956"),
	StatusInactive: uuid.MustParse("aca32352-08c1-40ee-8a6e-e95d77e68724"),
}

var IDMap map[uuid.UUID]string

func init() {
	IDMap = make(map[uuid.UUID]string, len(NameMap))
	for k, v := range NameMap {
		IDMap[v] = k
	}
}

// Document Fields
const (
	Agency        = "a"
	ArrivalStop   = "as"
	ArrivalTime   = "at"
	CreatedAt     = "ca"
	DepartureStop = "ds"
	DepartureTime = "dt"
	ID            = "id"
	Notes         = "n"
	Route         = "r"
	Status        = "s"
	Type          = "t"
	Trip          = "tr"
	UnitID        = "u"
	UpdatedAt     = "ua"
	UnitCount     = "uc"
	UnitPosition  = "up"
	UserID        = "uid"
)

// Deserializers
func GetUUID(data dynamodb.AttributeValue) uuid.UUID {
	if data == nil {
		return uuid.Nil
	}

	if id, ok := data.(*dynamodb.AttributeValueMemberB); ok {
		value, err := uuid.FromBytes(id.Value)
		if err != nil {
			return uuid.Nil
		}
		return value
	}
	return uuid.Nil
}

func GetString(data dynamodb.AttributeValue) string {
	if data == nil {
		return ""
	}

	if s, ok := data.(*dynamodb.AttributeValueMemberS); ok {
		return s.Value
	}
	return ""
}

func GetStringSlice(data dynamodb.AttributeValue) []string {
	if data == nil {
		return nil
	}

	if ss, ok := data.(*dynamodb.AttributeValueMemberSS); ok {
		return ss.Value
	}
	return nil
}

func GetTime(data dynamodb.AttributeValue) time.Time {
	if data == nil {
		return time.Time{}
	}

	if n, ok := data.(*dynamodb.AttributeValueMemberN); ok {
		if i, err := strconv.ParseInt(n.Value, 10, 64); err == nil {
			return time.Unix(i, 0)
		}
	}
	return time.Time{}
}

func GetIntPtr(data dynamodb.AttributeValue) *int {
	if data == nil {
		return nil
	}

	if n, ok := data.(*dynamodb.AttributeValueMemberN); ok {
		if i, err := strconv.Atoi(n.Value); err == nil {
			return &i
		}
	}
	return nil
}

// Errors
var (
	ErrBadDocID     = errors.New("bad document ID")
	ErrBadDocType   = errors.New("bad document type")
	ErrBadDocStatus = errors.New("bad document status")
	ErrBadCreatedAt = errors.New("bad created at time")
	ErrBadUpdatedAt = errors.New("bad updated at time")
	ErrBadUserID    = errors.New("bad user ID")
	ErrBadCarrier   = errors.New("bad carrier")
	ErrBadUnitID    = errors.New("bad unit ID")
)

// DynamoDB Constants
const (
	DynamoDBBatchWriteLimit = 25
)

type DBDocument map[string]dynamodb.AttributeValue
