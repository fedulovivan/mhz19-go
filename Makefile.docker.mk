#
# Here are all docker and docker-compose related targets 
# which are extracted into separate file and included into main Makefile.
#

docker-build:
	DOCKER_CLI_HINTS=false docker build --label "git.revision=${GIT_REV}" --tag $(NAME) .

docker-down:
	docker stop $(NAME) && docker rm $(NAME)-1

docker-up:
ifeq ($(OS_NAME), linux)
	docker run --detach --restart=always --env-file=$(CONF) -v ./sqlite:/app/sqlite --network=host --device /dev/snd:/dev/snd --name=$(NAME)-1 $(NAME)
else
	docker run --detach --restart=always --env-file=$(CONF) -v ./sqlite:/app/sqlite -p $(REST_API_PORT):$(REST_API_PORT) --name=$(NAME)-1 $(NAME)
endif
	
docker-logs:
	docker logs --follow $(NAME)-1

docker-shell:
	docker exec -it $(NAME)-1 /bin/sh

docker-dive:
	_dive $(NAME)

docker-logs-save:
	docker logs --timestamps $(NAME)-1 2>&1 | cat > log.txt
	
docker-stats:
	docker stats $(NAME)-1

compose-build:
	docker compose build

compose-up:
	docker compose up --no-build --detach

compose-down:
	docker compose down

# compose-up-dev:
# 	NAME=$(NAME) docker compose up --no-build --detach prometheus grafana

# compose-down-dev:
# 	NAME=$(NAME) docker compose down prometheus grafana

# utility command for getting shell in the "Docker Desktop"s linux vm on mac
# borrowed from https://gist.github.com/BretFisher/5e1a0c7bcca4c735e716abf62afad389
macos-docker-shell:
	docker run -it --rm --privileged --pid=host justincormack/nsenter1	
