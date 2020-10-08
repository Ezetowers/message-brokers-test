SHELL := /bin/bash
PWD := $(shell pwd)

GIT_REMOTE = github.com/Ezetowers/nsq-test

default: build

all:

deps:
	go mod tidy
	go mod vendor

build: deps
	GOOS=linux go build -o bin/consumer ${GIT_REMOTE}/consumer
	GOOS=linux go build -o bin/producer ${GIT_REMOTE}/producer
.PHONY: build

docker-image:
	docker build -f ./consumer/Dockerfile -t "consumer:latest" .
	docker build -f ./producer/Dockerfile -t "producer:latest" .
.PHONY: docker-image

docker-compose-up: docker-image
	docker-compose -f docker-compose-dev.yaml up -d --build
.PHONY: docker-compose-up

docker-compose-down:
	docker-compose -f docker-compose-dev.yaml stop -t 1
	docker-compose -f docker-compose-dev.yaml down
.PHONY: docker-compose-down

docker-compose-logs:
	docker-compose -f docker-compose-dev.yaml logs -f
.PHONY: docker-compose-logs
