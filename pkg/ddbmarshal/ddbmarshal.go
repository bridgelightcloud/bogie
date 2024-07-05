package ddbmarshal

import (
	"fmt"
	"reflect"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBMarshaler interface {
	MarshalAttributeValue() (dynamodb.AttributeValue, error)
}

func MarshalAttributeValue(data any) (dynamodb.AttributeValue, error) {
	if m, ok := data.(DynamoDBMarshaler); ok {
		return m.MarshalAttributeValue()
	}

	v := reflect.ValueOf(data)
	k := v.Kind()

	switch k {
	case reflect.Invalid, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Complex64, reflect.Complex128:
		return nil, fmt.Errorf("invalid data type")
	case reflect.Bool:
		return &dynamodb.AttributeValueMemberBOOL{Value: v.Bool()}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &dynamodb.AttributeValueMemberN{Value: fmt.Sprintf("%d", v.Int())}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &dynamodb.AttributeValueMemberN{Value: fmt.Sprintf("%d", v.Uint())}, nil
	case reflect.Float32, reflect.Float64:
		return &dynamodb.AttributeValueMemberN{Value: fmt.Sprintf("%f", v.Float())}, nil
	case reflect.Array, reflect.Slice:
		
	case reflect.Map:
		
	case reflect.String:
		return &dynamodb.AttributeValueMemberS{Value: v.String()}, nil
	case reflect.Struct:
	}

	return nil, nil
}
