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

YARNCMD=yarn
YARNJSON=$(YARNCMD) json

all: clean generate-json build

checkout-legion-data:
	$(GITCLONE) https://github.com/andrelind/legion-data.git

generate-json: checkout-legion-data
	cd legion-data && \
	$(YARNCMD) && \
	$(YARNJSON) && \
	cd .. && \
	mv legion-data/out/legion-data.json ./legion-data.json && \
	rm -rf legion-data

build:
	$(GOBUILD) .

run:
	./$(BINARY_NAME) -t $(TOKEN)

docker:
	$(DOCKERBUILD) --squash -t $(DOCKER_IMAGE_NAME):$(VERSION) .

clean:
	rm -f ./$(BINARY_NAME)
	rm -f ./legion-data.json
	rm -rf ./legion-data
