#!/bin/bash

. scripts/_config.sh

echo "Build code"
GOOS=linux GOARCH=amd64 go build cmd/captchas/main.go
zip captchas.zip main
rm main

echo "Create lambda"
aws lambda create-function \
    --region ${AWS_DEFAULT_REGION} \
    --function-name ${LAMBDA_CAPTCHAS} \
    --runtime go1.x \
    --handler main \
    --memory-size 128 \
    --zip-file fileb://captchas.zip \
    --role $ROLE_ARN
