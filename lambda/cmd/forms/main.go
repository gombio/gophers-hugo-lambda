package main

import (
	"context"
	"errors"
	"image/color"
	"log"

	"encoding/base64"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/afocus/captcha"

	"helloLambda/internal/font"

	"github.com/google/uuid"
)

type Request struct {
	Data    interface{} `json:"data"`
	Captcha Captcha     `json:"captcha"`
}

type Response struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type Captcha struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

type Form struct {
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

func FindCaptcha(captcha_id string) (*Captcha, error) {
	input := &dynamodb.QueryInput{
		Limit:     aws.Int64(1),
		TableName: aws.String("Captchas"),
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(captcha_id),
					},
				},
			},
		},
	}
	result, err := dynamo.Query(input)
	if err != nil {
		return nil, err
	}
	captchas := []Captcha{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &captchas)
	if err != nil {
		return nil, err
	}

	if len(captchas) == 0 {
		return nil, errors.New("Buu: ni ma takiej captchy")
	}

	captcha := captchas[0]

	return &captcha, nil
}

func HandleRequest(ctx context.Context, request Request) (interface{}, error) {
	log.Println("HandleRequest")

	//Find existing captcha
	captcha, err := FindCaptcha(request.Captcha.ID)
	if err != nil {
		return Response{
			ID:      "",
			Message: err.Error(),
		}, nil
	}

	//Validate captcha
	if request.Captcha.Secret != captcha.Secret {
		return Response{
			ID:      "",
			Message: "Buu: a umisz czytac? Skup sie!",
		}, nil
	}

	//Save form to DynamoDB
	form := Form{
		ID:   uuid.New().String(),
		Data: request.Data,
	}
	fv, err := dynamodbattribute.MarshalMap(form)
	if err != nil {
		return Response{
			ID:      "",
			Message: err.Error(),
		}, nil
	}
	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("Forms"),
		Item:      fv,
	})
	if err != nil {
		return Response{
			ID:      "",
			Message: err.Error(),
		}, nil
	}

	//Remove used captcha
	_, err = dynamo.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("Captchas"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(request.Captcha.ID),
			},
			"secret": {
				S: aws.String(request.Captcha.Secret),
			},
		},
	})
	if err != nil {
		return Response{
			ID:      "",
			Message: err.Error(),
		}, nil
	}

	return Response{
		ID:      form.ID,
		Message: "Yay",
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
