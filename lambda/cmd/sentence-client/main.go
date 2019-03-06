package main

import (
	"log"

	"github.com/djhworld/go-lambda-invoke/golambdainvoke"
)

func main() {
	response, err := golambdainvoke.Run(golambdainvoke.Input{
		Port: 8001,
		Payload: map[string]string{
			"name": "Bob",
		},
	})

	if err != nil {
		log.Println(err)
	}

	log.Println(string(response))
}
