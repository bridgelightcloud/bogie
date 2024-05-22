package api

import "github.com/aws/aws-lambda-go/events"

func dummy() events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{}
}
