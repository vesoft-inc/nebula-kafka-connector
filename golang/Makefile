.PHONY: build unit test e2e fmt run-examples

default: build

build: fmt
	go mod tidy
	go build

test:
	go mod tidy
	go test -v -race --covermode=atomic --coverprofile coverage.out

e2e:
	@echo "TODO: Add e2e tests"

fmt:
	go fmt $(shell go list ./... | grep -v /generated_code/)

run-examples:
	go run examples/basic_example.go
