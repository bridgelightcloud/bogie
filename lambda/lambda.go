package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/seannyphoenix/bogie/internal/api"
)

func main() {
	lambda.Start(api.Route)
}
