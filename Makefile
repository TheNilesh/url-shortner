GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOCLEAN = $(GOCMD) clean
GOFMT = $(GOCMD) fmt
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

BINARY_NAME = url-shortner
BIN = ./bin

DOCKER_IMAGE_NAME = thenilesh/url-shortner
DOCKER_IMAGE_TAG ?= latest

build:
	$(GOBUILD) -o $(BIN)/$(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BIN)/

fmt:
	$(GOFMT) ./...

run:
	$(GORUN) main.go

docker-build: clean
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

redis-run:
	docker run -d -p 6379:6379 redis

all: test build

.PHONY: build test clean fmt run docker-build docker-run docker-push redis-run

