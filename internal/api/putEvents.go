package api

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/seannyphoenix/bogie/internal/event"
)

type PutEventsRequest struct {
	Body string
}

func putEvents(r PutEventsRequest) events.LambdaFunctionURLResponse {
	var evs []event.Event
	body := r.Body

	err := json.Unmarshal([]byte(body), &evs)

	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Error unmarshalling events",
		}
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Events received: %d", len(evs)),
	}
}
