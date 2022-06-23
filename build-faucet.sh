#!/bin/bash

docker run --rm -v $(pwd):/app golang:1.18 sh -c 'cd /app && GOOS=linux GOARCH=amd64 go build -o build/faucet .'
