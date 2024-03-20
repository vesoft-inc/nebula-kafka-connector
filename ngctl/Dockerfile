# Node: should build in nebula-ng-tools folder
# e.g.
# cd nebula-ng-tools
# docker build -t reg.vesoft-inc.com/vesoft-ng/ngctl:nightly -f console/Dockerfile  --target ngctl .
FROM golang:1.20.3-alpine as builder

COPY ngctl /usr/src
COPY golang /usr/golang

RUN cd /usr/src && apk add --no-cache git make && make

FROM alpine as ngctl

COPY --from=builder /usr/src/ngctl /usr/local/bin/ngctl

ENTRYPOINT ["ngctl"]
