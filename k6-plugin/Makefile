all: build
.PHONY: build

pairs := darwin/amd64 linux/amd64 linux/arm64
GOPATH ?= ~/go
export GO111MODULE=on
VERSION ?= v1.0.1
K6_VERSION ?= v0.42.0

fmt:
	find . -name '*.go' -exec gofmt -s -w {} +

lint :
	golangci-lint run --out-format=tab ./...

build: 
	go install go.k6.io/xk6/cmd/xk6@v0.9.2
	$(GOPATH)/bin/xk6 build $(K6_VERSION) \ 
		--with github.com/vesoft-inc/nebula-ng-tools/k6-plugin@$(VERSION) \
		--replace github.com/vesoft-inc/nebula-ng-tools/golang=${PWD}/../golang; 

build-all: build-arm-v7

	go install go.k6.io/xk6/cmd/xk6@v0.9.2
	for  pair in $(pairs);do echo $$pair; \
			os=`echo $$pair | cut -d / -f 1 ` ;\
			arch=`echo $$pair | cut -d / -f 2 ` ;\
			GOOS=$$os GOARCH=$$arch $(GOPATH)/bin/xk6 build $(K6_VERSION) \
				--with github.com/vesoft-inc/nebula-ng-tools/k6-plugin@$(VERSION) \
				--replace github.com/vesoft-inc/nebula-ng-tools/golang=${PWD}/../golang; 
			mv k6 k6-$$os-$$arch; \
	done

build-arm-v7:
	go install go.k6.io/xk6/cmd/xk6@v0.9.2
	GOOS=linux GOARCH=arm64 GOARM=7 $(GOPATH)/bin/xk6 build $(K6_VERSION) \
		--with github.com/vesoft-inc/nebula-ng-tools/k6-plugin@$(VERSION) \
		--replace github.com/vesoft-inc/nebula-ng-tools/golang=${PWD}/../golang; 
	mv k6 k6-linux-arm64-v7

build-dev:
	go install go.k6.io/xk6/cmd/xk6@v0.9.2
	$(GOPATH)/bin/xk6 build $(K6_VERSION) \
		--with github.com/vesoft-inc/nebula-ng-tools/k6-plugin@latest=${PWD}/../k6-plugin \
		--replace github.com/vesoft-inc/nebula-ng-tools/golang=${PWD}/../golang; 	

build-dev-all:
	go install go.k6.io/xk6/cmd/xk6@v0.9.2
	for  pair in $(pairs);do echo $$pair; \
			os=`echo $$pair | cut -d / -f 1 ` ;\
			arch=`echo $$pair | cut -d / -f 2 ` ;\
			GOOS=$$os GOARCH=$$arch $(GOPATH)/bin/xk6 build $(K6_VERSION) \
				--with github.com/vesoft-inc/nebula-ng-tools/k6-plugin@latest=${PWD}/../k6-plugin \
				--replace github.com/vesoft-inc/nebula-ng-tools/golang=${PWD}/../golang; 	

			mv k6 k6-$$os-$$arch; \
	done	

.PHONY: format
format:
	find . -name '*.go' -exec gofmt -s -w {} +
