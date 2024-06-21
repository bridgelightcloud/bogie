package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bridgelightcloud/bogie/internal/api"
	_ "golang.org/x/crypto/x509roots/fallback"
)

func main() {
	api.Setup()
	lambda.Start(api.Route)
}
