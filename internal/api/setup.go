package api

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamoDBClient *dynamodb.Client
var bogieTable string

func Setup() {

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	dynamoDBClient = dynamodb.NewFromConfig(cfg)

	bogieTable = os.Getenv("BOGIE_TABLE")
}
