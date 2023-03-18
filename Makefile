SHELL := /bin/bash

.DEFAULT_GOAL := build

IMAGE ?= ghcr.io/thevilledev/learn-k8s-admission-webhooks:latest

.PHONY: build
build: $(shell find . -name '*.go')
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o $@ ./cmd/webhook

.PHONY: build-local
build-local: $(shell find . -name '*.go')
        CGO_ENABLED=0 go build -ldflags="-s -w" -o $@ ./cmd/webhook

.PHONY: run-local
run-local: build-local
		LISTEN_ADDR=":8443" \
		CERT_PATH="third_party/tls/local/webhook-server-tls.crt" \
		KEY_PATH="third_party/tls/local/webhook-server-tls.key" \
		./webhook

.PHONY: docker-image
docker-image:
	docker build -f build/package/docker/Dockerfile -t $(IMAGE) .

.PHONY: push-image
push-image: docker-image
	docker push $(IMAGE)
