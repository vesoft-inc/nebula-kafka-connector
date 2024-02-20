export GO111MODULE := on
GOENV  := GO15VENDOREXPERIMENT="1" CGO_ENABLED=0
GO     := $(GOENV) go
GO_BUILD := $(GO) build -trimpath
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

build:
	$(GO_BUILD) -ldflags '$(LDFLAGS)' -o bin/ngadmin ./cmd/

test:
	go test -gcflags=all="-l" -v -covermode=atomic ./pkg/... 

test-install:build
	./bin/ngadmin install -c ./examples/nebula.yaml --force

test-uninstall:build
	./bin/ngadmin uninstall -c ./examples/nebula.yaml

test-verify:build
	./bin/ngadmin verify -c ./examples/nebula.yaml