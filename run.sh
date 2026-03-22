#!/bin/bash

ENV=$1

if [ "$ENV" = "prod" ]; then
    echo "Running in production..."
    go run cmd/cli/main.go

elif [ "$ENV" = "dev" ]; then
    echo "Running in development..."
    go run cmd/playground/main.go

elif [ "$ENV" = "build" ]; then
    echo "Building binary..."
    go build -o build/sysmontui cmd/cli/main.go
    echo "Done: ./build/sysmontui"
else
    echo "Usage: ./run.sh [dev|prod|build]"
    exit 1
fi
