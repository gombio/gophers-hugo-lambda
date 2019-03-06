#!/bin/bash

. scripts/_config.sh

echo "Build code"
GOOS=linux GOARCH=amd64 go build cmd/forms/main.go
zip forms.zip main
rm main

echo "Create lambda"
aws lambda create-function \
    --region ${AWS_DEFAULT_REGION} \
    --function-name ${LAMBDA_FORMS} \
    --runtime go1.x \
    --handler main \
    --memory-size 128 \
    --zip-file fileb://forms.zip \
    --role $ROLE_ARN
