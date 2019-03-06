package main

import (
	"bytes"
	"context"
	"image/color"
	"image/png"
	"log"

	"encoding/base64"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/afocus/captcha"
	"github.com/google/uuid"

	"helloLambda/internal/font"
)

type Request struct{}

type Response struct {
	ID      string `json:"id"`
	Captcha string `json:"captcha"`
}

type Captcha struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

func HandleRequest(ctx context.Context, request Request) (interface{}, error) {
	log.Println("HandleRequest")

	var img *captcha.Image //NOTE: *image.RGBA
	var str string

	//Generate image
	img, str = cap.Create(6, captcha.ALL)
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		log.Println(err.Error())

		return "Buu", err
	}

	//Save image to DB
	captcha := Captcha{
		ID:     uuid.New().String(),
		Secret: str,
	}
	av, err := dynamodbattribute.MarshalMap(captcha)
	if err != nil {
		log.Println(err.Error())

		return "Buu", err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Captchas"),
	}
	_, err = dynamo.PutItem(input)
	if err != nil {
		log.Println(err.Error())

		return "Buu", err
	}

	return Response{
		ID:      captcha.ID,
		Captcha: base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}

//NOTE: "global"
var cap *captcha.Captcha
var dynamo *dynamodb.DynamoDB

func main() {
	log.Println("Startup")

	log.Println("Load font")
	fontBytes, err := base64.StdEncoding.DecodeString(font.FONT)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Setup captcha service")
	cap = captcha.New() //NOTE set "global" value
	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	err = cap.AddFontFromBytes(fontBytes)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Setup DynamoDB connection")
	dynamo = dynamodb.New(session.New())

	lambda.Start(HandleRequest)
}
