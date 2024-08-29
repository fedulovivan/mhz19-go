CONF ?= .env
NAME ?= backend
GIT_REV ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date +%FT%T)
NUM_MIGRATION ?= 00
API_URL ?= http://localhost:8888/rules
API_LOAD_COUNT ?= 1000
API_LOAD_THREADS ?= 10

default: lint test build

.PHONY: build
build:
	CGO_ENABLED=1 go build -o ./bin/$(NAME) ./cmd/$(NAME)

.PHONY: clean
clean:
	rm ./bin/$(NAME)

.PHONY: api-load-read
api-load-read:
	ab -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(API_URL)

.PHONY: api-load-write
api-load-write:
	ab -T application/json -u ./assets/create.json -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(API_URL)

.PHONY: api-load-once
api-load-once:
	wget -O /dev/null $(API_URL)

.PHONY: run
run:
	go run ./cmd/$(NAME)

.PHONY: tidy
tidy:
	go mod tidy	

.PHONY: lint
lint:
	golangci-lint run

.PHONY: migrate-reset
migrate-reset: migrate-down migrate-up

.PHONY: migrate-up
migrate-up:
	sqlite3 ./database.bin < ./sql/$(NUM_MIGRATION)-up.sql

.PHONY: migrate-down
migrate-down:
	sqlite3 ./database.bin < ./sql/$(NUM_MIGRATION)-down.sql

.PHONY: migrate-dump
migrate-dump:
	sqlite3 ./database.bin .dump > ./sql/$(DATE)-dump.sql

.PHONY: test
test:
	SQLITE_FILENAME=$(PWD)/database_ut.bin go test -cover -race -count 1 ./...


