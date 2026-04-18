.PHONY: build run clean test

APP_NAME=daily_task

build:
	go build -o $(APP_NAME) cmd/main.go

run:
	go run cmd/main.go

clean:
	rm -f $(APP_NAME)
	rm -rf log/

test:
	go test ./...

deps:
	go mod download
	go mod tidy