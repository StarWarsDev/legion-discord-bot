BINARY_NAME=legion-discord-bot
DOCKER_IMAGE_NAME=stevegood/legion-discord-bot
TOKEN=your_mom_has_a_token
VERSION=SNAPSHOT

DOCKERCMD=docker
DOCKERBUILD=$(DOCKERCMD) build

GITCMD=git
GITCLONE=$(GITCMD) clone

GOCMD=go
GOBUILD=$(GOCMD) build -v -o $(BINARY_NAME)

all: clean build

build:
	$(GOBUILD) .

run:
	DISCORD_TOKEN=$(TOKEN) LEGION_DATA_VERSION=$(LEGION_DATA_VERSION) ./$(BINARY_NAME)

test:
	CGO_ENABLED=0 go test ./...

up:
	docker-compose up --build

docker:
	$(DOCKERBUILD) --build-arg LEGION_DATA_VERSION=$(LEGION_DATA_VERSION) -t $(DOCKER_IMAGE_NAME):$(VERSION) . && \
	$(DOCKERBUILD) --build-arg LEGION_DATA_VERSION=$(LEGION_DATA_VERSION) -t $(DOCKER_IMAGE_NAME):latest .

docker-push:
	$(DOCKERCMD) push $(DOCKER_IMAGE_NAME):$(VERSION) && \
	$(DOCKERCMD) push $(DOCKER_IMAGE_NAME):latest

clean:
	rm -f ./$(BINARY_NAME)
