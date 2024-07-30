#!/bin/sh

set -e

echo "Building Go project..."
go build -v -o ./bin/main ./cmd/main.go
./bin/main