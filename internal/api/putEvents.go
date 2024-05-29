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

func putEvents(ev PutEventsRequest) events.LambdaFunctionURLResponse {
	var evs []event.Event
	body := ev.Body

	err := json.Unmarshal([]byte(body), &evs)

	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       "Error unmarshalling events",
		}
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Events received: %d", len(evs)),
	}
}
