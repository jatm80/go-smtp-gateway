#!/bin/bash

## lint
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.64.8 golangci-lint run -v

## build
go build

export SMTP_ADDR="localhost:2525"
export SMTP_DOMAIN="gateway.home"
export SMTP_USER="gateway"
export SMTP_PASS="gateway"
export TELEGRAM_TOKEN=$(pass JT-01/exports/sms-gateway/TELEGRAM_TOKEN)
export TELEGRAM_CHAT_ID=$(pass JT-01/exports/sms-gateway/TELEGRAM_CHAT_ID)

./go-smtp-gateway