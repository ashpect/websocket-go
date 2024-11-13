#!/bin/bash

if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go to continue."
    exit 1
fi

echo "Tidying Go modules..."
go mod tidy

echo "Building the application..."
go build -o websocket-go .

echo "Running the application..."
./websocket-go