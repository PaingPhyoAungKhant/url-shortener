.PHONY: run test lint format build

BUILD_DIR := bin
APP := url-shortener
GOOS := linux
GOARCH := amd64

run:
	go run ./cmd/api/main.go

test: 
	go test ./...

lint:
	golangci-lint run ./...

format:
	golangci-lint fmt ./...

build:
	GOOS=$(GOOS) GOARCH=${GOARCH} go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP) ./cmd/api/main.go



