# Node: should build in nebula-ng-tools folder
# e.g.
# cd nebula-ng-tools
# docker build -t reg.vesoft-inc.com/vesoft-ng/nebula-console:nightly -f console/Dockerfile  --target nebula-console .
FROM golang:1.20.3-alpine as builder

COPY console /usr/src
COPY golang /usr/golang

RUN cd /usr/src && apk add --no-cache git make && make

FROM alpine as nebula-console

COPY --from=builder /usr/src/nebula-console /usr/local/bin/nebula-console
COPY --from=builder /usr/src/data/ /data/


ENTRYPOINT ["nebula-console"]


FROM alpine as meta-console

COPY --from=builder /usr/src/meta-console /usr/local/bin/meta-console
COPY --from=builder /usr/src/data/ /data/

ENTRYPOINT ["meta-console"]
