.PHONY: build dev test lint migrate-up migrate-down jwt clean

APP_NAME := gobase
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) .

dev:
	go run . serve

test:
	go test ./... -v

lint:
	golangci-lint run

migrate-up:
	go run . migrate up

migrate-down:
	go run . migrate down

jwt:
	go run . jwt generate --user-id=$(id)

clean:
	rm -rf $(BUILD_DIR)
