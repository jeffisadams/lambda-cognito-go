package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleUnauthenticatedRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "{\"message\":\"This endpoint does not require any authentication\",\"structure\":\"This field was added just to prove it's not an error\"}",
		// This is important as part of the CORS config.
		// Again you should know the security implications of CORS before implementing this
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
		},
	}, nil
}

func main() {
	lambda.Start(handleUnauthenticatedRequest)
}
