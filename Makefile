GO := go
BUILD := build -mod=vendor
GOROOT := $(shell go env GOROOT)
LDFLAGS := -s -w
GCFLAGS := -trimpath=$(CURDIR);$(GOROOT)/src
VERSION := `git describe --tags`
IMAGE_NAME := pizzabot
IMAGE_REGISTRY ?= ghcr.io/razzie
FULL_IMAGE_NAME := $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(VERSION)

.PHONY: build
build:
	$(GO) $(BUILD) -ldflags="$(LDFLAGS)" -gcflags=all="$(GCFLAGS)" ./cmd/pizzabot

.PHONY: docker-build
docker-build:
	docker build . -t $(FULL_IMAGE_NAME)

.PHONY: docker-push
docker-push: docker-build
	docker push $(FULL_IMAGE_NAME)
