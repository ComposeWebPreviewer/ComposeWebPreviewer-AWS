package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var tableName = os.Getenv("CODESPACEDB_TABLE_NAME")

type details struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Wasm string `json:"wasm"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to load configuration"}, err
	}

	id := request.QueryStringParameters["id"]

	client := dynamodb.NewFromConfig(cfg)
	output, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	detailsObject := &details{
		Id:   id,
		Code: output.Item["code"].(*types.AttributeValueMemberS).Value,
		Wasm: output.Item["wasm"].(*types.AttributeValueMemberS).Value,
	}
	detailsJson, err := json.Marshal(detailsObject)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: string(detailsJson),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
