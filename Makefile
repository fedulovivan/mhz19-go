CONF ?= .env
include $(CONF)
include Makefile.docker.mk

NAME ?= mhz19-go-backend
GIT_REV ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date +%FT%T)
REST_API_URL ?= http://localhost:$(REST_API_PORT)$(REST_API_PATH)
API_LOAD_COUNT ?= 1000
API_LOAD_THREADS ?= 10
API_RPS ?= 50
OS_NAME := $(shell uname -s | tr A-Z a-z)

default: lint test build

build:
	GORACE="halt_on_error=1" CGO_ENABLED=1 go build -race -o ./bin/backend ./cmd/backend

clean:
	rm ./bin/backend

api-load-rules-read:
	oha -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules

api-load-stats-read:
	oha -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/stats

api-load-rules-write:
	curl -H 'Content-Type: application/json' -X PUT -d '{"deviceId":"DeviceId(0x00158d0004244bda)","deviceClass":1}' $(REST_API_URL)/devices
	curl -H 'Content-Type: application/json' -X PUT -d '{"deviceId":"DeviceId(10011cec96)","deviceClass":6}' $(REST_API_URL)/devices
	oha --method PUT -H 'Content-Type: application/json' -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) -D ./assets/load/create-rule.json --rand-regex-url $(REST_API_URL)/rules/name-[a-z0-9]{16}

api-load-push-message-write:
	oha --method PUT -H 'Content-Type: application/json' -D ./assets/load/push-message.json -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) -q $(API_RPS) $(REST_API_URL)/push-message

api-load-once:
	wget -O /dev/null $(REST_API_URL)/rules

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
	export DB_REV=01 && make migrate-down-single
	export DB_REV=00 && make migrate-down-single

migrate-up:
	export DB_REV=00 && make migrate-up-single
	export DB_REV=01 && make migrate-up-single

migrate-up-single:
	sqlite3 ./database.bin < ./sql/$(DB_REV)-up.sql

migrate-down-single:
	sqlite3 ./database.bin < ./sql/$(DB_REV)-down.sql

migrate-dump:
	sqlite3 ./database.bin .dump > ./sql/$(DATE)-dump.sql

test:
	go test -cover -race -count 1 ./...

bench:
	rm -f *.prof && go test ./internal/counters -bench=^Benchmark30$$ -run=^$$ -benchmem -cpuprofile cpu.prof -memprofile=mem.prof

pprof-mem:
	go tool pprof mem.prof

pprof-cpu:
	go tool pprof cpu.prof

test-one:
	go test ./internal/engine -run TestMappings -v

# ab -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules
# ab -T application/json -u ./assets/load/create-rule.json -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules
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
# oha --method PUT -H 'Content-Type: application/json' -d "{\"name\":\"name-`uuidgen`\"}" -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules
# watch -n 0.5 