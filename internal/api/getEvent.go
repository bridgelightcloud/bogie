package api

import (
	"encoding/json"

	lambdaEvents "github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/seannyphoenix/bogie/pkg/event"
)

func getEvent(id string) lambdaEvents.LambdaFunctionURLResponse {
	thisUUID, err := uuid.Parse(id)
	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: 404,
			Body:       "Invalid UUID",
		}
	}

	event := event.GetExampleEvent(thisUUID)
	body, err := json.Marshal(event)

	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error marshalling event",
		}
	}

	return lambdaEvents.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(body),
	}
}
