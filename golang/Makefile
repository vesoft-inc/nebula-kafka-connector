.PHONY: build unit test e2e fmt run-examples

default: build

build: fmt
	go mod tidy
	go build

test:
	go mod tidy
	go test -v -race -timeout 30s  --covermode=atomic --coverprofile coverage.out

e2e-up:
	cd ../docker-compose && docker-compose pull && docker-compose up -d

e2e-down:
	cd ../docker-compose && docker-compose down

e2e:
	cd e2e && go test -v -race

fmt:
	go fmt $(shell go list ./... | grep -v /generated_code/)

run-examples:
	go run examples/basic_example.go
