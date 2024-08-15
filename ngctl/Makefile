metaCtl ?= ngctl

# build with verison infos
buildTime = $(shell TZ=UTC date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%h' -1)
gitTag = $(shell git describe --tags --abbrev=0 --tag 2>/dev/null)

ldflags="-X github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/version.GitTag=${gitTag} \
-X github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/version.BuildTime=${buildTime} \
-X github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/version.GitCommit=${gitCommit}"

.PHONY: build vendorbuild clean fmt gen

default: build

build: clean fmt
	@CGO_ENABLED=0 go build -o ${metaCtl} -ldflags ${ldflags} ./

vendorbuild: clean fmt
	@CGO_ENABLED=0 go mod vendor && go build -mod vendor -o ${metaCtl} --tags netgo -ldflags ${ldflags} ./

clean:
	@rm -rf ${metaCtl} vendor

clean-all:
	@rm -rf ${metaCtl} vendor box/blob.go

fmt:
	@go mod tidy && find . -path vendor -prune -o -type f -iname '*.go' -exec go fmt {} \;

# generate box/blob.go(NOTE: `go generate` command may doesn't support cross-platform)
gen:
	@CGO_ENABLED=0 go generate ./...
