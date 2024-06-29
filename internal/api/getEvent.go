package api

import (
	"context"
	"encoding/json"

	lambdaEvents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/models"
	"github.com/google/uuid"
)

func getEvent(id string) lambdaEvents.LambdaFunctionURLResponse {
	thisUUID, err := uuid.Parse(id)
	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Invalid UUID",
		}
	}

	item, err := dynamoDBClient.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(bogieTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberB{
				Value: thisUUID[:],
			},
		},
	})
	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error reading from DynamoDB",
		}
	}

	var event models.Event
	err = event.UnmarshalDynamoDB(item.Item)
	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error unmarshaling event",
		}
	}

	body, err := json.Marshal(event)
	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error marshaling event",
		}
	}

	return lambdaEvents.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(body),
	}
}
