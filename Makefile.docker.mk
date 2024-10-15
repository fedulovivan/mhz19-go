#
# Here are all docker and docker-compose related targets 
# which are extracted into separate file and included into main Makefile.
#

build-norace:
	CGO_ENABLED=1 go build -o ./bin/backend ./cmd/backend

docker-build:
	DOCKER_CLI_HINTS=false docker build --label "git.revision=${GIT_REV}" --tag $(NAME) .

docker-down:
	docker stop $(NAME) && docker rm $(NAME)

docker-up:
ifeq ($(OS_NAME), linux)
	docker run --detach --restart=always --env-file=$(CONF) -v ./database.bin:/app/database.bin --network=host --device /dev/snd:/dev/snd --name=$(NAME) $(NAME)
else
	docker run --detach --restart=always --env-file=$(CONF) -v ./database.bin:/app/database.bin -p $(REST_API_PORT):$(REST_API_PORT) --name=$(NAME) $(NAME)
endif
	
docker-logs:
	docker logs --follow $(NAME)

docker-shell:
	docker exec -it $(NAME) /bin/sh

docker-dive:
	_dive $(NAME)

docker-logs-save:
	docker logs --timestamps $(NAME) 2>&1 | cat > log.txt
	
docker-stats:
	docker stats $(NAME)

compose-up:
	NAME=$(NAME) docker compose up --no-build --detach

compose-down:
	NAME=$(NAME) docker compose down

compose-up-dev:
	NAME=$(NAME) docker compose up --no-build --detach prometheus grafana

compose-down-dev:
	NAME=$(NAME) docker compose down prometheus grafana

update:
	git pull && make docker-build && make docker-down && make docker-up && make docker-logs

# utility command for getting shell in the "Docker Desktop"s linux vm on mac
# borrowed from https://gist.github.com/BretFisher/5e1a0c7bcca4c735e716abf62afad389
macos-docker-shell:
	docker run -it --rm --privileged --pid=host justincormack/nsenter1	
