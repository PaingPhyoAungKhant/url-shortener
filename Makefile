.PHONY: run test lint

BUILD_DIR := bin
APP := omnipocket
GOOS := linux
GOARCH := amd64

run:
	go run ./cmd/server/main.go

test: 
	go test ./...

lint:
	golangci-lint run ./...

build:
	GOOS=$(GOOS) GOARCH=${GOARCH} go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP) ./cmd/$(APP)



