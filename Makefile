CONF ?= .env
include $(CONF)
include Makefile.docker.mk
include Makefile.load.mk

.DEFAULT_GOAL := default

NAME ?= mhz19-go-backend
GIT_REV ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date +%FT%T)
OS_NAME := $(shell uname -s | tr A-Z a-z)

default: lint test build

build:
	GORACE="halt_on_error=1" CGO_ENABLED=1 go build -race -o ./bin/backend ./cmd/backend

build-norace:
	CGO_ENABLED=1 go build -o ./bin/backend ./cmd/backend

clean:
	rm ./bin/backend

run:
	GORACE="halt_on_error=1" go run -race ./cmd/backend

run-norace:
	go run ./cmd/backend

provision:
	go run ./cmd/provision

tidy:
	go mod tidy	

lint:
	golangci-lint run

migrate-reset: migrate-down migrate-up

migrate-down:
	export DB_REV=03 && make migrate-down-single
	export DB_REV=02 && make migrate-down-single
	export DB_REV=01 && make migrate-down-single
	export DB_REV=00 && make migrate-down-single

migrate-up:
	export DB_REV=00 && make migrate-up-single
	export DB_REV=01 && make migrate-up-single
	export DB_REV=02 && make migrate-up-single
	export DB_REV=03 && make migrate-up-single

migrate-up-single:
	sqlite3 ./sqlite/database.bin < ./sql/$(DB_REV)-up.sql

migrate-down-single:
	sqlite3 ./sqlite/database.bin < ./sql/$(DB_REV)-down.sql

migrate-dump:
	sqlite3 ./sqlite/database.bin .dump > ./sql/$(DATE)-dump.sql

test:
	go test -cover -race -count 1 ./...

bench:
	rm -f *.prof && go test ./internal/engine -bench=^Benchmark10$$ -run=^$$ -benchmem -cpuprofile cpu.prof -memprofile=mem.prof

pprof-mem:
	go tool pprof -http=:7272 mem.prof

pprof-cpu:
	go tool pprof -http=:7272 cpu.prof

test-one:
	go test ./internal/engine -run TestMappings -v

update:
	git -C . pull && git -C ../mhz19-front pull && git -C ../device-pinger pull && make compose-build && make compose-up

# git pull && make docker-build && make docker-down && make docker-up && make docker-logs
# rm -f *.prof && go test ./internal/entities/rules -bench=^Benchmark20$$ -run=^$$ -benchmem -cpuprofile cpu.prof -memprofile=mem.prof
# rm -f *.prof && go test ./internal/counters -bench=^Benchmark30$$ -run=^$$ -benchmem -cpuprofile cpu.prof -memprofile=mem.prof
# CGO_ENABLED=1 go test -benchmem ./...
# CGO_ENABLED=1 go test -cover -race -count 1 ./...
# curl -X PUT -H "Content-Type: application/json" -d @assets/push-message.json $(REST_API_URL)/push-message
# go test -bench Benchmark10 -run=^$ -benchmem ./internal/counters -cpuprofile=cpu.prof
# go test -bench Benchmark10 -run=^$ ./internal/counters -memprofile=mem.prof
# go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out
# go test -bench=^Benchmark10$ -run=^$ . -memprofile=mem.prof ./internal/counters
# go test -benchmem -run=^$ -bench "^Benchmark10$" github.com/fedulovivan/mhz19-go/internal/counters -cpuprofile=cpuprofile.prof
# go test ./internal/counters -bench=^Benchmark20$$ -run=^$$ 
# https://gist.github.com/ungoldman/11282441
# https://stackoverflow.com/questions/978142/how-to-benchmark-apache-with-delays
# oha --method PUT -H 'Content-Type: application/json' -D ./assets/load/create-rule.json $(REST_API_URL)/rules
# watch -n 0.5 