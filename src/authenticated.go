package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "{\"message\":\"Marshalling a return body is a problem for another day.  But the request was successful\",\"structure\":\"See this is actually not an error\"}",
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
	lambda.Start(handleRequest)
}
