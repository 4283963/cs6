#!/bin/bash
cd "$(dirname "$0")/backend"
go mod tidy
go run main.go
