package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

var tableName = os.Getenv("CODESPACEDB_TABLE_NAME")

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to load configuration"}, err
	}

	id := uuid.New()

	client := dynamodb.NewFromConfig(cfg)
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to put item in DynamoDB"}, err
	}
	
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: id.String(),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
