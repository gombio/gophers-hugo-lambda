package main

import (
	"context"
	"log"

	// "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Sentence string `json:"sentence"`
}

func HandleRequest(ctx context.Context, request Request) (interface{}, error) {
	log.Println("HandleRequest")

	r := Response{
		Sentence: "Hell " + request.Name,
	}

	return r, nil
}

func main() {
	log.Println("Startup")

	lambda.Start(HandleRequest)
}
