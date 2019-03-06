#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o main cmd/hello-world/main.go
zip hello-world.zip main
rm main
