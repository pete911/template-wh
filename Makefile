#!make
SHELL := /bin/bash
.SHELLFLAGS := -ec
IMAGE := pete911/template-wh
VERSION ?= dev

VERSION ?= dev

test:
	$(info *** [go test] ***)
	go clean -testcache && go test -cover ./...

e2e-test:
	$(info *** [end to end tests] ***)
	./e2e/e2e

build: test
	$(info *** [go build] ***)
	go build -mod vendor

image:
	docker build -t ${IMAGE}:${VERSION} .
	docker tag ${IMAGE}:${VERSION} ${IMAGE}:latest

push-image:
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest
