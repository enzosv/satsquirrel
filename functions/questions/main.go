package main

import (
	"context"
	"encoding/json"
	"strings"

	"satsquirrel/pkg/opensat"

	"slices"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// Get question IDs from query parameter
	ids := request.QueryStringParameters["ids"]
	if ids == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing ids parameter",
		}, nil
	}

	questionIDs := strings.Split(ids, ",")

	// Load all questions
	allQuestions, err := opensat.LoadOpenSAT()
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	// Filter questions by ID
	filteredQuestions := make([]opensat.OpenSATQuestion, len(questionIDs))
	for _, topic := range allQuestions {
		for _, question := range topic {
			if slices.Contains(questionIDs, question.ID) {
				filteredQuestions = append(filteredQuestions, question)
			}
		}
	}

	jsonData, err := json.Marshal(filteredQuestions)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		Body: string(jsonData),
	}, nil
}
