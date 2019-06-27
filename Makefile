BINARY_NAME=legion-discord-bot

GOCMD=go
GOBUILD=$(GOCMD) build -v -o $(BINARY_NAME)

all: build

build:
	$(GOBUILD) .

clean:
	rm -rf ./$(BINARY_NAME)
