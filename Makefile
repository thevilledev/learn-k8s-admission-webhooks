.DEFAULT_GOAL := cmd/webhook

IMAGE ?= stackrox/admission-controller-webhook-demo:latest

cmd/webhook: $(shell find . -name '*.go')
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $@ ./cmd/webhook

.PHONY: docker-image
docker-image: image/webhook-server
	docker build -t $(IMAGE) image/

.PHONY: push-image
push-image: docker-image
	docker push $(IMAGE)
