# Node: should build in nebula-ng-tools folder
# e.g.
# cd nebula-ng-tools
# docker build -t reg.vesoft-inc.com/vesoft-ng/ngctl:nightly -f console/Dockerfile  --target ngctl .
FROM golang:1.22.4-alpine as builder

COPY ngctl /usr/src
COPY golang /usr/golang
COPY ngadm /usr/ngadm
COPY br-ent /usr/br-ent
COPY agent /usr/agent

# https://github.com/alpinelinux/docker-alpine/issues/98
RUN sed -i 's/https/http/' /etc/apk/repositories
RUN cd /usr/src && apk add --no-cache git make && make

FROM alpine as ngctl

COPY --from=builder /usr/src/ngctl /usr/local/bin/ngctl

ENTRYPOINT ["ngctl"]
