# Node: should build in nebula-ng-tools folder
# e.g.
# cd nebula-ng-tools
# docker build -t reg.vesoft-inc.com/vesoft-ng/ngql:nightly -f ngql/Dockerfile  --target ngql .
FROM golang:1.20.3-alpine as builder

COPY ngql /usr/src
COPY golang /usr/golang

RUN cd /usr/src && apk add --no-cache git make && make

FROM alpine as ngql

COPY --from=builder /usr/src/ngql /usr/local/bin/ngql

ENTRYPOINT ["ngql"]
