GOCMD = go
BINARY_NAME = url-shortner
BIN = ./bin

DOCKER_IMAGE_NAME = thenilesh/url-shortner
DOCKER_IMAGE_TAG ?= latest

build:
	$(GOCMD) build -o $(BIN)/$(BINARY_NAME) -v

test:
	$(GOCMD) test -v ./...

cover:
	$(GOCMD) test -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	$(GOCMD) clean
	rm -rf $(BIN)/

fmt:
	$(GOCMD) fmt ./...

run:
	$(GOCMD) run main.go

docker-build: clean
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

redis-run:
	docker run -d -p 6379:6379 redis

all: test build

.PHONY: build test cover clean fmt run docker-build docker-run docker-push redis-run

