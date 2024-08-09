CONF ?= .env
NAME ?= backend
GIT_REV ?= $(shell git rev-parse --short HEAD)

default: build-backend

.PHONY: build-backend
build-backend: lint test
	CGO_ENABLED=0 go build -o ./bin/$(NAME) ./cmd/$(NAME)

.PHONY: run
run:
	go run ./cmd/$(NAME)

.PHONY: tidy
tidy:
	go mod tidy	

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -race -count 1 ./...	
