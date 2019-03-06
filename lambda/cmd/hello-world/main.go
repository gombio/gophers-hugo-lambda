package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {
	log.Println("HANDLE")

	return "Hello World!", nil
}

func main() {
	log.Println("START")

	lambda.Start(HandleRequest)
}
