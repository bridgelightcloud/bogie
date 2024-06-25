package db

import (
	"strconv"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	ArrivalStop   = "as"
	ArrivalTime   = "at"
	Carrier       = "c"
	CreatedAt     = "ca"
	DepartureStop = "ds"
	DepartureTime = "dt"
	ID            = "id"
	Line          = "l"
	Notes         = "n"
	Type          = "t"
	Trip          = "tr"
	UnitID        = "u"
	UpdatedAt     = "ua"
	UnitCount     = "uc"
	UnitPosition  = "up"
	UserID        = "uid"
)

func GetUUID(data dynamodb.AttributeValue) uuid.UUID {
	if data == nil {
		return uuid.Nil
	}

	if id, ok := data.(*dynamodb.AttributeValueMemberB); ok {
		return uuid.UUID(id.Value)
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
