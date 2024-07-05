package models

import dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type DDBDocument map[string]dynamodb.AttributeValue
