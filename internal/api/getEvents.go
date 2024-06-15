package api

import (
	"context"
	"encoding/json"
	"net/http"

	lambdaEvents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/seannyphoenix/bogie/internal/event"
)

func getEvents() lambdaEvents.LambdaFunctionURLResponse {
	res, err := dynamoDBClient.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: aws.String(bogieTable),
	})
	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error scanning events",
		}
	}

	var events []event.Event
	for _, item := range res.Items {
		var event event.Event
		err = event.UnmarshalDynamoDB(item)
		if err != nil {
			return lambdaEvents.LambdaFunctionURLResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Error unmarshalling event",
			}
		}
		events = append(events, event)
	}

	body, err := json.Marshal(events)
	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling events",
		}
	}

	return lambdaEvents.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}
}
