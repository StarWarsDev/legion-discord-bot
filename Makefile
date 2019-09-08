BINARY_NAME=legion-discord-bot
DOCKER_IMAGE_NAME=stevegood/legion-discord-bot
TOKEN=your_mom_has_a_token
VERSION=SNAPSHOT
LEGION_DATA_VERSION=master

DOCKERCMD=docker
DOCKERBUILD=$(DOCKERCMD) build

GITCMD=git
GITCLONE=$(GITCMD) clone

GOCMD=go
GOBUILD=$(GOCMD) build -v -o $(BINARY_NAME)

YARNCMD=yarn
YARNJSON=$(YARNCMD) json

all: clean generate-json build

checkout-legion-data:
	$(GITCLONE) https://github.com/andrelind/legion-data.git

checkout-ld-tag: checkout-legion-data
	cd legion-data && \
	$(GITCMD) fetch --all && \
	$(GITCMD) checkout $(LEGION_DATA_VERSION)

generate-json: checkout-ld-tag
	cd legion-data && \
	$(YARNCMD) && \
	$(YARNJSON) && \
	cd .. && \
	mv legion-data/out/legion-data.json ./legion-data.json && \
	rm -rf legion-data

build:
	$(GOBUILD) .

run:
	DISCORD_TOKEN=$(TOKEN) LEGION_DATA_VERSION=$(LEGION_DATA_VERSION) ./$(BINARY_NAME)

docker:
	$(DOCKERBUILD) --build-arg LEGION_DATA_VERSION=$(LEGION_DATA_VERSION) -t $(DOCKER_IMAGE_NAME):$(VERSION) . && \
	$(DOCKERBUILD) --build-arg LEGION_DATA_VERSION=$(LEGION_DATA_VERSION) -t $(DOCKER_IMAGE_NAME):latest .

docker-push:
	$(DOCKERCMD) push $(DOCKER_IMAGE_NAME):$(VERSION) && \
	$(DOCKERCMD) push $(DOCKER_IMAGE_NAME):latest

clean:
	rm -f ./$(BINARY_NAME)
	rm -f ./legion-data.json
	rm -rf ./legion-data
