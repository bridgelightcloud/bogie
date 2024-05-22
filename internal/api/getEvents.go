package api

import (
	"encoding/json"

	lambdaEvents "github.com/aws/aws-lambda-go/events"
	"github.com/seannyphoenix/bogie/pkg/event"
)

func getEvents() lambdaEvents.LambdaFunctionURLResponse {
	events := event.GetExampleEventArray(3)
	body, err := json.Marshal(events)

	if err != nil {
		return lambdaEvents.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error marshalling events",
		}
	}

	return lambdaEvents.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(body),
	}
}
