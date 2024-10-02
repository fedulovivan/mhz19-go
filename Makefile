CONF ?= .env
include $(CONF)

NAME ?= mhz19-go-backend
GIT_REV ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date +%FT%T)
REST_API_URL ?= http://localhost:$(REST_API_PORT)$(REST_API_PATH)
API_LOAD_COUNT ?= 1000
API_LOAD_THREADS ?= 10
OS_NAME := $(shell uname -s | tr A-Z a-z)

default: lint test build

.PHONY: build
build:
	CGO_ENABLED=1 go build -o ./bin/backend ./cmd/backend

.PHONY: docker-build
docker-build:
	DOCKER_CLI_HINTS=false docker build --label "git.revision=${GIT_REV}" --tag $(NAME) .

.PHONY: docker-down
docker-down:
	docker stop $(NAME) && docker rm $(NAME)

.PHONY: docker-up
docker-up:
ifeq ($(OS_NAME), linux)
	docker run -d --env-file=$(CONF) -v ./database.bin:/database.bin --network=host --device /dev/snd:/dev/snd --name=$(NAME) $(NAME)
else
	docker run -d --env-file=$(CONF) -v ./database.bin:/database.bin --network=host --name=$(NAME) $(NAME)
endif
	
# docker run -d --env-file=$(CONF) -v ./database.bin:/database.bin -p $(REST_API_PORT):$(REST_API_PORT) --name=$(NAME) --network=host $(NAME)
	
.PHONY: docker-logs
docker-logs:
	docker logs --follow $(NAME)

.PHONY: docker-logs-save
docker-logs-save:
	docker logs --timestamps $(NAME) 2>&1 | cat > log.txt

.PHONY: clean
clean:
	rm ./bin/backend

.PHONY: api-load-read-rules
api-load-read-rules:
	ab -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules

.PHONY: api-load-write-rules
api-load-write-rules:
	ab -T application/json -u ./assets/create-rule.json -n $(API_LOAD_COUNT) -c $(API_LOAD_THREADS) $(REST_API_URL)/rules

.PHONY: api-load-push-messages
api-load-push-messages:
	hey -T application/json -m PUT -D ./assets/push-message.json -n 1000 -c 10 -q 50 $(REST_API_URL)/push-message

# https://stackoverflow.com/questions/978142/how-to-benchmark-apache-with-delays
# https://gist.github.com/ungoldman/11282441
# watch -n 0.5 
# .PHONY: api-load-write-3
# api-load-write-3:
# 	curl -X PUT -H "Content-Type: application/json" -d @assets/push-message.json $(REST_API_URL)/push-message

.PHONY: api-load-once
api-load-once:
	wget -O /dev/null $(REST_API_URL)/rules

.PHONY: run
run:
	go run ./cmd/backend

.PHONY: provision
provision:
	go run ./cmd/provision

.PHONY: tidy
tidy:
	go mod tidy	

.PHONY: lint
lint:
	golangci-lint run

.PHONY: migrate-reset
migrate-reset: migrate-down migrate-up

.PHONY: migrate-down
migrate-down:
	export DB_REV=01 && make migrate-down-single
	export DB_REV=00 && make migrate-down-single

.PHONY: migrate-up
migrate-up:
	export DB_REV=00 && make migrate-up-single
	export DB_REV=01 && make migrate-up-single

.PHONY: migrate-up-single
migrate-up-single:
	sqlite3 ./database.bin < ./sql/$(DB_REV)-up.sql

.PHONY: migrate-down-single
migrate-down-single:
	sqlite3 ./database.bin < ./sql/$(DB_REV)-down.sql

.PHONY: migrate-dump
migrate-dump:
	sqlite3 ./database.bin .dump > ./sql/$(DATE)-dump.sql

.PHONY: test
test:
	CGO_ENABLED=1 go test -cover -race ./...
# CGO_ENABLED=1 go test -cover -race -count 1 ./...

# .PHONY: bench
# bench:
# 	CGO_ENABLED=1 go test -benchmem ./...

.PHONY: test-one
test-one:
	go test -v github.com/fedulovivan/mhz19-go/internal/engine -run "TestMappings"
