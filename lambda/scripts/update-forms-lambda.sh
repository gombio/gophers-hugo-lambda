#!/bin/bash

. scripts/_config.sh

echo "Build code"
go build cmd/forms/main.go
zip forms.zip main
rm main

echo "Create lambda"
aws lambda update-function-code \
    --region ${AWS_DEFAULT_REGION} \
    --function-name ${LAMBDA_FORMS} \
    --zip-file fileb://forms.zip \
