.PHONY: build test fmt
GO      ?= go
PKG := ./pkg/...

IMAGE_REGISTRY ?= reg.vesoft-inc.com
IMAGE_REPO ?= ${IMAGE_REGISTRY}/cloud-dev
IMAGE_TAG ?= v3.7.0

BUILDTARGET := bin/br-ent
SUFFIX := $(GOEXE)

GITSHA := $(shell git describe --no-match --always --dirty)
GITREF := $(shell git rev-parse --abbrev-ref HEAD)

REPO := github.com/vesoft-inc/br-ent

LDFLAGS += -X $(REPO)/pkg/version.GitSha=$(GITSHA)
LDFLAGS += -X $(REPO)/pkg/version.GitRef=$(GITREF)

docker-build: ## Build docker images.
	docker build -t "${IMAGE_REPO}/br-ent:${IMAGE_TAG}" .

docker-push: ## Push docker images.
	docker push "${IMAGE_REPO}/br-ent:${IMAGE_TAG}"

build:
	$(GO) build -trimpath -ldflags '$(LDFLAGS)' -o $(BUILDTARGET) main.go
	chmod +x $(BUILDTARGET)

test:
	$(GO) test -v $(PKG) -short

fmt:
	$(GO) mod tidy && find . -path vendor -prune -o -type f -iname '*.go' -exec goimports -w {} \;