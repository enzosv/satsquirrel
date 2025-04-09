package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Your server-side functionality
	questions, err := loadOpenSAT("OpenSAT.json")
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	topics := map[string]int{"math": 2, "english": 2}

	shuffled := shuffleSubset(ctx, questions, topics)

	jsonData, err := json.Marshal(shuffled)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json; charset=utf-8"},
		Body:       string(jsonData),
	}, nil
}

func main() {
	lambda.Start(handler)
}
