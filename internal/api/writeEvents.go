package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bridgelightcloud/bogie/internal/db"
	"github.com/bridgelightcloud/bogie/internal/models"
	"github.com/bridgelightcloud/bogie/internal/util"
)

type PutEventsRequest struct {
	Body string
}

func writeEvents(r PutEventsRequest) events.LambdaFunctionURLResponse {
	var evs []models.Event
	body := r.Body

	err := json.Unmarshal([]byte(body), &evs)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Error unmarshaling events",
		}
	}

	for _, chunk := range util.ChunkifySlice(evs, db.DynamoDBBatchWriteLimit) {
		println("Chunk size: ", len(chunk))
		input := dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{},
		}

		for _, ev := range chunk {
			item, err := ev.MarshalDynamoDB()
			if err != nil {
				println("Error marshaling event: ", err.Error())
				return events.LambdaFunctionURLResponse{
					StatusCode: 500,
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
			println("Error writing to DynamoDB: ", err.Error())
			return events.LambdaFunctionURLResponse{
				StatusCode: 500,
			}
		}
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Events received: %d", len(evs)),
	}
}
