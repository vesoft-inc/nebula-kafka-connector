.PHONY: build unit test fmt run-examples

default: build

build: fmt
	go mod tidy
	go build

unit:
	go mod tidy
	go test -v -race --covermode=atomic --coverprofile coverage.out

test:
	go mod tidy
	go test -v -race --tags=integration --covermode=atomic --coverprofile coverage.out

fmt:
	go fmt $(shell go list ./... | grep -v /generated_code/)

run-examples:
	go run examples/basic_example.go
