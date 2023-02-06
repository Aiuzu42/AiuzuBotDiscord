#!/bin/bash
echo "Formatting code...\n"
go fmt ./...
echo "Running linter...\n"
golangci-lint run
echo "Running tests...\n"
go test ./...
echo "Done...\n"