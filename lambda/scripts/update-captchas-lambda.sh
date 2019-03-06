#!/bin/bash

. scripts/_config.sh

echo "Build code"
go build cmd/captchas/main.go
zip captchas.zip main
rm main

echo "Create lambda"
aws lambda update-function-code \
    --region ${AWS_DEFAULT_REGION} \
    --function-name ${LAMBDA_CAPTCHAS} \
    --zip-file fileb://captchas.zip \
