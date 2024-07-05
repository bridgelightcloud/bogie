package ddbmarshal

import (
	"fmt"
	"strconv"
	"time"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UnixTimestamp time.Time

func (t *UnixTimestamp) UnmarshalAttributeValue(v dynamodb.AttributeValue) error {
	if attr, ok := v.(*dynamodb.AttributeValueMemberN); ok {
		i, err := strconv.ParseInt(attr.Value, 10, 64)
		if err != nil {
			return err
		}

		*t = UnixTimestamp(time.Unix(i, 0))
		return nil
	} else {
		return fmt.Errorf("expected N, got %T", v)
	}
}

func (t UnixTimestamp) MarshalAttributeValue() (dynamodb.AttributeValue, error) {
	return &dynamodb.AttributeValueMemberN{Value: strconv.FormatInt(time.Time(t).Unix(), 10)}, nil
}

func (t UnixTimestamp) String() string {
	return time.Time(t).Format(time.RFC3339)
}
