GO_RUN:=go run
GO_BUILD:=go build

DOCKER_RUN:=docker run
DOCKER_BUILD:=docker build

CONTAINER_NAME:=log-server
TAG:=0.1

PHONY: help run build

help:
	@echo Makefile Command Reference
	@echo wip

run:
	$(GO_RUN) .

build:
	$(GO_BUILD) -o serverd -ldflags '-s -w'

docker_run: docker_build
	$(DOCKER_RUN) -p 1323:1323 --name $(CONTAINER_NAME) $(CONTAINER_NAME):$(TAG)

docker_build:
	$(DOCKER_BUILD) --tag $(CONTAINER_NAME):$(TAG) .
