package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	for _, chunk := range chunkifyEvents(evs) {
		println("Chunk size: ", len(chunk))
		input := dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{},
		}

		for _, ev := range chunk {
			item, err := ev.MarshalDynamoDB()
			if err != nil {
				return events.LambdaFunctionURLResponse{
					StatusCode: 500,
					Body:       "Error marshalling event",
				}
			}
			input.RequestItems[bogieTable] = append(input.RequestItems[bogieTable], types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: item,
				},
			})
		}

		_, err := dynamoDBClient.BatchWriteItem(context.Background(), &input)
		if err != nil {
			return events.LambdaFunctionURLResponse{
				StatusCode: 500,
				Body:       "Error writing to DynamoDB",
			}
		}
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Events received: %d", len(evs)),
	}
}

func chunkifyEvents(evs []event.Event) [][]event.Event {
	chunkSize := 25
	chunks := make([][]event.Event, 0)
	for i := 0; i < len(evs); i += chunkSize {
		end := i + chunkSize

		if end > len(evs) {
			end = len(evs)
		}

		chunks = append(chunks, evs[i:end])
	}

	return chunks
}
