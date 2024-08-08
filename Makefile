CONF ?= .env
NAME ?= mhz19go
GIT_REV ?= $(shell git rev-parse --short HEAD)

default: build

.PHONY: build
build: lint test
	CGO_ENABLED=0 go build -o $(NAME)

.PHONY: run
run:
	go run .

.PHONY: tidy
tidy:
	go mod tidy	

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -race -count 1 ./...	
