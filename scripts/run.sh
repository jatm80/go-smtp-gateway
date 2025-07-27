#!/bin/bash

## lint
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.64.8 golangci-lint run -v

## build
go build