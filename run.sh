#!/bin/bash

ENV=$1

if [ "$ENV" = "prod" ]; then
    echo "Running in production..."
    go run cmd/cli/main.go
elif [ "$ENV" = "dev" ]; then
    echo "Running in development..."
    go run cmd/playground/main.go
else
    echo "Usage: ./run.sh [dev|prod]"
    exit 1
fi
