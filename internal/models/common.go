package models

import dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type DBDocument map[string]dynamodb.AttributeValue
